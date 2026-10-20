package analyzer

import (
	"context"
	"io/fs"
	"path/filepath"

	"github.com/nodelike/diffloc/internal/model"
)

// AnalyzeFiles counts lines in a directory that is not a git repository.
func AnalyzeFiles(ctx context.Context, root string, filter *Filter) (*model.Stats, error) {
	if err := filter.LoadGitignore(root); err != nil {
		return nil, err
	}

	var paths []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if path == root {
				return err
			}
			return nil
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}

		rel, err := filepath.Rel(root, path)
		if err != nil || rel == "." {
			return nil
		}
		rel = filepath.ToSlash(rel)

		if d.IsDir() {
			if !filter.ShouldDescend(rel) {
				return filepath.SkipDir
			}
			return nil
		}
		if d.Type().IsRegular() && filter.ShouldInclude(rel) {
			paths = append(paths, rel)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	counts, err := countAll(ctx, root, paths)
	if err != nil {
		return nil, err
	}

	stats := &model.Stats{}
	for i, p := range paths {
		c := counts[i]
		if c.err != nil {
			stats.Skipped = append(stats.Skipped, p)
			continue
		}
		if c.missing {
			continue
		}
		stats.UnchangedFiles = append(stats.UnchangedFiles, &model.FileInfo{Path: p, Lines: c.lines})
	}

	finalize(stats)
	return stats, nil
}
