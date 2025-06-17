# Contributing to git-copy

Thank you for your interest in contributing to git-copy! This document provides guidelines and information for contributors.

## Table of Contents
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project mastertainers.

## Getting Started

### Prerequisites
- Go 1.21 or later
- Git
- Make
- GitHub account

### Development Tools (Optional but Recommended)
- [golangci-lint](https://golangci-lint.run/) for linting
- [actionlint](https://github.com/rhymond/actionlint) for GitHub Actions validation
- [govulncheck](https://golang.org/x/vuln/cmd/govulncheck) for security scanning

## Development Setup

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/YOUR-USERNAME/git-copy.git
   cd git-copy
   ```

2. **Install dependencies**
   ```bash
   go mod download
   make install
   ```

3. **Verify setup**
   ```bash
   make test-all
   make build
   ```

## Making Changes

### Branch Naming Convention
- `feature/your-feature-name` - for new features
- `bugfix/issue-description` - for bug fixes
- `docs/documentation-update` - for documentation changes
- `chore/mastertenance-task` - for mastertenance tasks

### Commit Message Guidelines
Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

Examples:
```
feat: add support for custom commit messages
fix: resolve race condition in file processing
docs: update README with new examples
test: add integration tests for directory operations
```

### Code Style
- Follow Go formatting conventions (`gofmt`)
- Use meaningful variable and function names
- Add comments for complex logic
- Ensure all exported functions have documentation

## Testing

### Running Tests
```bash
# Run all tests
make test-all

# Run specific test types
make test           # Unit tests
make test-race      # Race condition detection
make test-coverage  # Coverage analysis
make test-startup   # Application startup validation
```

### Writing Tests
- Place tests in the `test/` directory
- Use table-driven tests when appropriate
- Test both success and error cases
- Add benchmarks for performance-critical code

### Test Categories
1. **Unit Tests** (`test/cmd_test.go`) - Test individual functions
2. **Integration Tests** (`test/integration_test.go`) - Test workflows
3. **Git Operations Tests** (`test/git_operations_test.go`) - Test Git interactions
4. **Edge Cases Tests** (`test/edge_cases_test.go`) - Test unusual scenarios

## Code Quality

### Before Submitting
Run the following commands to ensure code quality:

```bash
# Format code
make fmt

# Run linting
make lint

# Run all tests
make test-all

# Validate CI configuration
make validate-ci

# Build application
make build
```

### Required Checks
- [ ] All tests pass
- [ ] Code is formatted (`gofmt`)
- [ ] Linting passes
- [ ] No race conditions detected
- [ ] Security scan passes
- [ ] Documentation updated

## Submitting Changes

### Pull Request Process

1. **Create a descriptive PR title**
   ```
   feat: add support for recursive directory exclusions
   fix: handle empty files correctly
   docs: improve troubleshooting section
   ```

2. **Fill out the PR template completely**
   - Describe your changes
   - Link related issues
   - Confirm all tests pass
   - Add screenshots if applicable

3. **Ensure CI passes**
   - All GitHub Actions workflows must pass
   - Address any failing checks

4. **Request review**
   - Tag appropriate reviewers
   - Respond to feedback promptly
   - Update your PR based on review comments

### Review Criteria
- Code quality and adherence to Go best practices
- Test coverage for new functionality
- Documentation updates
- Backward compatibility
- Performance impact
- Security considerations

## Release Process

### Versioning
We follow [Semantic Versioning](https://semver.org/):
- `MAJOR.MINOR.PATCH`
- Major: Breaking changes
- Minor: New features (backward compatible)
- Patch: Bug fixes (backward compatible)

### Release Steps
1. Update version in relevant files
2. Update CHANGELOG.md
3. Create release PR
4. Tag release after merge
5. GitHub Actions will handle the rest

## Project Structure

```
git-copy/
├── .github/              # GitHub templates and workflows
├── cmd/                  # Application entry point
├── internal/gitcopy/     # Core application logic
├── test/                 # Test files
├── action.yml           # GitHub Action metadata
├── Dockerfile           # Container configuration
├── Makefile             # Build and test commands
└── README.md            # Project documentation
```

## Getting Help

### Documentation
- Read the [README.md](README.md) for usage examples
- Check existing [issues](https://github.com/pal-paul/git-copy/issues)
- Review [pull requests](https://github.com/pal-paul/git-copy/pulls)

### Communication
- Create an issue for bugs or feature requests
- Use discussions for questions and general topics
- Tag mastertainers for urgent issues

### mastertainers
- @pal-paul - Project owner and primary mastertainer

## Recognition

Contributors will be recognized in:
- Release notes
- Contributors section
- GitHub contributor graph

Thank you for contributing to git-copy! 🎉
