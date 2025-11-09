package analyzer

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/nodelike/diffloc/internal/model"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

// AnalyzeFiles analyzes files in a non-git directory
func AnalyzeFiles(ctx context.Context, rootPath string, filter *Filter) (*model.Stats, error) {
	stats := &model.Stats{
		ChangedFiles:   make([]*model.FileInfo, 0),
		UnchangedFiles: make([]*model.FileInfo, 0),
	}

	// Collect file paths first
	type fileJob struct {
		fullPath string
		relPath  string
	}
	
	fileJobs := make([]fileJob, 0)
	
	err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Skip files with errors
		}

		// Skip directories
		if d.IsDir() {
			// Check if directory should be excluded
			relPath, _ := filepath.Rel(rootPath, path)
			if relPath != "." && !filter.ShouldInclude(relPath+"/dummy.go") {
				return filepath.SkipDir
			}
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			return nil
		}

		// Check if file should be included
		if !filter.ShouldInclude(relPath) {
			return nil
		}

		fileJobs = append(fileJobs, fileJob{fullPath: path, relPath: relPath})
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Process files in parallel
	var statsMu sync.Mutex
	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(16) // Limit concurrent workers

	// Show progress bar for large repos
	totalFiles := len(fileJobs)
	var bar *progressbar.ProgressBar
	if totalFiles > 1000 {
		bar = progressbar.NewOptions(totalFiles,
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionSetDescription("Analyzing files"),
			progressbar.OptionShowCount(),
			progressbar.OptionSetWidth(40),
			progressbar.OptionThrottle(100),
		)
		defer bar.Finish()
	}

	for _, job := range fileJobs {
		job := job // Capture loop variable
		eg.Go(func() error {
			// Check for context cancellation
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			default:
			}
			// Count lines
			lines, err := CountLines(job.fullPath)
			if err != nil {
				return nil // Skip files with errors
			}

			fileInfo := &model.FileInfo{
				Path:      job.relPath,
				Lines:     lines,
				Additions: 0,
				Deletions: 0,
				IsChanged: false,
			}

			statsMu.Lock()
			stats.UnchangedFiles = append(stats.UnchangedFiles, fileInfo)
			stats.TotalLines += lines
			if bar != nil {
				bar.Add(1)
			}
			statsMu.Unlock()

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	stats.TotalFiles = len(stats.UnchangedFiles)
	stats.UnchangedCount = len(stats.UnchangedFiles)

	return stats, nil
}

