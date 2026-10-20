package analyzer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/nodelike/diffloc/internal/model"
	"golang.org/x/sync/errgroup"
)

// emptyTree is git's well-known empty tree object, used as the base before the first commit.
const emptyTree = "4b825dc642cb6eb9a060e54bf8d69288fbee4904"

// IsGitRepo reports whether path is inside a git work tree and git is installed.
func IsGitRepo(ctx context.Context, path string) bool {
	out, err := runGit(ctx, path, "rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(string(out)) == "true"
}

// AnalyzeGit reports per-file line counts for path, with additions and deletions
// of the working tree (staged and unstaged) against HEAD. Paths are relative to path,
// so running inside a subdirectory scopes the report to it.
func AnalyzeGit(ctx context.Context, path string, filter *Filter) (*model.Stats, error) {
	base := "HEAD"
	if _, err := runGit(ctx, path, "rev-parse", "--verify", "--quiet", "HEAD^{commit}"); err != nil {
		base = emptyTree
	}

	untrackedArgs := []string{"ls-files", "-z", "--others"}
	if filter.respectGitignore {
		untrackedArgs = append(untrackedArgs, "--exclude-standard")
	}

	var numstat, untrackedOut, trackedOut []byte
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() (err error) {
		numstat, err = runGit(egCtx, path, "diff", "--numstat", "-z", "--relative", "--no-renames", "--no-ext-diff", "--no-textconv", base)
		return err
	})
	eg.Go(func() (err error) {
		untrackedOut, err = runGit(egCtx, path, untrackedArgs...)
		return err
	})
	eg.Go(func() (err error) {
		trackedOut, err = runGit(egCtx, path, "ls-files", "-z", "--cached", "--stage")
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	changes, err := parseNumstat(numstat)
	if err != nil {
		return nil, err
	}

	type entry struct {
		path      string
		changed   bool
		untracked bool
		diff      numstatEntry
	}
	var entries []entry
	seen := make(map[string]bool)
	add := func(e entry) {
		if seen[e.path] || !filter.ShouldInclude(e.path) {
			return
		}
		seen[e.path] = true
		entries = append(entries, e)
	}

	for p, d := range changes {
		add(entry{path: p, changed: true, diff: d})
	}
	for _, p := range splitNul(untrackedOut) {
		add(entry{path: p, changed: true, untracked: true})
	}
	for _, rec := range splitNul(trackedOut) {
		// "<mode> <object> <stage>\t<path>"; gitlinks (submodules) have mode 160000.
		meta, p, ok := strings.Cut(rec, "\t")
		if !ok || strings.HasPrefix(meta, "160000 ") {
			continue
		}
		add(entry{path: p})
	}

	paths := make([]string, len(entries))
	for i, e := range entries {
		paths[i] = e.path
	}
	counts, err := countAll(ctx, path, paths)
	if err != nil {
		return nil, err
	}

	stats := &model.Stats{IsGitRepo: true}
	for i, e := range entries {
		c := counts[i]
		if c.err != nil {
			stats.Skipped = append(stats.Skipped, e.path)
			continue
		}
		if c.missing && !e.changed {
			// Tracked but absent without a diff, e.g. outside a sparse checkout.
			continue
		}

		info := &model.FileInfo{Path: e.path, Lines: c.lines, IsChanged: e.changed}
		switch {
		case e.untracked:
			info.Additions = c.lines
		case e.changed:
			info.Additions = e.diff.additions
			info.Deletions = e.diff.deletions
		}

		if e.changed {
			stats.ChangedFiles = append(stats.ChangedFiles, info)
		} else {
			stats.UnchangedFiles = append(stats.UnchangedFiles, info)
		}
	}

	finalize(stats)
	return stats, nil
}

type numstatEntry struct {
	additions int
	deletions int
}

// parseNumstat parses `git diff --numstat -z --no-renames`: "<add>\t<del>\t<path>\0" per file.
// Binary files report "-" for both counts and are recorded as zero.
func parseNumstat(out []byte) (map[string]numstatEntry, error) {
	changes := make(map[string]numstatEntry)
	for _, rec := range splitNul(out) {
		fields := strings.SplitN(rec, "\t", 3)
		if len(fields) != 3 {
			return nil, fmt.Errorf("unexpected git numstat record %q", rec)
		}
		var e numstatEntry
		if fields[0] != "-" {
			n, err := strconv.Atoi(fields[0])
			if err != nil {
				return nil, fmt.Errorf("unexpected git numstat record %q", rec)
			}
			e.additions = n
		}
		if fields[1] != "-" {
			n, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, fmt.Errorf("unexpected git numstat record %q", rec)
			}
			e.deletions = n
		}
		changes[fields[2]] = e
	}
	return changes, nil
}

func splitNul(out []byte) []string {
	out = bytes.TrimSuffix(out, []byte{0})
	if len(out) == 0 {
		return nil
	}
	return strings.Split(string(out), "\x00")
}

// runGit runs git in dir without taking optional locks, so it never blocks or
// disturbs a concurrent git process in the same repository.
func runGit(ctx context.Context, dir string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"--no-optional-locks", "-C", dir}, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("git %s: %s", args[0], strings.TrimSpace(stderr.String()))
		}
		return nil, fmt.Errorf("git %s: %w", args[0], err)
	}
	return out, nil
}
