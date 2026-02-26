# Action Binary Management Strategy

## Overview

This repository uses a two-phase binary management strategy to keep development clean while ensuring production reliability.

## Strategy

### 🛠️ Development Phase

- **Binary Location**: `cmd/app-git-copy`
- **Git Status**: **Excluded** (in `.gitignore`)
- **Purpose**: Local development and testing
- **Build Command**: `make build` or `go build -o ./cmd/app-git-copy ./cmd`

### 🚀 Release Phase

- **Binary Location**: `cmd/app-git-copy` (same location)
- **Git Status**: **Committed** (force-added during release)
- **Purpose**: Production use in GitHub Actions
- **Build**: Automated via release workflows

## Workflows

### 1. `release.yml`

- **Trigger**: Push to `master` (including merged PRs)
- **Purpose**: Prepares release artifacts in two steps:
  1. Builds Linux action binary and opens a dedicated PR containing only `cmd/app-git-copy`
  2. Creates (or updates) a GitHub draft release for manual publishing
- **Features**:
  - Runs test suite before building
  - Builds optimized binary with `-ldflags="-s -w"`
  - Opens an isolated binary-update PR automatically
  - Keeps release as draft so maintainers control final publish timing

## Benefits

### ✅ Development Benefits

- Clean development environment (no binary in git)
- Fast git operations (no large binary files)
- No merge conflicts on binary files
- Developers build fresh binaries locally

### ✅ Production Benefits

- Users get pre-built, optimized binaries
- No build requirements for action users
- Consistent binary across all action usage
- Smaller binary size with optimization flags

### ✅ Release Benefits

- Automated binary refresh after every merge to `master`
- Clear PR-first process for binary updates
- Draft-first release flow for manual control and validation
- Binary always produced in Linux runner-compatible format
- **No release-PR coupling required**

## File Structure

```bash
cmd/
├── app-git-copy          # Binary (ignored in dev, committed in releases)
├── cmd.go               # Application entry point
└── ...

.github/workflows/
├── release.yml                 # Builds binary PR + creates/updates release draft
├── test.yml                    # CI testing and validation
└── ...

.gitignore
├── cmd/app-git-copy      # Binary excluded during development
```

## Usage

### For Developers

```bash
# Build for local development
make build

# Run locally
./cmd/app-git-copy --help

# Clean up
make clean
```

### For Action Users

```yaml
# In GitHub Actions workflows
- uses: pal-paul/git-copy@v2.1.4
  with:
    # ... action inputs
```

### For Maintainers

1. Merge code into `master`
2. `release.yml` automatically opens a binary-update PR
3. `release.yml` creates or updates a draft release
4. Merge the binary PR
5. Publish the prepared draft release manually

## Troubleshooting

### "Action binary not found" Error

This error occurs when:

1. Using a development branch that doesn't have the binary committed
2. Using a version before the binary management strategy was implemented

**Solution**: Use a proper release tag (e.g., `@v4.0.0`) instead of `@main` or development branches.

### Local Development

If you need the binary for local development:

```bash
# Build fresh binary
make build

# Or with Go directly
go build -o ./cmd/app-git-copy ./cmd
```

### Binary Size Optimization

The release binary is built with optimization flags:

- `-ldflags="-s -w"`: Strips debug information and symbol tables
- Results in smaller binary size for distribution

## Implementation Notes

1. **Isolated Binary PR**: Release workflow keeps binary updates in a separate PR
2. **Draft-first Release**: Workflow creates a draft release instead of auto-publishing
3. **Error Handling**: Action includes binary existence check with helpful error messages
4. **Optimization**: Release binaries are optimized for size and performance
