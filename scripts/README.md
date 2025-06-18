# Developer Scripts

This directory contains scripts for local development and testing.

## Available Scripts

### `lint.sh`
Comprehensive linting script for local development that includes:
- Basic Go checks (`go vet`, `go fmt`)
- Advanced linting with golangci-lint (with version compatibility)
- Security checks for hardcoded credentials and unsafe operations
- Vulnerability scanning (if tools are available)

**Usage:**
```bash
# Run directly
./scripts/lint.sh

# Or use make target
make lint-local
```

**Requirements:**
- Go toolchain
- Optional: golangci-lint, govulncheck

**Installation of tools:**
```bash
# Install all development tools
make install-tools
```

## Why Separate from CI?

The advanced linting (golangci-lint) has been moved to local development only because:

1. **Version Compatibility**: Different golangci-lint versions have incompatible configurations
2. **Go Version Conflicts**: CI environments may have version mismatches
3. **Developer Flexibility**: Developers can install and use the latest tools locally
4. **CI Stability**: CI focuses on core checks (go vet, go fmt) that are stable across versions

The CI pipeline still runs basic but reliable checks:
- `go vet` for static analysis
- `go fmt` for code formatting
- Security scans
- All unit and integration tests
