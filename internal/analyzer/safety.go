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
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Clean the path
	cleanPath := filepath.Clean(absPath)

	// Check if running in root directory
	if cleanPath == "/" {
		return fmt.Errorf("refusing to run in root directory (/). This could scan your entire system")
	}

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		cleanHome := filepath.Clean(homeDir)

		// Check if running directly in home directory
		if cleanPath == cleanHome {
			return fmt.Errorf("refusing to run in home directory (%s). Please run in a specific project directory", cleanHome)
		}
	}

	// Check for dangerous system directories
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
		// Check if path is the dangerous directory or a direct subdirectory
		if cleanPath == cleanDangerous {
			return fmt.Errorf("refusing to run in system directory (%s). This is a protected system location", cleanPath)
		}
	}

	// Additional check: warn if path looks like it might be too broad
	// Count the depth from root - if it's too shallow, it might be dangerous
	pathDepth := strings.Count(cleanPath, string(os.PathSeparator))
	if pathDepth <= 2 && cleanPath != "/" {
		// This catches things like /Users or /home which have depth of 2
		// But we want to allow things like /Users/username/project (depth 3+)
		parentDir := filepath.Dir(cleanPath)
		if parentDir == "/" || (homeDir != "" && parentDir == filepath.Dir(homeDir)) {
			return fmt.Errorf("refusing to run at (%s). This directory is too broad. Please run in a specific project directory", cleanPath)
		}
	}

	return nil
}

// ShouldWarnLargeDirectory checks if we should warn the user about analyzing a large directory
func ShouldWarnLargeDirectory(path string) (bool, string) {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, ""
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, ""
	}

	// Warn if analyzing the entire home directory (even subdirectories with lots of content)
	if strings.HasPrefix(absPath, homeDir) {
		// If it's directly under home (like ~/Documents, ~/Desktop)
		relPath, err := filepath.Rel(homeDir, absPath)
		if err == nil && !strings.Contains(relPath, string(os.PathSeparator)) {
			return true, fmt.Sprintf("Warning: Analyzing a top-level directory in your home folder (%s). This may take a while.", filepath.Base(absPath))
		}
	}

	return false, ""
}
