# git-copy

[![Test and Build](https://github.com/pal-paul/git-copy/actions/workflows/test.yml/badge.svg)](https://github.com/pal-paul/git-copy/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/pal-paul/git-copy)](https://goreportcard.com/report/github.com/pal-paul/git-copy)
[![codecov](https://codecov.io/gh/pal-paul/git-copy/branch/master/graph/badge.svg)](https://codecov.io/gh/pal-paul/git-copy)

A GitHub Action for copying files and directories between repositories with batch operations and pull request automation.

## Features

- **Batch File Operations**: Copy multiple files or entire directories in a single commit
- **Pull Request Automation**: Automatically create pull requests with reviewers
- **Cross-platform Support**: Works on Linux, macOS, and Windows
- **Error Resilience**: Continues processing even when some files fail
- **Comprehensive Logging**: Detailed logging for debugging and monitoring

## Usage

### Basic File Copy

Copy a single file from one repository to another:

```yaml
name: Copy Configuration File
on:
  push:
    branches: [master]
    paths: ["config/**"]

jobs:
  copy-config:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout source repo
        uses: actions/checkout@v4

      - name: Copy config file to destination repo
        uses: pal-paul/git-copy@v2.1.4
        with:
          owner: "your-org"
          repo: "destination-repo"
          token: "${{ secrets.GITHUB_TOKEN }}"
          file_path: "config/app.json"
          destination_file_path: "configs/production.json"
          branch: "config-update-${{ github.sha }}"
          pull_message: "🔧 Update production config from master repo"
          pull_description: |
            Automated config sync from ${{ github.repository }}
            
            **Changes:**
            - Updated production configuration
            - Synced from commit: ${{ github.sha }}
            - Branch: ${{ github.ref_name }}
          reviewers: "devops-team,config-maintainer"
```

### Directory Synchronization

Copy an entire directory structure with comprehensive workflow:

```yaml
name: Sync Documentation
on:
  push:
    branches: [master]
    paths: ["docs/**"]
  schedule:
    # Sync docs daily at 2 AM UTC
    - cron: "0 2 * * *"

jobs:
  sync-docs:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout source repo
        uses: actions/checkout@v4

      - name: Copy documentation to public repo
        uses: pal-paul/git-copy@v2.1.4
        with:
          owner: "your-org"
          repo: "public-docs"
          token: "${{ secrets.GITHUB_TOKEN }}"
          directory: "docs/"
          destination_directory: "api-docs/"
          branch: "docs-sync-${{ github.run_number }}"
          pull_message: "📚 Sync API documentation"
          pull_description: |
            🔄 **Automated Documentation Sync**
            
            **Source:** ${{ github.repository }}
            **Commit:** ${{ github.sha }}
            **Triggered by:** ${{ github.event_name }}
            
            **Updated Documentation:**
            - API specifications
            - Code examples  
            - Integration guides
            - Troubleshooting sections
            
            **Review Required:** Documentation team review needed before merge.
          team_reviewers: "docs-team"
          reviewers: "tech-writer-lead"
```

### Cross-Organization Copy

Copy files between different GitHub organizations:

```yaml
name: Sync Shared Libraries
on:
  release:
    types: [published]

jobs:
  sync-to-partner-org:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout source repo
        uses: actions/checkout@v4

      - name: Copy shared libraries to partner organization
        uses: pal-paul/git-copy@v2.1.4
        with:
          owner: "partner-org"
          repo: "shared-components"
          token: "${{ secrets.PARTNER_ORG_TOKEN }}"  # Important: Use cross-org token
          directory: "lib/shared/"
          destination_directory: "vendor/your-org-libs/"
          branch: "update-shared-libs-${{ github.event.release.tag_name }}"
          pull_message: "📦 Update shared libraries to ${{ github.event.release.tag_name }}"
          pull_description: |
            🚀 **Shared Library Update**
            
            **Release:** ${{ github.event.release.name }}
            **Tag:** ${{ github.event.release.tag_name }}
            **Release Notes:** ${{ github.event.release.html_url }}
            
            **Updated Components:**
            - Core utilities
            - Authentication helpers
            - Data processing modules
            
            **Breaking Changes:** See release notes for migration guide.
          reviewers: "integration-team"
          team_reviewers: "shared-lib-maintainers"
```

### Multi-Environment Deployment

Deploy different configurations to multiple environments with conditional logic:

```yaml
name: Deploy Settings to Multiple Environments
on:
  push:
    branches: [master]
    paths: ["settings/**"]
  workflow_dispatch:
    inputs:
      environment:
        description: 'Target environment'
        required: true
        default: 'dev'
        type: choice
        options:
        - dev
        - staging
        - production

jobs:
  deploy-dev:
    if: github.event_name == 'push' || github.event.inputs.environment == 'dev'
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repo
        uses: actions/checkout@v4

      - name: Deploy to Dev Environment
        uses: pal-paul/git-copy@v2.1.4
        with:
          owner: "your-org"
          repo: "infrastructure-configs"
          token: "${{ secrets.DEPLOY_TOKEN }}"
          directory: "settings/dev/"
          destination_directory: "environments/dev/configs/"
          branch: "dev-config-update-${{ github.run_number }}"
          pull_message: "🚀 Update dev environment settings"
          pull_description: |
            **Development Environment Update**
            
            **Source:** ${{ github.repository }}
            **Environment:** Development
            **Auto-merge:** Safe for dev environment
          reviewers: "dev-team"

  deploy-staging:
    if: github.event.inputs.environment == 'staging'
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repo
        uses: actions/checkout@v4

      - name: Deploy to Staging Environment
        uses: pal-paul/git-copy@v2.1.4
        with:
          owner: "your-org"
          repo: "infrastructure-configs"
          token: "${{ secrets.DEPLOY_TOKEN }}"
          directory: "settings/staging/"
          destination_directory: "environments/staging/configs/"
          branch: "staging-config-update-${{ github.run_number }}"
          pull_message: "🧪 Update staging environment settings"
          pull_description: |
            **Staging Environment Update**
            
            **Environment:** Staging
            **Approval Required:** DevOps team approval needed
            **Testing:** Requires validation before production
          reviewers: "staging-team,devops-lead"

  deploy-production:
    if: github.event.inputs.environment == 'production'
    runs-on: ubuntu-latest
    environment: production  # Requires approval
    steps:
      - name: Checkout repo
        uses: actions/checkout@v4

      - name: Deploy to Production Environment
        uses: pal-paul/git-copy@v2.1.4
        with:
          owner: "your-org"
          repo: "infrastructure-configs"
          token: "${{ secrets.PRODUCTION_DEPLOY_TOKEN }}"
          directory: "settings/production/"
          destination_directory: "environments/production/configs/"
          branch: "production-config-update-${{ github.run_number }}"
          pull_message: "🔥 PRODUCTION: Update production environment settings"
          pull_description: |
            🚨 **PRODUCTION ENVIRONMENT UPDATE**
            
            **Environment:** Production
            **Approval Required:** Manager approval mandatory
            **Impact:** Live system configuration changes
            **Rollback Plan:** Previous config available in git history
            
            **Checklist:**
            - [ ] Configuration validated in staging
            - [ ] Backup plan confirmed
            - [ ] Monitoring alerts configured
            - [ ] Team notified of deployment
          reviewers: "production-team"
          team_reviewers: "devops-managers,security-team"
```

### Conditional File Copy

Copy files based on specific conditions and path filters:

```yaml
name: Conditional Config Sync
on:
  push:
    branches: [master]

jobs:
  sync-configs:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        config:
          - { file: "database.yml", env: "production", reviewers: "dba-team" }
          - { file: "api.yml", env: "staging", reviewers: "backend-team" }
          - { file: "frontend.yml", env: "development", reviewers: "frontend-team" }
    steps:
      - name: Checkout repo
        uses: actions/checkout@v4

      - name: Check if config file changed
        id: changes
        run: |
          if git diff --name-only ${{ github.event.before }} ${{ github.sha }} | grep -q "configs/${{ matrix.config.file }}"; then
            echo "changed=true" >> $GITHUB_OUTPUT
          else
            echo "changed=false" >> $GITHUB_OUTPUT
          fi

      - name: Copy changed config to environment repo
        if: steps.changes.outputs.changed == 'true'
        uses: pal-paul/git-copy@v2.1.4
        with:
          owner: "your-org"
          repo: "${{ matrix.config.env }}-configs"
          token: "${{ secrets.CONFIG_SYNC_TOKEN }}"
          file_path: "configs/${{ matrix.config.file }}"
          destination_file_path: "app/${{ matrix.config.file }}"
          branch: "update-${{ matrix.config.file }}-${{ github.sha }}"
          pull_message: "🔧 Update ${{ matrix.config.file }} for ${{ matrix.config.env }}"
          reviewers: "${{ matrix.config.reviewers }}"
```

### Cross-Organization Copy

Copy files between different GitHub organizations:

```yaml
name: Sync Shared Libraries
on:
  release:
    types: [published]

jobs:
  sync-to-partner-org:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout source repo
        uses: actions/checkout@v4

      - name: Copy shared libraries to partner organization
        uses: pal-paul/git-copy@v2.1.4
        with:
          owner: "partner-org"
          repo: "shared-components"
          token: "${{ secrets.PARTNER_ORG_TOKEN }}"
          directory: "lib/shared/"
          destination_directory: "vendor/your-org-libs/"
          branch: "update-shared-libs-${{ github.event.release.tag_name }}"
          pull_message: "📦 Update shared libraries to ${{ github.event.release.tag_name }}"
          pull_description: |
            Updated shared libraries from ${{ github.repository }}
            
            Release: ${{ github.event.release.name }}
            Tag: ${{ github.event.release.tag_name }}
            Release Notes: ${{ github.event.release.html_url }}
          reviewers: "integration-team"
```

### Conditional Copy with Path Filtering

Copy different files based on what changed:

```yaml
name: Smart Configuration Sync
on:
  push:
    branches: [master]

jobs:
  detect-changes:
    runs-on: ubuntu-latest
    outputs:
      database-changed: ${{ steps.changes.outputs.database }}
      api-changed: ${{ steps.changes.outputs.api }}
      frontend-changed: ${{ steps.changes.outputs.frontend }}
    steps:
      - uses: actions/checkout@v4
      - uses: dorny/paths-filter@v2
        id: changes
        with:
          filters: |
            database:
              - 'config/database/**'
            api:
              - 'config/api/**'
            frontend:
              - 'config/frontend/**'

  sync-database-config:
    needs: detect-changes
    if: needs.detect-changes.outputs.database-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repo
        uses: actions/checkout@v4

      - name: Copy database config
        uses: pal-paul/git-copy@v2.1.4
        with:
          owner: "your-org"
          repo: "database-infrastructure"
          token: "${{ secrets.GITHUB_TOKEN }}"
          directory: "config/database/"
          destination_directory: "configs/"
          branch: "database-config-sync"
          pull_message: "🗄️ Update database configuration"

  sync-api-config:
    needs: detect-changes
    if: needs.detect-changes.outputs.api-changed == 'true'
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repo
        uses: actions/checkout@v4

      - name: Copy API config
        uses: pal-paul/git-copy@v2.1.4
        with:
          owner: "your-org"
          repo: "api-gateway"
          token: "${{ secrets.GITHUB_TOKEN }}"
          directory: "config/api/"
          destination_directory: "gateway-config/"
          branch: "api-config-sync"
          pull_message: "🔌 Update API gateway configuration"
```

## Input Parameters

### Required Parameters

| Parameter | Description | Example |
|-----------|-------------|---------|
| `owner` | Target repository owner/organization | `"your-org"` |
| `repo` | Target repository name | `"destination-repo"` |
| `token` | GitHub token with repository access | `"${{ secrets.GITHUB_TOKEN }}"` |

### File Operation Parameters (choose one)

| Parameter | Description | Example |
|-----------|-------------|---------|
| `file_path` | Source file path to copy | `"config/app.json"` |
| `destination_file_path` | Target file path (required with `file_path`) | `"configs/production.json"` |
| `directory` | Source directory to copy | `"docs/"` |
| `destination_directory` | Target directory (required with `directory`) | `"api-docs/"` |

### Optional Parameters

| Parameter | Description | Default | Example |
|-----------|-------------|---------|---------|
| `branch` | Target branch name | `"update-branch"` | `"feature/config-update"` |
| `ref_branch` | Source branch to branch from | `"master"` | `"master"` |
| `pull_message` | Pull request title | Auto-generated | `"Update configuration"` |
| `pull_description` | Pull request description | Auto-generated | `"Automated sync from master repo"` |
| `reviewers` | Comma-separated list of reviewers | None | `"user1,user2,user3"` |
| `team_reviewers` | Comma-separated list of team reviewers | None | `"team1,team2"` |

### Example with All Parameters

```yaml
- name: Complete example
  uses: pal-paul/git-copy@v2.1.4
  with:
    # Required
    owner: "your-org"
    repo: "destination-repo"
    token: "${{ secrets.GITHUB_TOKEN }}"
    
    # File operation (choose file OR directory)
    file_path: "config/app.json"
    destination_file_path: "configs/production.json"
    
    # Optional
    branch: "config-update-${{ github.sha }}"
    ref_branch: "master"
    pull_message: "🚀 Update production configuration"
    pull_description: |
      Automated configuration update from ${{ github.repository }}
      
      Changes include:
      - Updated API endpoints
      - New feature flags
      - Performance optimizations
    reviewers: "devops-lead,config-admin"
    team_reviewers: "platform-team,security-team"
```

## Parameters

All input parameters for the git-copy action:

| Parameter | Description | Required | Default | Example |
|-----------|-------------|----------|---------|---------|
| `owner` | GitHub owner/organization name of destination repo | ✅ Yes | - | `"your-org"` |
| `repo` | GitHub repository name of destination repo | ✅ Yes | - | `"target-repo"` |
| `token` | GitHub token with repo access | ✅ Yes | - | `"${{ secrets.GITHUB_TOKEN }}"` |
| `ref_branch` | Base branch of destination repo | ❌ No | `master` | `"main"`, `"develop"` |
| `branch` | Branch name for the pull request | ❌ No | Auto-generated | `"config-update-123"` |
| `file_path` | Path to source file (for single file copy) | ❌ No* | - | `"config/app.json"` |
| `destination_file_path` | Destination path for the file | ❌ No* | Same as source | `"configs/production.json"` |
| `directory` | Path to source directory (for directory copy) | ❌ No* | - | `"docs/"` |
| `destination_directory` | Destination path for directory | ❌ No* | Same as source | `"public-docs/"` |
| `pull_message` | Pull request title | ❌ No | Auto-generated | `"Update configuration"` |
| `pull_description` | Pull request description | ❌ No | Auto-generated | `"Automated sync from master repo"` |
| `reviewers` | Comma-separated list of reviewers | ❌ No | None | `"user1,user2,user3"` |
| `team_reviewers` | Comma-separated list of team reviewers | ❌ No | None | `"team1,team2"` |

**\* Note:** Either `file_path` OR `directory` must be provided (not both).

### Parameter Examples

#### File Copy Parameters

```yaml
# Copy single file with rename
file_path: "src/config.yaml"
destination_file_path: "config/production.yaml"
```

#### Directory Copy Parameters

```yaml
# Copy entire directory structure
directory: "docs/"
destination_directory: "public-docs/"
```

#### Pull Request Parameters

```yaml
# Comprehensive PR configuration
pull_message: "🔧 Update production configuration"
pull_description: |
  **Configuration Update**
  
  - Updated production settings
  - Synced from ${{ github.repository }}
  - Commit: ${{ github.sha }}
  
  **Review Checklist:**
  - [ ] Configuration values are correct
  - [ ] No sensitive data exposed
  - [ ] Testing completed
reviewers: "config-admin,devops-lead"
team_reviewers: "platform-team,security-team"
```

#### Branch and Token Parameters

```yaml
# Advanced branch and token configuration
ref_branch: "main"                    # Target base branch
branch: "auto-update-${{ github.run_number }}"  # Dynamic branch name
token: "${{ secrets.CROSS_ORG_TOKEN }}"         # Cross-org token
```

## Common Use Cases

### 1. Configuration Management

```yaml
# Sync environment-specific configs
- uses: pal-paul/git-copy@v2.1.4
  with:
    owner: "your-org"
    repo: "config-repo"
    token: "${{ secrets.GITHUB_TOKEN }}"
    directory: "configs/${{ matrix.environment }}/"
    destination_directory: "deployments/${{ matrix.environment }}/"
```

### 2. Documentation Sync

```yaml
# Keep documentation in sync across repos
- uses: pal-paul/git-copy@v2.1.4
  with:
    owner: "your-org"
    repo: "docs-site"
    token: "${{ secrets.GITHUB_TOKEN }}"
    directory: "docs/"
    destination_directory: "content/api/"
```

### 3. Shared Component Distribution

```yaml
# Distribute shared libraries
- uses: pal-paul/git-copy@v2.1.4
  with:
    owner: "your-org"
    repo: "component-library"
    token: "${{ secrets.GITHUB_TOKEN }}"
    directory: "dist/"
    destination_directory: "vendor/shared-components/"
```

### 4. Infrastructure as Code

```yaml
# Deploy Terraform configurations
- uses: pal-paul/git-copy@v2.1.4
  with:
    owner: "your-org"
    repo: "terraform-infrastructure"
    token: "${{ secrets.GITHUB_TOKEN }}"
    directory: "terraform/modules/"
    destination_directory: "infrastructure/modules/"
```

## Troubleshooting

### Common Issues

#### Permission Denied

```
Error: 403 Forbidden
```

**Solution**: Ensure your token has the necessary permissions:

- `contents: write` - To create/update files
- `pull-requests: write` - To create pull requests
- `metadata: read` - To read repository metadata

#### Branch Already Exists

```
Error: Reference already exists
```

**Solution**: Use dynamic branch names:

```yaml
branch: "update-${{ github.sha }}"
# or
branch: "sync-${{ github.run_number }}"
```

#### File Not Found

```
Error: could not read file
```

**Solution**: Verify the source file path exists and use relative paths:

```yaml
file_path: "./config/app.json"  # ✅ Good
file_path: "/config/app.json"   # ❌ Avoid absolute paths
```

#### Large File Issues

For files larger than GitHub's limits:

```yaml
# Consider splitting large directories
directory: "docs/api/"           # ✅ Specific subdirectory
directory: "docs/"               # ❌ Might be too large
```

### Debug Mode

Enable verbose logging by setting debug environment variables:

```yaml
env:
  ACTIONS_STEP_DEBUG: true
  ACTIONS_RUNNER_DEBUG: true
```

### Token Permissions

Required token scopes:

- `repo` (for private repositories)
- `public_repo` (for public repositories)
- `workflow` (if modifying workflow files)

**For Cross-Organization Actions:**
When copying to a different organization, ensure the token has access to the target repository:

1. **Personal Access Token**: Create a PAT with repo access to the target organization
2. **Organization Token**: Use an organization-level token or app token
3. **Repository Access**: Token owner must have write access to the target repository

**Common Permission Issues:**

```yaml
# ❌ Using default GITHUB_TOKEN for cross-org copy
token: "${{ secrets.GITHUB_TOKEN }}"  # Only works within same repository

# ✅ Using organization or personal token 
token: "${{ secrets.CROSS_ORG_TOKEN }}"  # Has access to target repository
```

Example token creation:

```bash
# GitHub CLI
gh auth token --scopes repo,workflow

# Or create via GitHub Settings > Developer settings > Personal access tokens
```

### 401 Unauthorized Error

```
2025/06/18 11:13:11 failed to get branch master: 401 Unauthorized
```

This error occurs when the GitHub token doesn't have proper permissions to access the target repository.

**Common Causes & Solutions:**

1. **Cross-Organization Access**:

   ```yaml
   # ❌ Default GITHUB_TOKEN doesn't work across organizations
   token: "${{ secrets.GITHUB_TOKEN }}"
   
   # ✅ Use Personal Access Token or Organization Token
   token: "${{ secrets.PERSONAL_ACCESS_TOKEN }}"
   ```

2. **Missing Repository Access**:
   - Ensure the token owner has **write access** to target repository
   - For organization repos, token may need organization permissions
   - Check if repository is private and token has appropriate scope

3. **Incorrect Token Scopes**:
   Required scopes depend on repository type:

   ```bash
   # For private repositories
   repo
   
   # For public repositories  
   public_repo
   
   # For organization repositories
   repo, read:org
   ```

4. **Repository Name Issues**:
   - Verify `owner` and `repo` parameters are correct
   - Check for typos in repository names
   - Ensure target repository exists

**Debugging Steps:**

```yaml
# Add debug logging to your workflow
- name: Debug Repository Access
  run: |
    echo "Target: ${{ inputs.owner }}/${{ inputs.repo }}"
    # Test token access (replace with your values)
    curl -H "Authorization: token ${{ secrets.YOUR_TOKEN }}" \
         https://api.github.com/repos/${{ inputs.owner }}/${{ inputs.repo }}
```

**Creating Proper Tokens:**

For cross-organization use:

1. Go to GitHub Settings → Developer settings → Personal access tokens
2. Create token with `repo` scope (and `read:org` if needed)
3. Add token to secrets in source repository
4. Use the secret token in your action

## Contributing

We welcome contributions! Please see our [Contributing Guide](.github/CONTRIBUTING.md) for details on:

- Setting up the development environment
- Code style guidelines
- Testing requirements
- Pull request process

### Quick Start for Contributors

1. **Fork and clone the repository**

   ```bash
   git clone https://github.com/YOUR-USERNAME/git-copy.git
   cd git-copy
   ```

2. **Set up development environment**

   ```bash
   make install
   make test-all
   ```

3. **Make your changes and test**

   ```bash
   make fmt        # Format code
   make lint       # Run linting
   make test-all   # Run all tests
   make build      # Build application
   ```

4. **Submit a pull request**
   - Follow our [PR template](.github/pull_request_template.md)
   - Ensure all CI checks pass
   - Request review from mastertainers

## CI/CD Pipeline

This project uses GitHub Actions for continuous integration and deployment:

### Automated Testing (`test.yml`)

Triggered on: Push to `master`/`master`, Pull requests to `master`/`master`

**Test Matrix:**

- **Go versions**: 1.21, 1.22
- **Operating System**: Ubuntu Latest
- **Test Types**: Unit, Integration, Race Detection, Coverage

**Quality Checks:**

- Code formatting (`gofmt`)

- Linting (`golangci-lint`)
- Security scanning (`gosec`)
- Vulnerability checking (`govulncheck`)
- Application startup validation

**Artifacts:**

- Test coverage reports

- Coverage HTML reports
- Built binaries

### Dependency Updates (`dependency-update.yml`)

Triggered on: Weekly schedule (Mondays), Manual dispatch

**Automated Updates:**

- Go dependencies (`go get -u ./...`)

- GitHub Actions versions
- Security vulnerability scanning

**Pull Request Creation:**

- Automatic PR creation for dependency updates
- Comprehensive testing before merge
- Reviewer assignment and labeling

### Release Pipeline (`release.yml`)

Triggered on: Git tag push (`v*`)

**Release Process:**

1. **Pre-release Testing**: Full test suite on multiple Go versions

2. **Multi-platform Builds**: Linux, macOS, Windows (AMD64/ARM64)
3. **GitHub Release**: Automatic changelog generation and asset upload
4. **Docker Images**: Build and push to GitHub Container Registry
5. **Version Management**: Update major version tags for GitHub Actions

**Release Artifacts:**

- Binary executables for all platforms
- Checksums for verification
- Docker images with multiple tags
- Automated changelog from git history

### Security Monitoring

- Weekly vulnerability scans
- Dependency security audits
- Automatic issue creation for security alerts
- SARIF upload for GitHub Security tab

## Project Structure

```
git-copy/
├── .github/                    # GitHub configuration
│   ├── ISSUE_TEMPLATE/        # Bug report and feature request templates
│   ├── workflows/             # GitHub Actions workflows
│   ├── CONTRIBUTING.md        # Contribution guidelines
│   ├── SECURITY.md           # Security policy
│   └── pull_request_template.md
├── cmd/                       # Application entry point
│   └── cmd.go
├── internal/                  # Internal packages
│   └── gitcopy/              # Core application logic
│       └── gitcopy.go
├── test/                     # Test files
│   ├── cmd_test.go          # Core functionality tests
│   ├── integration_test.go   # Integration tests
│   ├── git_operations_test.go # Git operations tests
│   └── edge_cases_test.go    # Edge case tests
├── action.yml               # GitHub Action metadata
├── Dockerfile              # Container configuration
├── Makefile                # Build and development commands
├── .golangci.yml          # Linting configuration
└── README.md              # Project documentation
```

## Security

This project takes security seriously. Please see our [Security Policy](.github/SECURITY.md) for:

- Supported versions
- Vulnerability reporting process
- Security best practices
- Contact information

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Development

### Building and Running

The project includes a comprehensive Makefile with multiple targets for development, testing, and deployment:

#### Build Targets

- `make build` - Build the application binary
- `make install` - Install the application and dependencies
- `make clean` - Clean build artifacts

#### Test Targets

- `make test` - Run unit tests
- `make test-coverage` - Run tests with coverage analysis (generates `coverage.html`)
- `make test-race` - Run tests with race condition detection
- `make test-bench` - Run performance benchmarks
- `make test-all` - Run all tests including race detection, coverage, and startup validation
- `make test-startup` - Test application startup and environment validation

#### Run Targets

- `make run` - Build and run with test environment variables (validates startup)
- `make run-with-env` - Run with actual environment variables (requires proper GitHub token)

#### Example Usage

```bash
# Run all tests
make test-all

# Check test coverage
make test-coverage
open coverage.html

# Test application startup
make test-startup

# Build and run with test data
make run
```

### Testing

This project includes a comprehensive test suite covering:

#### Core Functionality Tests (`cmd_test.go`)

- Environment configuration and initialization
- File reading operations with error handling
- Directory traversal and recursive file discovery
- Input validation for file and directory operations
- Branch initialization logic
- Reviewer parsing (users and teams)
- Cross-platform path handling
- Default value assignment
- Performance benchmarks

#### Integration Tests (`integration_test.go`)

- End-to-end file processing workflows
- Batch operation simulation and validation
- Input validation scenarios
- Branch initialization logic
- Error handling in directory processing with permission issues

#### Git Operations Tests (`git_operations_test.go`)

- Git batch operation data structures
- Pull request creation and validation
- GitHub workflow validation
- Concurrent file operations
- Rate limiting scenarios
- Error recovery with exponential backoff retry logic

#### Edge Cases Tests (`edge_cases_test.go`)

- Complex nested directory structures
- Large file handling (1MB+ files)
- Special characters in file paths
- Empty file processing
- Symbolic link handling (Unix systems)
- Concurrent directory access
- File permission variations
- Resource cleanup validation
- Path normalization across platforms

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detection
make test-race

# Run benchmark tests
make test-bench

# Run all test targets
make test-all

# Or use Go directly
go test -v ./cmd
go test -v ./cmd -cover
go test -v ./cmd -race
go test -v ./cmd -bench=.
```

### Test Coverage

Current test coverage: **17.8%** of statements

The test suite includes:

- **29 test functions** covering various scenarios
- **Benchmark tests** for performance validation
- **Concurrent tests** for race condition detection
- **Error simulation** for resilience testing
- **Cross-platform compatibility** tests

### Test Categories

1. **Unit Tests**: Test individual functions in isolation
2. **Integration Tests**: Test complete workflows and interactions
3. **Performance Tests**: Validate performance under load
4. **Edge Case Tests**: Handle unusual or extreme conditions
5. **Concurrency Tests**: Ensure thread safety
6. **Error Handling Tests**: Validate graceful degradation

### Fixed Issues

Recent improvements include:

- **Branch Initialization Fix**: Fixed critical bug where directory operations didn't initialize git branches properly
- **Error Handling**: Added `continue` statements to skip failed files instead of stopping entire operations
- **Resource Management**: Added proper file handle cleanup with `defer file.Close()`
- **Cross-platform Paths**: Replaced string concatenation with `filepath.Join()` for proper path handling
- **Enhanced Logging**: Added informative batch operation messages

### Build and Test Commands

Available Makefile targets:

- `make test`: Run basic tests
- `make test-coverage`: Generate coverage report
- `make test-race`: Run with race detection
- `make test-bench`: Run benchmark tests
- `make test-all`: Run all test types
- `make lint`: Run basic linting (CI-compatible)
- `make lint-local`: Run comprehensive local linting
- `make fmt`: Format Go code
- `make build`: Build the application
- `make install-tools`: Install development tools
- `make tidy`: Update and tidy Go modules

### Linting

The project uses a two-tier linting approach:

#### CI/CD Linting (`make lint`)

- Basic, stable checks: `go vet` and `go fmt`
- No external dependencies or version conflicts
- Reliable across different Go and tool versions

#### Local Development Linting (`make lint-local`)

- Comprehensive checks using golangci-lint
- Security scanning and vulnerability detection
- Version-aware configuration
- Enhanced developer feedback

```bash

# Install development tools first
make install-tools

# Run comprehensive local linting
make lint-local

# Or run the script directly
./scripts/lint.sh

# Basic linting (CI-style)
make lint
```

**Note**: Advanced linting is now local-only to avoid CI version conflicts. The CI pipeline focuses on reliable, core checks while developers get the full linting experience locally.

- `make test-bench`: Run benchmark tests
- `make test-all`: Run all test types
- `make lint`: Run basic linting (CI-compatible)
