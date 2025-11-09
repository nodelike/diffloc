# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.4] - 2025-11-09

### Fixed
- **Viewport scrolling now works properly** using `charmbracelet/bubbles/viewport`
- Fixed terminal scrollback jumping issues with proper alt screen handling
- Fixed header cropping when scrolling to top (adjusted footer height)
- Fixed mouse escape sequences leaking into terminal prompt
- Full output now prints to terminal after exiting interactive mode

### Added
- Mouse wheel scrolling support in interactive mode
- Proper viewport integration for smooth scrolling experience

### Changed
- Replaced manual scroll logic with battle-tested viewport component
- Cleaner code (110 lines removed)

## [1.0.3] - 2025-11-09

### Changed
- **Defaults to bottom** on startup (summary visible first, scroll up for details)
- Stay at bottom when changing sort order
- Cleaned up CONTRIBUTING.md (192 → 40 lines, removed verbose content)

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

