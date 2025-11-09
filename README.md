# diffloc ✨

A beautiful TUI (Terminal User Interface) tool that displays git diff statistics with line counts, additions, and deletions. Perfect for quickly understanding what changed in your codebase.

## Features

- 🎨 **Beautiful TUI**: Clean, colorful interface with Bubble Tea
- 📊 **Detailed Statistics**: View line counts, additions, and deletions for each file
- 🔄 **Interactive Sorting**: Sort files by name, lines, additions, or deletions with keyboard shortcuts
- 🚫 **Smart Filtering**: Automatically excludes common artifacts (node_modules, build dirs, lock files, images, binaries)
- 📁 **.gitignore Support**: Respects .gitignore patterns (optional)
- 🔧 **Configurable**: Override file extensions and add custom exclusions
- 🌍 **Non-Git Support**: Works in non-git directories, just shows file line counts
- 🎯 **Project Focus**: Optimized for React/NodeJS/Go/Python projects

## Installation

### Homebrew (macOS/Linux)

```bash
brew install nodelike/tap/diffloc
```

### From Source

Requires Go 1.24 or later:

```bash
go install github.com/nodelike/diffloc/cmd/diffloc@latest
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/nodelike/diffloc/releases).

## Usage

### Basic Usage

Run in current directory:

```bash
diffloc
```

Analyze a specific directory:

```bash
diffloc -path /path/to/project
```

### Keyboard Controls

- `n` - Sort by name
- `l` - Sort by lines
- `a` - Sort by additions (+)
- `d` - Sort by deletions (-)
- `q` or `Ctrl+C` - Quit

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-path` | Path to analyze | `.` (current directory) |
| `-no-gitignore` | Ignore .gitignore patterns (always-excluded patterns still apply) | `false` |
| `-exclude` | Additional exclusion pattern (can be repeated) | - |
| `-ext` | Override allowed file extensions (can be repeated) | `.go, .py, .js, .jsx, .ts, .tsx, .vue, .svelte, .mjs, .cjs` |

### Examples

Ignore .gitignore patterns:

```bash
diffloc -no-gitignore
```

Add custom exclusion patterns:

```bash
diffloc -exclude "test_.*\.go" -exclude ".*_test\.go"
```

Override allowed extensions:

```bash
diffloc -ext .go -ext .mod
```

Analyze a specific project:

```bash
diffloc -path ~/projects/myapp
```

## What Gets Excluded

### Always Excluded (Regardless of Flags)

**Directories:**
- `node_modules`, `venv`, `.venv`, `__pycache__`, `.git`
- `dist`, `build`, `.egg-info`, `.tox`, `coverage`
- `.next`, `vendor`, `bin`, `tmp`

**Lock Files:**
- `*.lock`, `*-lock.json`, `*-lock.yaml`
- `Pipfile.lock`, `.gitignore`

**Binaries:**
- `*.exe`, `*.so`, `*.dylib`, `*.dll`
- `*_templ.go` (templ generated files)

**Images:**
- `*.jpg`, `*.jpeg`, `*.png`, `*.gif`, `*.bmp`
- `*.svg`, `*.ico`, `*.webp`, `*.tiff`, `*.tif`
- `*.psd`, `*.raw`, `*.heic`, `*.avif`

### Allowed File Extensions (Default)

- `.go` - Go
- `.py` - Python
- `.js`, `.jsx`, `.mjs`, `.cjs` - JavaScript
- `.ts`, `.tsx` - TypeScript
- `.vue` - Vue.js
- `.svelte` - Svelte

You can override these with the `-ext` flag.

## Development

### Prerequisites

- Go 1.24 or later
- Make (optional, for convenience)

### Building

```bash
# Using Make
make build

# Or directly with Go
go build -o bin/diffloc ./cmd/diffloc
```

### Running

```bash
# Using Make
make run

# Or directly
go run ./cmd/diffloc
```

### Testing

```bash
make test
```

### Local Release Test

Test the GoReleaser configuration locally:

```bash
make release-test
```

## How It Works

### Git Mode

When run in a git repository, `diffloc`:

1. Compares the working directory against HEAD
2. Identifies changed files (modified tracked files)
3. Identifies untracked files (new files)
4. Lists unchanged tracked files
5. Calculates line additions/deletions using git diff
6. Displays everything in a beautiful TUI

### Non-Git Mode

When run outside a git repository:

1. Recursively scans all files in the directory
2. Counts lines for each file
3. Displays files as "unchanged" with 0 additions/deletions

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Roadmap

- [ ] Submit to homebrew-core (requires 75+ stars, 30+ days old)
- [ ] Add more output formats (JSON, plain text)
- [ ] Support comparing between specific commits/branches
- [ ] Add configuration file support (.difflocrc)
- [ ] Add file size filtering
- [ ] Add search/filter functionality in TUI
- [ ] Add pagination for large file lists
- [ ] Support for more language ecosystems

## License

MIT License - See [LICENSE](LICENSE) for details.

## Credits

Built with:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Styling
- [go-git](https://github.com/go-git/go-git) - Git operations
- [glob](https://github.com/gobwas/glob) - Pattern matching

## Support

- 🐛 [Report a bug](https://github.com/nodelike/diffloc/issues/new)
- 💡 [Request a feature](https://github.com/nodelike/diffloc/issues/new)
- 📖 [Read the docs](https://github.com/nodelike/diffloc)

---

Made with ❤️ by [@nodelike](https://github.com/nodelike)

