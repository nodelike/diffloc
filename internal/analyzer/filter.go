package analyzer

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gobwas/glob"
)

// Filter handles file exclusion logic
type Filter struct {
	allowedExts      map[string]bool
	excludePatterns  []*regexp.Regexp
	gitignoreGlobs   []glob.Glob
	respectGitignore bool
}

// NewFilter creates a new filter with default or custom settings
func NewFilter(allowedExts []string, customExcludes []string, respectGitignore bool) *Filter {
	f := &Filter{
		allowedExts:      make(map[string]bool),
		respectGitignore: respectGitignore,
	}

	// Set allowed extensions (default if not provided)
	if len(allowedExts) == 0 {
		defaultExts := []string{".go", ".py", ".js", ".jsx", ".ts", ".tsx", ".vue", ".svelte", ".mjs", ".cjs"}
		for _, ext := range defaultExts {
			f.allowedExts[ext] = true
		}
	} else {
		for _, ext := range allowedExts {
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			f.allowedExts[ext] = true
		}
	}

	// Always excluded patterns (regardless of flags)
	alwaysExclude := []string{
		`(^|/)node_modules($|/)`,
		`(^|/)venv($|/)`,
		`(^|/)\.venv($|/)`,
		`(^|/)__pycache__($|/)`,
		`(^|/)\.git($|/)`,
		`(^|/)dist($|/)`,
		`(^|/)build($|/)`,
		`(^|/)\.egg-info($|/)`,
		`(^|/)\.tox($|/)`,
		`(^|/)coverage($|/)`,
		`(^|/)\.next($|/)`,
		`(^|/)vendor($|/)`,
		`(^|/)bin($|/)`,
		`(^|/)tmp($|/)`,
		`\.lock$`,
		`-lock\.json$`,
		`-lock\.yaml$`,
		`Pipfile\.lock$`,
		`\.gitignore$`,
		`\.exe$`,
		`\.so$`,
		`\.dylib$`,
		`\.dll$`,
		`_templ\.go$`,
		`\.(jpg|jpeg|png|gif|bmp|svg|ico|webp|tiff|tif|psd|raw|heic|avif)$`,
	}

	f.excludePatterns = make([]*regexp.Regexp, 0, len(alwaysExclude)+len(customExcludes))
	for _, pattern := range alwaysExclude {
		if re, err := regexp.Compile(pattern); err == nil {
			f.excludePatterns = append(f.excludePatterns, re)
		}
	}

	// Add custom excludes
	for _, pattern := range customExcludes {
		if re, err := regexp.Compile(pattern); err == nil {
			f.excludePatterns = append(f.excludePatterns, re)
		}
	}

	return f
}

// LoadGitignore parses .gitignore file and loads patterns
func (f *Filter) LoadGitignore(repoRoot string) error {
	if !f.respectGitignore {
		return nil
	}

	gitignorePath := filepath.Join(repoRoot, ".gitignore")
	file, err := os.Open(gitignorePath)
	if err != nil {
		// No .gitignore is fine
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Convert gitignore pattern to glob pattern
		pattern := line
		
		// Handle negation (we'll just skip these for simplicity)
		if strings.HasPrefix(pattern, "!") {
			continue
		}

		// Handle directory patterns
		if strings.HasSuffix(pattern, "/") {
			pattern = pattern + "**"
		}

		// Handle patterns starting with /
		if strings.HasPrefix(pattern, "/") {
			pattern = strings.TrimPrefix(pattern, "/")
		} else {
			// Patterns without / should match anywhere
			pattern = "**/" + pattern
		}

		// Compile glob pattern
		if g, err := glob.Compile(pattern, '/'); err == nil {
			f.gitignoreGlobs = append(f.gitignoreGlobs, g)
		}
	}

	return scanner.Err()
}

// ShouldInclude checks if a file should be included based on all filters
func (f *Filter) ShouldInclude(path string) bool {
	// Normalize path separators
	path = filepath.ToSlash(path)

	// Check always-excluded patterns first
	for _, re := range f.excludePatterns {
		if re.MatchString(path) {
			return false
		}
	}

	// Check gitignore patterns
	if f.respectGitignore {
		for _, g := range f.gitignoreGlobs {
			if g.Match(path) {
				return false
			}
		}
	}

	// Check if extension is allowed
	ext := filepath.Ext(path)
	if ext == "" || !f.allowedExts[ext] {
		return false
	}

	return true
}

// CountLines counts the number of lines in a file
func CountLines(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		count++
	}

	return count, scanner.Err()
}

