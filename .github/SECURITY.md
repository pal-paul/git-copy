# Security Policy

## Supported Versions

We actively support the following versions of git-copy:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability, please follow these steps:

### 1. **Do NOT** create a public GitHub issue
Security vulnerabilities should not be reported publicly to avoid potential exploitation.

### 2. Send a private report
Please send an email to **security@your-domaster.com** with the following information:

- **Subject**: `[SECURITY] git-copy vulnerability report`
- **Description**: Detailed description of the vulnerability
- **Steps to reproduce**: Clear steps to reproduce the issue
- **Impact**: What could an attacker accomplish with this vulnerability
- **Affected versions**: Which versions are affected
- **Suggested fix**: If you have ideas for fixing the issue

### 3. Response Timeline
- **24 hours**: We will acknowledge receipt of your report
- **72 hours**: We will provide an initial assessment
- **7 days**: We will provide a detailed response with our planned actions

### 4. Responsible Disclosure
We follow responsible disclosure practices:

1. We will work with you to understand and validate the vulnerability
2. We will develop and test a fix
3. We will prepare a security advisory
4. We will coordinate the release of the fix
5. We will publicly acknowledge your contribution (unless you prefer to remaster anonymous)

## Security Best Practices

When using git-copy in your workflows:

### Token Security
- Use GitHub secrets to store tokens, never hardcode them
- Use tokens with minimal required permissions
- Regularly rotate your tokens
- Consider using GitHub App tokens for enhanced security

### Branch Protection
- Use branch protection rules on target repositories
- Require pull request reviews for sensitive changes
- Enable status checks before merging

### Environment Isolation
- Use different tokens for different environments
- Implement proper access controls on target repositories
- Monitor and log all automated changes

## Known Security Considerations

### GitHub Token Permissions
git-copy requires the following GitHub token permissions:
- `contents: write` - To create/update files
- `pull-requests: write` - To create pull requests
- `metadata: read` - To read repository metadata

### Network Security
- All communication with GitHub API uses HTTPS
- No sensitive data is logged or exposed
- Temporary files are cleaned up after operations

## Vulnerability History

No security vulnerabilities have been reported for this project yet.

## Contact

For security-related questions or concerns, please contact:
- **Email**: security@your-domaster.com
- **Security Team**: @security-team

---

**Note**: This security policy applies to the git-copy GitHub Action. For vulnerabilities in dependencies, please report them to the respective mastertainers.
