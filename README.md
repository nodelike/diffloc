# diffloc

Git diff & line count statistics with a clean TUI. Shows line counts, additions, and deletions for your codebase. Works in non-git repos too. Currently supports Golang, JavaScript, and Python projects, but you can extend it to your needs with a `.diffloc.yaml` config file in your project folder. See [instructions below](#config-file).

## Installation

```bash
# Homebrew
brew install nodelike/tap/diffloc

# Go
go install github.com/nodelike/diffloc/cmd/diffloc@latest
```

## Usage

```bash
diffloc                    # Current directory
diffloc /path/to/project   # Specific path
diffloc --json             # JSON output
diffloc --static           # Non-interactive output
```

### Keyboard Controls

- `n` - Sort by name
- `l` - Sort by lines  
- `a` - Sort by additions
- `d` - Sort by deletions
- `q` - Quit

### Flags

| Flag | Description |
|------|-------------|
| `--no-gitignore` | Ignore .gitignore patterns |
| `--exclude-tests` | Exclude test files |
| `--exclude <pattern>` | Custom exclusion regex (repeatable) |
| `--ext <ext>` | Override allowed extensions (repeatable) |
| `--max-depth <n>` | Limit directory depth (0 = unlimited) |
| `--json` | Output as JSON |
| `--static` | Non-interactive output |
| `--config <file>` | Config file path (default: `.diffloc.yaml`) |

### Config File

Create `.diffloc.yaml` in your project root:

```yaml
exclude-tests: true
no-gitignore: false
max-depth: 0
exclude:
  - "vendor/"
  - "\\.pb\\.go$"
ext:
  - ".go"
  - ".py"
```

## What Gets Excluded

**Directories:** `node_modules`, `venv`, `.venv`, `__pycache__`, `.git`, `dist`, `build`, `.egg-info`, `.tox`, `coverage`, `.next`, `vendor`, `bin`, `tmp`

**Files:** Lock files (`*.lock`, `*-lock.json`), binaries (`*.exe`, `*.so`, `*.dylib`, `*.dll`), images (`*.jpg`, `*.png`, `*.svg`, etc.), generated files (`*_templ.go`, `*.pb.go`, `*.min.js`)

**Default Extensions:** `.go`, `.py`, `.js`, `.jsx`, `.ts`, `.tsx`, `.vue`, `.svelte`, `.mjs`, `.cjs`

## Features

- Works in git repos and non-git directories
- Respects .gitignore (optional)
- Smart filtering for common artifacts
- Interactive sorting
- JSON and static output modes
- Configurable via file or flags

## License

MIT


made with ❤️ by [@nodelike](https://nodelike.com/)
