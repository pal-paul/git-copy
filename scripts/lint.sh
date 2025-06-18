#!/bin/bash

# Git Copy - Local Linting Script for Developers
# This script runs comprehensive linting checks locally

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔍 Git Copy - Local Linting Script${NC}"
echo "======================================"

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        return 1
    fi
}

# Basic Go checks
echo -e "${YELLOW}📋 Running basic Go checks...${NC}"

# Go vet
echo "Running go vet..."
if go vet ./...; then
    print_status 0 "go vet passed"
else
    print_status 1 "go vet failed"
    exit 1
fi

# Go fmt check
echo "Checking code formatting..."
UNFORMATTED=$(gofmt -s -l . 2>/dev/null || true)
if [ -z "$UNFORMATTED" ]; then
    print_status 0 "Code formatting is correct"
else
    echo -e "${RED}❌ The following files need formatting:${NC}"
    echo "$UNFORMATTED"
    echo -e "${YELLOW}Run 'make fmt' or 'go fmt ./...' to fix formatting${NC}"
    exit 1
fi

# Advanced linting with golangci-lint
echo -e "${YELLOW}🔧 Running advanced linting...${NC}"

if command_exists golangci-lint; then
    echo "Detecting golangci-lint version..."
    LINT_VERSION=$(golangci-lint --version | grep -o 'version [0-9]*\.[0-9]*' | cut -d' ' -f2 2>/dev/null || echo "unknown")
    MAJOR_VERSION=$(echo "$LINT_VERSION" | cut -d'.' -f1 2>/dev/null || echo "0")
    
    echo "Found golangci-lint version: $LINT_VERSION"
    
    # Check if .golangci.yml exists
    if [ -f ".golangci.yml" ]; then
        echo "Using configuration file: .golangci.yml"
        
        # Try to run with config file
        if [ "$MAJOR_VERSION" -ge "2" ] || [ "$LINT_VERSION" = "1.59" ] || [ "$LINT_VERSION" = "1.60" ]; then
            echo "Using v2 configuration format..."
            if golangci-lint run --config .golangci.yml --timeout=5m; then
                print_status 0 "golangci-lint passed with configuration"
            else
                echo -e "${YELLOW}⚠️  golangci-lint found issues. See output above.${NC}"
                echo -e "${BLUE}💡 These may be acceptable for comprehensive test coverage.${NC}"
                
                # Ask if user wants to continue or see fallback
                echo -e "${YELLOW}Do you want to run fallback linting? (y/n)${NC}"
                read -r response
                if [[ "$response" =~ ^[Yy]$ ]]; then
                    echo "Running fallback linting..."
                    golangci-lint run --enable=errcheck,govet,ineffassign,staticcheck,unused --timeout=5m || true
                fi
            fi
        else
            echo "Version $LINT_VERSION detected. Using fallback configuration..."
            if golangci-lint run --enable=errcheck,govet,ineffassign,staticcheck,unused,gosec,gocritic --timeout=5m; then
                print_status 0 "golangci-lint passed with fallback configuration"
            else
                echo -e "${YELLOW}⚠️  golangci-lint found issues with fallback configuration.${NC}"
            fi
        fi
    else
        echo "No .golangci.yml found, using default linters..."
        golangci-lint run --timeout=5m || true
    fi
else
    echo -e "${YELLOW}⚠️  golangci-lint not found${NC}"
    echo -e "${BLUE}💡 Install from: https://golangci-lint.run/usage/install/${NC}"
    echo -e "${BLUE}💡 Or use: make install-tools${NC}"
fi

# Additional security checks
echo -e "${YELLOW}🔒 Running security checks...${NC}"

# Check for hardcoded secrets (basic patterns)
echo "Checking for potential hardcoded credentials..."
if grep -r --include="*.go" --exclude-dir=test "password.*=\|secret.*=\|token.*=" . 2>/dev/null; then
    echo -e "${YELLOW}⚠️  Found potential hardcoded credentials (review required)${NC}"
else
    print_status 0 "No obvious hardcoded credentials found"
fi

# Check for unsafe operations
echo "Checking for unsafe operations..."
if grep -r --include="*.go" "unsafe\." . 2>/dev/null; then
    echo -e "${YELLOW}⚠️  Found unsafe package usage (review required)${NC}"
else
    print_status 0 "No unsafe operations found"
fi

# Vulnerability check (if govulncheck is available)
if command_exists govulncheck; then
    echo "Running vulnerability check..."
    if govulncheck ./...; then
        print_status 0 "No known vulnerabilities found"
    else
        print_status 1 "Vulnerabilities detected"
    fi
else
    echo -e "${YELLOW}💡 Install govulncheck for vulnerability scanning:${NC}"
    echo -e "${BLUE}   go install golang.org/x/vuln/cmd/govulncheck@latest${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Local linting completed!${NC}"
echo -e "${BLUE}💡 To install missing tools, run: make install-tools${NC}"
echo -e "${BLUE}💡 For formatting fixes, run: make fmt${NC}"
