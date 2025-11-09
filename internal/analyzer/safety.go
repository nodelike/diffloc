package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidatePath checks if the given path is safe to analyze
// Returns an error if the path is considered dangerous
func ValidatePath(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	cleanPath := filepath.Clean(absPath)

	if cleanPath == "/" {
		return fmt.Errorf("refusing to run in root directory (/). This could scan your entire system")
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		cleanHome := filepath.Clean(homeDir)

		if cleanPath == cleanHome {
			return fmt.Errorf("refusing to run in home directory (%s). Please run in a specific project directory", cleanHome)
		}
	}

	dangerousDirs := []string{
		"/usr",
		"/etc",
		"/var",
		"/bin",
		"/sbin",
		"/boot",
		"/sys",
		"/proc",
		"/dev",
		"/System",       // macOS
		"/Library",      // macOS
		"/Applications", // macOS
		"/Volumes",      // macOS
		"/private",      // macOS
		"/opt",
		"/root",
		"/tmp",
		"/Windows",             // Windows
		"/Program Files",       // Windows
		"/Program Files (x86)", // Windows
	}

	for _, dangerousDir := range dangerousDirs {
		cleanDangerous := filepath.Clean(dangerousDir)
		if cleanPath == cleanDangerous {
			return fmt.Errorf("refusing to run in system directory (%s). This is a protected system location", cleanPath)
		}
	}

	pathDepth := strings.Count(cleanPath, string(os.PathSeparator))
	if pathDepth <= 2 && cleanPath != "/" {
		parentDir := filepath.Dir(cleanPath)
		if parentDir == "/" || (homeDir != "" && parentDir == filepath.Dir(homeDir)) {
			return fmt.Errorf("refusing to run at (%s). This directory is too broad. Please run in a specific project directory", cleanPath)
		}
	}

	return nil
}

// ShouldWarnLargeDirectory checks if we should warn the user about analyzing a large directory
func ShouldWarnLargeDirectory(path string) (bool, string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, ""
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, ""
	}

	if strings.HasPrefix(absPath, homeDir) {
		relPath, err := filepath.Rel(homeDir, absPath)
		if err == nil && !strings.Contains(relPath, string(os.PathSeparator)) {
			return true, fmt.Sprintf("Warning: Analyzing a top-level directory in your home folder (%s). This may take a while.", filepath.Base(absPath))
		}
	}

	return false, ""
}
