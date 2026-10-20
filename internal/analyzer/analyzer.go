package analyzer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"

	"github.com/nodelike/diffloc/internal/model"
)

// Analyze counts lines under path, using git when path is inside a work tree.
func Analyze(ctx context.Context, path string, filter *Filter) (*model.Stats, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", path)
	}

	if _, err := exec.LookPath("git"); err == nil && IsGitRepo(ctx, path) {
		return AnalyzeGit(ctx, path, filter)
	}
	return AnalyzeFiles(ctx, path, filter)
}

// finalize sorts files by path and fills in the totals.
func finalize(stats *model.Stats) {
	byPath := func(files []*model.FileInfo) {
		sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	}
	byPath(stats.ChangedFiles)
	byPath(stats.UnchangedFiles)
	sort.Strings(stats.Skipped)

	if stats.ChangedFiles == nil {
		stats.ChangedFiles = []*model.FileInfo{}
	}
	if stats.UnchangedFiles == nil {
		stats.UnchangedFiles = []*model.FileInfo{}
	}

	for _, files := range [][]*model.FileInfo{stats.ChangedFiles, stats.UnchangedFiles} {
		for _, f := range files {
			stats.TotalLines += f.Lines
			stats.TotalAdditions += f.Additions
			stats.TotalDeletions += f.Deletions
		}
	}
	stats.ChangedCount = len(stats.ChangedFiles)
	stats.UnchangedCount = len(stats.UnchangedFiles)
	stats.TotalFiles = stats.ChangedCount + stats.UnchangedCount
	stats.NetChange = stats.TotalAdditions - stats.TotalDeletions
}
