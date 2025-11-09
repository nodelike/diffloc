package analyzer

import (
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/nodelike/diffloc/internal/model"
)

// AnalyzeGit analyzes a git repository for changes
func AnalyzeGit(rootPath string, filter *Filter) (*model.Stats, error) {
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

	// Get HEAD commit
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

	// Process changed and untracked files
	for path, fileStatus := range status {
		// Skip if file should be excluded
		if !filter.ShouldInclude(path) {
			continue
		}

		changedPaths[path] = true
		fullPath := filepath.Join(rootPath, path)

		// Count total lines in current file
		lines, err := CountLines(fullPath)
		if err != nil {
			// File might be deleted
			lines = 0
		}

		fileInfo := &model.FileInfo{
			Path:      path,
			Lines:     lines,
			Additions: 0,
			Deletions: 0,
			IsChanged: true,
		}

		// Handle different file statuses
		switch {
		case fileStatus.Staging == git.Untracked || fileStatus.Worktree == git.Untracked:
			// Untracked files: all lines are additions
			fileInfo.Additions = lines
		case fileStatus.Worktree == git.Deleted || fileStatus.Staging == git.Deleted:
			// Deleted files
			fileInfo.Deletions = lines
			fileInfo.Lines = 0
		default:
			// Modified files: calculate diff
			additions, deletions := calculateDiff(headTree, worktree, path)
			fileInfo.Additions = additions
			fileInfo.Deletions = deletions
		}

		stats.ChangedFiles = append(stats.ChangedFiles, fileInfo)
		stats.TotalLines += fileInfo.Lines
		stats.TotalAdditions += fileInfo.Additions
		stats.TotalDeletions += fileInfo.Deletions
	}

	// Process unchanged tracked files
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

		stats.UnchangedFiles = append(stats.UnchangedFiles, fileInfo)
		stats.TotalLines += lines

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Calculate totals
	stats.ChangedCount = len(stats.ChangedFiles)
	stats.UnchangedCount = len(stats.UnchangedFiles)
	stats.TotalFiles = stats.ChangedCount + stats.UnchangedCount
	stats.NetChange = stats.TotalAdditions - stats.TotalDeletions

	return stats, nil
}

// calculateDiff calculates additions and deletions for a modified file
func calculateDiff(headTree *object.Tree, worktree *git.Worktree, path string) (additions, deletions int) {
	// Get file from HEAD
	headFile, err := headTree.File(path)
	if err != nil {
		return 0, 0
	}

	headContent, err := headFile.Contents()
	if err != nil {
		return 0, 0
	}

	// Get file from worktree
	fs := worktree.Filesystem
	worktreeFile, err := fs.Open(path)
	if err != nil {
		return 0, 0
	}
	defer worktreeFile.Close()

	worktreeContent := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := worktreeFile.Read(buf)
		if n > 0 {
			worktreeContent = append(worktreeContent, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	// Simple line-by-line diff
	headLines := strings.Split(headContent, "\n")
	worktreeLines := strings.Split(string(worktreeContent), "\n")

	// Use a simple diff algorithm
	additions, deletions = simpleDiff(headLines, worktreeLines)

	return additions, deletions
}

// simpleDiff performs a basic diff between two sets of lines
func simpleDiff(oldLines, newLines []string) (additions, deletions int) {
	oldMap := make(map[string]int)
	newMap := make(map[string]int)

	for _, line := range oldLines {
		oldMap[line]++
	}

	for _, line := range newLines {
		newMap[line]++
	}

	// Count deletions (lines in old but not in new, or fewer in new)
	for line, oldCount := range oldMap {
		newCount := newMap[line]
		if newCount < oldCount {
			deletions += oldCount - newCount
		}
	}

	// Count additions (lines in new but not in old, or more in new)
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
func Analyze(rootPath string, filter *Filter) (*model.Stats, error) {
	if IsGitRepo(rootPath) {
		return AnalyzeGit(rootPath, filter)
	}
	return AnalyzeFiles(rootPath, filter)
}

