package analyzer

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nodelike/diffloc/internal/model"
)

func gitCmd(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir, "-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false"}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func write(t *testing.T, root, rel, content string) {
	t.Helper()
	p := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func newRepo(t *testing.T) string {
	t.Helper()
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	dir := t.TempDir()
	gitCmd(t, dir, "init", "-q")
	return dir
}

func lines(n int, prefix string) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString(prefix)
		b.WriteString(strings.Repeat("x", i%7))
		b.WriteString("\n")
	}
	return b.String()
}

func mustFilter(t *testing.T, opts FilterOptions) *Filter {
	t.Helper()
	opts.RespectGitignore = true
	f, err := NewFilter(opts)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func analyze(t *testing.T, dir string, opts FilterOptions) *model.Stats {
	t.Helper()
	stats, err := Analyze(context.Background(), dir, mustFilter(t, opts))
	if err != nil {
		t.Fatal(err)
	}
	return stats
}

func byPath(stats *model.Stats) map[string]*model.FileInfo {
	m := make(map[string]*model.FileInfo)
	for _, f := range append(append([]*model.FileInfo{}, stats.ChangedFiles...), stats.UnchangedFiles...) {
		m[f.Path] = f
	}
	return m
}

func expect(t *testing.T, f *model.FileInfo, path string, lines, add, del int, changed bool) {
	t.Helper()
	if f == nil {
		t.Fatalf("%s: missing from results", path)
	}
	if f.Lines != lines || f.Additions != add || f.Deletions != del || f.IsChanged != changed {
		t.Errorf("%s: got lines=%d +%d -%d changed=%v, want lines=%d +%d -%d changed=%v",
			path, f.Lines, f.Additions, f.Deletions, f.IsChanged, lines, add, del, changed)
	}
}

func TestGitRepo(t *testing.T) {
	dir := newRepo(t)
	write(t, dir, "same.go", "package a\n\nfunc A() {}\n")
	write(t, dir, "moved.go", "a\nb\nc\nd\n")
	write(t, dir, "edit.go", "one\ntwo\nthree\n")
	write(t, dir, "gone.go", lines(5, "g"))
	write(t, dir, "staged_gone.go", lines(3, "s"))
	write(t, dir, ".gitignore", "ignored.go\n")
	write(t, dir, "node_modules/x.js", "x\n")
	gitCmd(t, dir, "add", "-A")
	gitCmd(t, dir, "commit", "-qm", "init")

	write(t, dir, "moved.go", "c\nd\na\nb\n")
	write(t, dir, "edit.go", "one\n2\nthree\nfour")
	write(t, dir, "new.go", "n1\nn2")
	write(t, dir, "ignored.go", "i\n")
	write(t, dir, "image.png", "png")
	os.Remove(filepath.Join(dir, "gone.go"))
	gitCmd(t, dir, "rm", "-q", "staged_gone.go")

	stats := analyze(t, dir, FilterOptions{})
	if !stats.IsGitRepo {
		t.Fatal("expected git repo")
	}
	files := byPath(stats)
	expect(t, files["same.go"], "same.go", 3, 0, 0, false)
	expect(t, files["moved.go"], "moved.go", 4, 2, 2, true)
	expect(t, files["edit.go"], "edit.go", 4, 2, 1, true)
	expect(t, files["new.go"], "new.go", 2, 2, 0, true)
	expect(t, files["gone.go"], "gone.go", 0, 0, 5, true)
	expect(t, files["staged_gone.go"], "staged_gone.go", 0, 0, 3, true)
	for _, p := range []string{"ignored.go", "image.png", "node_modules/x.js"} {
		if files[p] != nil {
			t.Errorf("%s should be excluded", p)
		}
	}
	if len(files) != 6 {
		t.Errorf("got %d files, want 6", len(files))
	}
	if stats.TotalAdditions != 6 || stats.TotalDeletions != 11 || stats.NetChange != -5 {
		t.Errorf("totals: +%d -%d net %d", stats.TotalAdditions, stats.TotalDeletions, stats.NetChange)
	}
	for i := 1; i < len(stats.ChangedFiles); i++ {
		if stats.ChangedFiles[i-1].Path > stats.ChangedFiles[i].Path {
			t.Error("changed files not sorted by path")
		}
	}
}

func TestGitSubdirectory(t *testing.T) {
	dir := newRepo(t)
	write(t, dir, "root.go", "r\n")
	write(t, dir, "pkg/a.go", "a\n")
	gitCmd(t, dir, "add", "-A")
	gitCmd(t, dir, "commit", "-qm", "init")
	write(t, dir, "pkg/a.go", "a\nb\n")
	write(t, dir, "pkg/sub/b.go", "b\n")
	write(t, dir, "root.go", "changed\n")

	stats := analyze(t, filepath.Join(dir, "pkg"), FilterOptions{})
	if !stats.IsGitRepo {
		t.Fatal("subdirectory of a repo should be analyzed with git")
	}
	files := byPath(stats)
	if len(files) != 2 {
		t.Fatalf("got %v, want only pkg files", files)
	}
	expect(t, files["a.go"], "a.go", 2, 1, 0, true)
	expect(t, files["sub/b.go"], "sub/b.go", 1, 1, 0, true)
}

func TestGitNoCommits(t *testing.T) {
	dir := newRepo(t)
	write(t, dir, "a.go", "a\nb\n")
	write(t, dir, "staged.go", "s\n")
	gitCmd(t, dir, "add", "staged.go")

	files := byPath(analyze(t, dir, FilterOptions{}))
	expect(t, files["a.go"], "a.go", 2, 2, 0, true)
	expect(t, files["staged.go"], "staged.go", 1, 1, 0, true)
}

func TestMaxDepthAndTests(t *testing.T) {
	dir := newRepo(t)
	write(t, dir, "a.go", "a\n")
	write(t, dir, "a_test.go", "a\n")
	write(t, dir, "x/b.go", "b\n")
	write(t, dir, "x/y/c.go", "c\n")
	write(t, dir, "tests/d.go", "d\n")
	gitCmd(t, dir, "add", "-A")
	gitCmd(t, dir, "commit", "-qm", "init")

	files := byPath(analyze(t, dir, FilterOptions{MaxDepth: 2, ExcludeTests: true}))
	if len(files) != 2 || files["a.go"] == nil || files["x/b.go"] == nil {
		t.Fatalf("got %v, want a.go and x/b.go", files)
	}
}

func TestPlainDirectory(t *testing.T) {
	t.Setenv("GIT_CEILING_DIRECTORIES", os.TempDir())
	dir := t.TempDir()
	write(t, dir, "a.go", "a\nb")
	write(t, dir, "bin.go", "x\x00y\n")
	write(t, dir, "build/out.go", "o\n")
	write(t, dir, "gen/skip.go", "s\n")
	write(t, dir, "deep/gen/skip2.go", "s\n")
	write(t, dir, "keep/k.py", "k\n")
	write(t, dir, ".gitignore", "gen/\n")

	stats := analyze(t, dir, FilterOptions{})
	if stats.IsGitRepo {
		t.Fatal("expected plain directory")
	}
	files := byPath(stats)
	if len(files) != 3 {
		t.Fatalf("got %v", files)
	}
	expect(t, files["a.go"], "a.go", 2, 0, 0, false)
	expect(t, files["bin.go"], "bin.go", 0, 0, 0, false)
	expect(t, files["keep/k.py"], "keep/k.py", 1, 0, 0, false)
}

func TestInvalidExclude(t *testing.T) {
	if _, err := NewFilter(FilterOptions{Excludes: []string{"("}}); err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestCountLines(t *testing.T) {
	dir := t.TempDir()
	cases := map[string]int{"": 0, "a": 1, "a\n": 1, "a\nb": 2, "\n\n": 2, strings.Repeat("x\n", 100000): 100000}
	for content, want := range cases {
		p := filepath.Join(dir, "f")
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		got, err := CountLines(p)
		if err != nil || got != want {
			t.Errorf("CountLines(%.10q) = %d, %v; want %d", content, got, err, want)
		}
	}
}
