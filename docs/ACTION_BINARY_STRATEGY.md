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

### 1. `prepare-action-release.yml`

- **Trigger**: Release Please PRs (branches starting with `release-please--`)
- **Purpose**: Builds and commits optimized binary to release PRs
- **Features**:
  - Runs full test suite before building
  - Creates optimized binary with `-ldflags="-s -w"`
  - Force-adds binary (overrides `.gitignore`)
  - Comments on PR with binary info

### 2. `build-action.yml`

- **Trigger**: Git tags (`v*`) or manual dispatch
- **Purpose**: Builds and commits binary for tagged releases
- **Features**:
  - Production-ready binary builds
  - Automatic commit to release tags

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

- Automated binary management
- No manual steps for maintainers
- Binary always matches released code
- Clear release process

## File Structure

```
cmd/
├── app-git-copy          # Binary (ignored in dev, committed in releases)
├── main.go              # Application entry point
└── ...

.github/workflows/
├── prepare-action-release.yml  # Builds binary for release PRs
├── build-action.yml           # Builds binary for tagged releases
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
- uses: pal-paul/git-copy@v4.0.0
  with:
    # ... action inputs
```

### For Maintainers

1. Create release PR (Release Please handles this)
2. `prepare-action-release.yml` automatically builds and commits binary
3. Merge release PR
4. `build-action.yml` ensures binary is committed to release tag
5. Users can immediately use the new version

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

1. **Force Add Strategy**: Release workflows use `git add -f` to override `.gitignore`
2. **Conditional Workflows**: Only runs on release branches to avoid unnecessary builds
3. **Error Handling**: Action includes binary existence check with helpful error messages
4. **Optimization**: Release binaries are optimized for size and performance
