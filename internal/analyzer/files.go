package analyzer

import (
	"os"
	"path/filepath"

	"github.com/nodelike/diffloc/internal/model"
)

// AnalyzeFiles analyzes files in a non-git directory
func AnalyzeFiles(rootPath string, filter *Filter) (*model.Stats, error) {
	stats := &model.Stats{
		ChangedFiles:   make([]*model.FileInfo, 0),
		UnchangedFiles: make([]*model.FileInfo, 0),
	}

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

		// Count lines
		lines, err := CountLines(path)
		if err != nil {
			return nil
		}

		fileInfo := &model.FileInfo{
			Path:      relPath,
			Lines:     lines,
			Additions: 0,
			Deletions: 0,
			IsChanged: false,
		}

		stats.UnchangedFiles = append(stats.UnchangedFiles, fileInfo)
		stats.TotalLines += lines

		return nil
	})

	if err != nil {
		return nil, err
	}

	stats.TotalFiles = len(stats.UnchangedFiles)
	stats.UnchangedCount = len(stats.UnchangedFiles)

	return stats, nil
}

