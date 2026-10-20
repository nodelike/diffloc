package analyzer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gobwas/glob"
)

// FilterOptions configures which files are counted.
type FilterOptions struct {
	Exts             []string
	Excludes         []string
	RespectGitignore bool
	ExcludeTests     bool
	MaxDepth         int // 0 = unlimited; 1 = only files directly in the root
}

// Filter handles file exclusion logic
type Filter struct {
	allowedExts      map[string]bool
	excludePatterns  []*regexp.Regexp
	gitignoreGlobs   []glob.Glob
	respectGitignore bool
	excludeTests     bool
	maxDepth         int
}

var defaultExts = []string{
	// Go
	".go",
	// Python
	".py", ".pyi", ".pyw",
	// JavaScript / TypeScript
	".js", ".jsx", ".ts", ".tsx", ".mjs", ".cjs", ".mts", ".cts",
	// Web components
	".vue", ".svelte", ".astro",
	// Rust
	".rs",
	// C / C++
	".c", ".h", ".cc", ".cpp", ".cxx", ".hpp", ".hxx", ".hh",
	// JVM languages
	".java", ".kt", ".kts", ".scala", ".groovy", ".clj", ".cljs", ".cljc",
	// .NET languages
	".cs", ".fs", ".fsx", ".vb",
	// Apple platforms
	".swift", ".m", ".mm",
	// Other popular languages
	".rb", ".php", ".dart", ".lua", ".pl", ".pm",
	".ex", ".exs", ".erl", ".hrl", ".hs", ".ml", ".mli",
	".jl", ".r", ".nim", ".zig",
	// Shell scripting
	".sh", ".bash", ".zsh", ".fish", ".ps1",
	// Styles
	".css", ".scss", ".sass", ".less", ".styl",
	// Markup
	".html", ".htm",
	// Data / query / schema
	".sql", ".graphql", ".gql", ".proto",
}

// excludedDirs are directory names skipped at any depth.
var excludedDirs = map[string]bool{
	"node_modules": true, "venv": true, ".venv": true, "__pycache__": true,
	".git": true, "dist": true, "build": true, "target": true, ".tox": true,
	"coverage": true, ".next": true, "vendor": true, "bin": true, "tmp": true,
	".eggs": true, ".pytest_cache": true, ".mypy_cache": true, ".nyc_output": true,
	".yarn": true, ".npm": true, "cypress": true, "jest-cache": true,
}

var testDirs = map[string]bool{"test": true, "tests": true}

// excludedSuffixes are generated, vendored, or binary file name endings.
var excludedSuffixes = []string{
	".lock", "-lock.json", "-lock.yaml", ".gitignore", ".eslintcache",
	".exe", ".so", ".dylib", ".dll", ".pyc", ".pyo", ".pyd", ".egg",
	"_templ.go", ".pb.go", "_gen.go", ".min.js", ".bundle.js",
	".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".ico", ".webp",
	".tiff", ".tif", ".psd", ".raw", ".heic", ".avif",
}

var testSuffixes = []string{
	"_test.go",
	".test.js", ".test.ts", ".test.jsx", ".test.tsx",
	".spec.js", ".spec.ts", ".spec.jsx", ".spec.tsx",
}

// NewFilter creates a filter. It fails if a custom exclude is not a valid regex.
func NewFilter(opts FilterOptions) (*Filter, error) {
	f := &Filter{
		allowedExts:      make(map[string]bool),
		respectGitignore: opts.RespectGitignore,
		excludeTests:     opts.ExcludeTests,
		maxDepth:         opts.MaxDepth,
	}

	exts := opts.Exts
	if len(exts) == 0 {
		exts = defaultExts
	}
	for _, ext := range exts {
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		f.allowedExts[strings.ToLower(ext)] = true
	}

	for _, pattern := range opts.Excludes {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid exclude pattern %q: %w", pattern, err)
		}
		f.excludePatterns = append(f.excludePatterns, re)
	}

	return f, nil
}

// LoadGitignore parses the root .gitignore. Only used outside git repos;
// inside a repo git itself applies every ignore source.
func (f *Filter) LoadGitignore(root string) error {
	if !f.respectGitignore {
		return nil
	}

	file, err := os.Open(filepath.Join(root, ".gitignore"))
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "!") {
			continue
		}

		pattern := strings.TrimSuffix(line, "/")
		anchored := strings.Contains(pattern, "/")
		pattern = strings.TrimPrefix(pattern, "/")

		candidates := []string{pattern, pattern + "/**"}
		if !anchored {
			candidates = append(candidates, "**/"+pattern, "**/"+pattern+"/**")
		}
		for _, c := range candidates {
			if g, err := glob.Compile(c, '/'); err == nil {
				f.gitignoreGlobs = append(f.gitignoreGlobs, g)
			}
		}
	}

	return scanner.Err()
}

// ShouldInclude reports whether a file (slash or OS separated, relative to the root) is counted.
func (f *Filter) ShouldInclude(path string) bool {
	path = filepath.ToSlash(path)

	if !f.allowedExts[strings.ToLower(filepath.Ext(path))] {
		return false
	}

	dirs := strings.Split(path, "/")
	name := dirs[len(dirs)-1]
	dirs = dirs[:len(dirs)-1]

	if f.maxDepth > 0 && len(dirs) >= f.maxDepth {
		return false
	}
	if f.excludedDirPath(dirs) {
		return false
	}

	for _, suffix := range excludedSuffixes {
		if strings.HasSuffix(name, suffix) {
			return false
		}
	}
	if f.excludeTests {
		for _, suffix := range testSuffixes {
			if strings.HasSuffix(name, suffix) {
				return false
			}
		}
	}

	return !f.matchesPatterns(path)
}

// ShouldDescend reports whether a directory walk should enter dir (relative to the root).
func (f *Filter) ShouldDescend(dir string) bool {
	dir = filepath.ToSlash(dir)
	dirs := strings.Split(dir, "/")
	if f.maxDepth > 0 && len(dirs) >= f.maxDepth {
		return false
	}
	if f.excludedDirPath(dirs) {
		return false
	}
	return !f.matchesPatterns(dir + "/")
}

func (f *Filter) excludedDirPath(dirs []string) bool {
	for _, d := range dirs {
		if excludedDirs[d] || strings.HasSuffix(d, ".egg-info") || (f.excludeTests && testDirs[d]) {
			return true
		}
	}
	return false
}

func (f *Filter) matchesPatterns(path string) bool {
	for _, re := range f.excludePatterns {
		if re.MatchString(path) {
			return true
		}
	}
	if f.respectGitignore {
		trimmed := strings.TrimSuffix(path, "/")
		for _, g := range f.gitignoreGlobs {
			if g.Match(trimmed) {
				return true
			}
		}
	}
	return false
}
