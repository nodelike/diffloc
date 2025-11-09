# Contributing

## Setup

```bash
git clone https://github.com/YOUR_USERNAME/diffloc.git
cd diffloc
go mod download
make build
```

## Pull Requests

1. Fork and create a feature branch
2. Make changes, test them
3. Commit using conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`, `chore:`
4. Push and open PR

## Project Structure

```
cmd/diffloc/          # Main entry
internal/
├── analyzer/         # Git & file analysis
├── model/            # Data structures
└── ui/               # TUI
```

## Standards

- Use `gofmt`
- Keep functions small
- Test your changes (git repos, non-git dirs, edge cases)
- Be respectful

## Issues

Found a bug or have an idea? Open an issue.

