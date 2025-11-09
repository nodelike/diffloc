package analyzer

import (
	"context"

	"github.com/nodelike/diffloc/internal/model"
)

// Analyzer defines the interface for analyzing file statistics
type Analyzer interface {
	Analyze(ctx context.Context, rootPath string, filter *Filter) (*model.Stats, error)
}

// GitAnalyzer implements Analyzer for Git repositories
type GitAnalyzer struct{}

// NewGitAnalyzer creates a new GitAnalyzer
func NewGitAnalyzer() *GitAnalyzer {
	return &GitAnalyzer{}
}

func (g *GitAnalyzer) Analyze(ctx context.Context, rootPath string, filter *Filter) (*model.Stats, error) {
	return AnalyzeGit(ctx, rootPath, filter)
}

// FileAnalyzer implements Analyzer for non-Git directories
type FileAnalyzer struct{}

func NewFileAnalyzer() *FileAnalyzer {
	return &FileAnalyzer{}
}

func (f *FileAnalyzer) Analyze(ctx context.Context, rootPath string, filter *Filter) (*model.Stats, error) {
	return AnalyzeFiles(ctx, rootPath, filter)
}

// GetAnalyzer returns the appropriate analyzer based on whether the path is a Git repository
func GetAnalyzer(rootPath string) Analyzer {
	if IsGitRepo(rootPath) {
		return NewGitAnalyzer()
	}
	return NewFileAnalyzer()
}

