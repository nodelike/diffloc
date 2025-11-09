# diffloc Quick Start Guide

This guide will get you up and running with diffloc in minutes.

## Build and Run

```bash
# Build the binary
make build

# Run diffloc in current directory
./bin/diffloc

# Or use make run
make run
```

## Try It Out

### In a Git Repository

```bash
cd /path/to/your/git/project
diffloc
```

You'll see:
- Changed files with additions/deletions
- Unchanged tracked files
- Total statistics
- Interactive sorting controls

### In a Non-Git Directory

```bash
cd /path/to/any/directory
diffloc
```

You'll see:
- All files with line counts
- No diff statistics (since there's no git)

## Keyboard Controls While Running

- Press `n` → Sort by name
- Press `l` → Sort by lines
- Press `a` → Sort by additions
- Press `d` → Sort by deletions
- Press `q` → Quit

## Common Usage Patterns

### Analyze a Specific Project

```bash
diffloc -path ~/projects/myapp
```

### Ignore .gitignore Patterns

Useful when you want to see files that are normally gitignored:

```bash
diffloc -no-gitignore
```

### Only Show Specific File Types

```bash
# Only Go files
diffloc -ext .go -ext .mod

# Only Python files
diffloc -ext .py
```

### Add Custom Exclusions

```bash
# Exclude test files
diffloc -exclude ".*_test\.go" -exclude "test_.*"
```

### Combine Options

```bash
diffloc -path ~/projects/myapp -no-gitignore -ext .go -ext .py
```

## What Gets Counted?

By default, diffloc counts these file types:
- `.go` (Go)
- `.py` (Python)
- `.js`, `.jsx`, `.mjs`, `.cjs` (JavaScript)
- `.ts`, `.tsx` (TypeScript)
- `.vue` (Vue.js)
- `.svelte` (Svelte)

And automatically excludes:
- `node_modules`, `venv`, `.git`, etc.
- Lock files
- Binary files
- Images
- Build artifacts

See README.md for the complete list.

## Next Steps

### Install System-Wide

```bash
make install
# Now you can use 'diffloc' from anywhere
```

### Create Your First Release

See [HOMEBREW_TAP_SETUP.md](HOMEBREW_TAP_SETUP.md) for detailed instructions on:
1. Creating a GitHub release
2. Setting up Homebrew tap
3. Publishing via GoReleaser

### Contribute

See [CONTRIBUTING.md](CONTRIBUTING.md) for:
- Development setup
- Coding standards
- How to submit PRs

## Troubleshooting

### "reference not found" Error

This happens in a new git repository with no commits. Make an initial commit:

```bash
git add -A
git commit -m "Initial commit"
```

### No Files Shown

Check if:
1. Files match the allowed extensions (use `-ext` to customize)
2. Files aren't in excluded directories
3. Files aren't in .gitignore (use `-no-gitignore` to override)

### "could not open a new TTY" Error

This happens in non-interactive environments (CI, pipes, etc.). The tool requires a terminal to display the TUI.

## Examples Output

### Changed Files

```
Changed Files:
  Lines    +Add     -Del     File
  ────────────────────────────────────────
  256      +45      -12      cmd/diffloc/main.go
  189      +23      -0       internal/ui/tui.go
  42       +42      -0       internal/model/types.go
```

### Summary

```
╭──────────────────────────────────╮
│ Net Change: +98 (increased)      │
│                                  │
│ Total Files: 23                  │
│  (3 changed, 20 unchanged)       │
│ Total Lines: 2,341               │
│ Insertions: +110                 │
│ Deletions: -12                   │
╰──────────────────────────────────╯
```

## More Information

- Full documentation: [README.md](README.md)
- Homebrew setup: [HOMEBREW_TAP_SETUP.md](HOMEBREW_TAP_SETUP.md)
- Contributing: [CONTRIBUTING.md](CONTRIBUTING.md)

