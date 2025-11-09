# Contributing to diffloc

Thank you for your interest in contributing to diffloc! This document provides guidelines and instructions for contributing.

## Code of Conduct

Be respectful and constructive. We're all here to build something useful.

## How to Contribute

### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/nodelike/diffloc/issues)
2. If not, create a new issue with:
   - Clear title and description
   - Steps to reproduce
   - Expected vs actual behavior
   - Your environment (OS, Go version)
   - Screenshots if applicable

### Suggesting Features

1. Check if the feature has already been requested
2. Create a new issue with:
   - Clear description of the feature
   - Use cases and benefits
   - Potential implementation approach (optional)

### Pull Requests

1. Fork the repository
2. Create a feature branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. Make your changes following our coding standards

4. Test your changes:
   ```bash
   make test
   make build
   ./bin/diffloc  # Test manually
   ```

5. Commit with clear messages:
   ```bash
   git commit -m "feat: add sorting by file size"
   ```

6. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

7. Create a Pull Request with:
   - Clear description of changes
   - Link to related issues
   - Screenshots/demos if applicable

## Development Setup

### Prerequisites

- Go 1.24 or later
- Git
- Make (optional)

### Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/diffloc.git
cd diffloc

# Install dependencies
go mod download

# Build
make build

# Run
make run
```

### Project Structure

```
diffloc/
├── cmd/diffloc/          # Main entry point
├── internal/
│   ├── analyzer/         # Git and file analysis logic
│   ├── model/            # Data structures
│   └── ui/               # TUI implementation
├── .goreleaser.yml       # Release configuration
└── .github/workflows/    # CI/CD
```

## Coding Standards

### Go Code Style

- Follow standard Go conventions (use `gofmt`)
- Write clear, self-documenting code
- Add comments for complex logic
- Keep functions small and focused

### Commit Messages

Use conventional commits format:

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting, etc.)
- `refactor:` Code refactoring
- `test:` Adding or updating tests
- `chore:` Maintenance tasks

Examples:
```
feat: add file size filtering
fix: handle empty git repositories
docs: update installation instructions
refactor: simplify diff calculation logic
```

### Testing

- Add tests for new features
- Ensure existing tests pass
- Test manually in various scenarios:
  - Git repositories
  - Non-git directories
  - Different project types (Go, Python, JS)
  - Edge cases (empty repos, large repos)

## Areas for Contribution

### Good First Issues

Look for issues labeled `good first issue` - these are great starting points!

### Priority Areas

1. **Testing**: More comprehensive test coverage
2. **Performance**: Optimize for large repositories
3. **Features**: See the Roadmap in README.md
4. **Documentation**: Improve docs, add examples
5. **Bug Fixes**: Check open issues

### Feature Ideas

Some ideas if you're looking for something to work on:

- JSON/plain text output formats
- Configuration file support
- File size filtering
- Compare between specific commits
- Search/filter in TUI
- Pagination for large file lists
- Support for more languages
- Integration with other tools

## Testing Releases Locally

Before submitting changes that affect releases:

```bash
# Test GoReleaser locally
make release-test

# Check the output in dist/
ls -la dist/
```

## Questions?

- Open a Discussion on GitHub
- Tag your issue with `question`
- Reach out to maintainers

## Recognition

Contributors will be recognized in:
- GitHub contributors page
- Release notes for significant contributions
- A future CONTRIBUTORS.md file

Thank you for contributing! 🎉

