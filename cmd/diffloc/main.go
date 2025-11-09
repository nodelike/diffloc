package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nodelike/diffloc/internal/analyzer"
	"github.com/nodelike/diffloc/internal/ui"
)

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	var (
		noGitignore    bool
		customExcludes stringSliceFlag
		allowedExts    stringSliceFlag
		path           string
	)

	flag.BoolVar(&noGitignore, "no-gitignore", false, "Ignore .gitignore patterns (always-excluded patterns still apply)")
	flag.Var(&customExcludes, "exclude", "Additional exclusion pattern (can be repeated)")
	flag.Var(&allowedExts, "ext", "Override allowed file extensions (can be repeated)")
	flag.StringVar(&path, "path", ".", "Path to analyze (defaults to current directory)")
	flag.Parse()

	// Get working directory if path is relative
	if path == "." {
		var err error
		path, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to get current directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Create filter
	filter := analyzer.NewFilter(allowedExts, customExcludes, !noGitignore)

	// Load gitignore if in a git repo and respecting gitignore
	if !noGitignore && analyzer.IsGitRepo(path) {
		repoRoot, err := analyzer.GetRepoRoot(path)
		if err == nil {
			filter.LoadGitignore(repoRoot)
		}
	}

	// Analyze the directory
	stats, err := analyzer.Analyze(path, filter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Run TUI
	if err := ui.Run(stats); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}

