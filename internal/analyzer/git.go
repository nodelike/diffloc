package analyzer

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/nodelike/diffloc/internal/model"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

// AnalyzeGit analyzes a git repository for changes
func AnalyzeGit(ctx context.Context, rootPath string, filter *Filter) (*model.Stats, error) {
	// Open the repository
	repo, err := git.PlainOpen(rootPath)
	if err != nil {
		return nil, err
	}

	// Get the worktree
	worktree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	// Get the HEAD reference
	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	// Get HEAD commit (will be used in calculateDiffNative)
	headCommit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return nil, err
	}

	// Get HEAD tree
	headTree, err := headCommit.Tree()
	if err != nil {
		return nil, err
	}

	// Get worktree status
	status, err := worktree.Status()
	if err != nil {
		return nil, err
	}

	stats := &model.Stats{
		ChangedFiles:   make([]*model.FileInfo, 0),
		UnchangedFiles: make([]*model.FileInfo, 0),
	}

	changedPaths := make(map[string]bool)
	var statsMu sync.Mutex

	// Collect changed files first
	type changedFileJob struct {
		path       string
		fileStatus *git.FileStatus
	}

	changedJobs := make([]changedFileJob, 0)
	for path, fileStatus := range status {
		// Skip if file should be excluded
		if !filter.ShouldInclude(path) {
			continue
		}
		changedPaths[path] = true
		changedJobs = append(changedJobs, changedFileJob{path: path, fileStatus: fileStatus})
	}

	// Process changed files in parallel
	eg, egCtx := errgroup.WithContext(ctx)
	eg.SetLimit(16) // Limit concurrent workers

	// Show progress bar for large repos (>1000 files)
	totalFiles := len(changedJobs)
	var bar *progressbar.ProgressBar
	if totalFiles > 1000 {
		bar = progressbar.NewOptions(totalFiles,
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionSetDescription("Analyzing changed files"),
			progressbar.OptionShowCount(),
			progressbar.OptionSetWidth(40),
			progressbar.OptionThrottle(100),
		)
		defer bar.Finish()
	}

	for _, job := range changedJobs {
		job := job // Capture loop variable
		eg.Go(func() error {
			// Check for context cancellation
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			default:
			}
			fullPath := filepath.Join(rootPath, job.path)

			// Count total lines in current file
			lines, err := CountLines(fullPath)
			if err != nil {
				// File might be deleted
				lines = 0
			}

			fileInfo := &model.FileInfo{
				Path:      job.path,
				Lines:     lines,
				Additions: 0,
				Deletions: 0,
				IsChanged: true,
			}

			// Handle different file statuses
			switch {
			case job.fileStatus.Staging == git.Untracked || job.fileStatus.Worktree == git.Untracked:
				// Untracked files: all lines are additions
				fileInfo.Additions = lines
			case job.fileStatus.Worktree == git.Deleted || job.fileStatus.Staging == git.Deleted:
				// Deleted files
				fileInfo.Deletions = lines
				fileInfo.Lines = 0
			default:
				// Modified files: calculate diff using git's native diff
				additions, deletions := calculateDiffNative(repo, headCommit, job.path)
				fileInfo.Additions = additions
				fileInfo.Deletions = deletions
			}

			statsMu.Lock()
			stats.ChangedFiles = append(stats.ChangedFiles, fileInfo)
			stats.TotalLines += fileInfo.Lines
			stats.TotalAdditions += fileInfo.Additions
			stats.TotalDeletions += fileInfo.Deletions
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

	// Collect unchanged files first
	unchangedPaths := make([]string, 0)
	err = headTree.Files().ForEach(func(f *object.File) error {
		path := f.Name

		// Skip if already processed as changed
		if changedPaths[path] {
			return nil
		}

		// Skip if file should be excluded
		if !filter.ShouldInclude(path) {
			return nil
		}

		unchangedPaths = append(unchangedPaths, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Process unchanged files in parallel
	eg2, eg2Ctx := errgroup.WithContext(ctx)
	eg2.SetLimit(16)

	// Show progress bar for large repos
	totalUnchanged := len(unchangedPaths)
	var bar2 *progressbar.ProgressBar
	if totalUnchanged > 1000 {
		bar2 = progressbar.NewOptions(totalUnchanged,
			progressbar.OptionSetWriter(os.Stderr),
			progressbar.OptionSetDescription("Analyzing unchanged files"),
			progressbar.OptionShowCount(),
			progressbar.OptionSetWidth(40),
			progressbar.OptionThrottle(100),
		)
		defer bar2.Finish()
	}

	for _, path := range unchangedPaths {
		path := path // Capture loop variable
		eg2.Go(func() error {
			// Check for context cancellation
			select {
			case <-eg2Ctx.Done():
				return eg2Ctx.Err()
			default:
			}
			fullPath := filepath.Join(rootPath, path)
			lines, err := CountLines(fullPath)
			if err != nil {
				// File might not exist in worktree
				return nil
			}

			fileInfo := &model.FileInfo{
				Path:      path,
				Lines:     lines,
				Additions: 0,
				Deletions: 0,
				IsChanged: false,
			}

			statsMu.Lock()
			stats.UnchangedFiles = append(stats.UnchangedFiles, fileInfo)
			stats.TotalLines += lines
			if bar2 != nil {
				bar2.Add(1)
			}
			statsMu.Unlock()

			return nil
		})
	}

	if err := eg2.Wait(); err != nil {
		return nil, err
	}

	// Calculate totals
	stats.ChangedCount = len(stats.ChangedFiles)
	stats.UnchangedCount = len(stats.UnchangedFiles)
	stats.TotalFiles = stats.ChangedCount + stats.UnchangedCount
	stats.NetChange = stats.TotalAdditions - stats.TotalDeletions

	return stats, nil
}

// calculateDiffNative uses go-git's native diff to calculate additions and deletions
func calculateDiffNative(repo *git.Repository, headCommit *object.Commit, path string) (additions, deletions int) {
	// Get HEAD tree
	headTree, err := headCommit.Tree()
	if err != nil {
		return 0, 0
	}

	// Get HEAD file
	headFile, err := headTree.File(path)
	if err != nil {
		return 0, 0
	}

	// Get worktree
	worktree, err := repo.Worktree()
	if err != nil {
		return 0, 0
	}

	// Get worktree file
	fs := worktree.Filesystem
	worktreeFile, err := fs.Open(path)
	if err != nil {
		return 0, 0
	}
	defer worktreeFile.Close()

	// Read worktree content
	worktreeContent := strings.Builder{}
	buf := make([]byte, 4096)
	for {
		n, err := worktreeFile.Read(buf)
		if n > 0 {
			worktreeContent.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	// Get HEAD content
	headContent, err := headFile.Contents()
	if err != nil {
		return 0, 0
	}

	// Calculate diff using simple line comparison
	// Split into lines
	headLines := strings.Split(headContent, "\n")
	workLines := strings.Split(worktreeContent.String(), "\n")

	// Use a simple but more accurate diff based on longest common subsequence approach
	additions, deletions = computeLineDiff(headLines, workLines)

	return additions, deletions
}

// computeLineDiff computes additions and deletions using a simplified diff algorithm
func computeLineDiff(oldLines, newLines []string) (additions, deletions int) {
	// Simple approach: create a map of lines and their counts
	oldMap := make(map[string]int)
	newMap := make(map[string]int)

	for _, line := range oldLines {
		oldMap[line]++
	}

	for _, line := range newLines {
		newMap[line]++
	}

	// Calculate additions and deletions
	for line, oldCount := range oldMap {
		newCount := newMap[line]
		if newCount < oldCount {
			deletions += oldCount - newCount
		}
	}

	for line, newCount := range newMap {
		oldCount := oldMap[line]
		if newCount > oldCount {
			additions += newCount - oldCount
		}
	}

	return additions, deletions
}

// IsGitRepo checks if the given path is a git repository
func IsGitRepo(path string) bool {
	_, err := git.PlainOpen(path)
	return err == nil
}

// GetRepoRoot returns the root path of the git repository
func GetRepoRoot(path string) (string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", err
	}

	return worktree.Filesystem.Root(), nil
}

// Analyze is the main entry point that decides between git and non-git analysis
func Analyze(ctx context.Context, rootPath string, filter *Filter) (*model.Stats, error) {
	if IsGitRepo(rootPath) {
		return AnalyzeGit(ctx, rootPath, filter)
	}
	return AnalyzeFiles(ctx, rootPath, filter)
}
