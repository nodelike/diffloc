# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.2] - 2025-11-09

### Added
- **Full viewport scrolling** in interactive TUI mode
  - Scroll through entire content including file lists, changed/unchanged sections, and summary
  - Arrow keys (↑/↓) or Vim keys (j/k) for line-by-line scrolling
  - Page Up/Page Down for faster navigation
  - Home/g to jump to top, End/G to jump to bottom
  - Scroll position indicator showing current line and total lines
  - Visual indicators when more content is available above/below

### Changed
- Improved keyboard controls section in README with dedicated scrolling section

## [1.0.1] - 2025-11-09

### Added
- Safety checks to prevent running diffloc in dangerous locations
  - Blocks execution in root directory (`/`)
  - Blocks execution in home directory (e.g., `~` or `/Users/username`)
  - Blocks execution in system directories (`/usr`, `/etc`, `/var`, `/System`, `/Library`, etc.)
  - Blocks execution in overly broad directories (e.g., `/Users`, `/home`)
  - Warning for top-level directories in home folder (e.g., `~/Documents`, `~/Desktop`)

### Security
- Prevents accidental scanning of entire filesystem or sensitive system directories
- Helps protect users from unintentionally triggering expensive operations on large directory trees

## [1.0.0] - 2025-11-09

### Initial Release
- Git diff and line count statistics with clean TUI
- Support for both Git repositories and non-Git directories
- Respects .gitignore patterns (optional)
- Smart filtering for common artifacts
- Interactive sorting by name, lines, additions, deletions
- JSON and static output modes
- Configurable via `.diffloc.yaml` file or command-line flags
- Support for Golang, JavaScript, TypeScript, Python, and Vue/Svelte projects
- Maximum depth limiting for directory traversal
- Customizable file extension and exclusion patterns

