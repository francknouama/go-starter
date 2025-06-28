# Security Policy

## Supported Versions

We actively support the latest major version of go-starter. Security updates are provided for:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Security Features

go-starter includes comprehensive security measures to protect against common vulnerabilities:

### Template Security
- **Template Injection Prevention**: All templates are validated for dangerous patterns
- **Function Whitelisting**: Only safe template functions are allowed
- **Path Traversal Protection**: Template paths are validated to prevent directory traversal
- **Resource Limits**: Templates have size and execution time limits

### Input Validation
- **Project Name Sanitization**: Dangerous characters and reserved names are rejected
- **Module Path Validation**: Go module paths are validated against malicious patterns
- **Path Validation**: Output paths are checked for traversal attempts
- **Resource Limits**: File count, directory count, and total size limits are enforced

### Generated Code Security
- **No Hardcoded Secrets**: Templates never include hardcoded credentials
- **Secure Defaults**: All generated code follows security best practices
- **Logger Security**: Structured logging prevents injection attacks
- **Configuration Validation**: All configuration values are sanitized

### CI/CD Security
- **Automated Scanning**: GitHub Actions workflows scan for vulnerabilities
- **Dependency Checking**: Regular scans for vulnerable dependencies
- **Template Validation**: All templates are automatically validated for security
- **Secret Detection**: Automated detection of potential hardcoded secrets

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security issue, please follow these steps:

### 1. Do Not Create Public Issues
**Please do not report security vulnerabilities through public GitHub issues.**

### 2. Report Privately
Send details to our security team at: **security@go-starter.dev**

Include the following information:
- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact assessment
- Suggested fix (if any)

### 3. Response Timeline
- **Initial Response**: Within 24 hours
- **Assessment**: Within 72 hours
- **Fix Development**: 1-7 days (depending on severity)
- **Release**: As soon as possible after fix verification

### 4. Disclosure Policy
We follow responsible disclosure:
- We will acknowledge receipt of your report
- We will investigate and confirm the issue
- We will develop and test a fix
- We will release the fix and notify users
- We will publicly disclose the vulnerability after users have had time to update

## Security Testing

### Running Security Tests
```bash
# Run all security tests
go test -v ./tests/security/...

# Scan templates for security issues
go run main.go security scan-templates

# Validate project configuration
go run main.go security scan-config project.yaml
```

### Manual Security Checks
```bash
# Check for vulnerable dependencies
go list -json -m all | nancy sleuth

# Run Gosec static analysis
gosec ./...

# Check for known vulnerabilities
govulncheck ./...
```

## Security Best Practices for Users

### Template Development
1. **Never include secrets**: Don't hardcode API keys, passwords, or tokens
2. **Validate all inputs**: Use the provided sanitization functions
3. **Limit complexity**: Avoid deeply nested loops or complex logic
4. **Use safe functions**: Stick to the whitelisted template functions

### Project Generation
1. **Review generated code**: Always review generated projects before use
2. **Update dependencies**: Keep generated project dependencies up to date
3. **Enable security features**: Use the security scanning commands regularly
4. **Follow Go security guidelines**: Apply standard Go security practices

### Configuration Management
1. **Use environment variables**: Don't hardcode sensitive values in configs
2. **Validate module paths**: Ensure module paths point to trusted repositories
3. **Limit permissions**: Use minimal file permissions for generated projects
4. **Regular updates**: Keep go-starter updated to the latest version

## Security Architecture

### Defense in Depth
go-starter implements multiple layers of security:

1. **Input Validation Layer**
   - Sanitizes all user inputs
   - Validates module paths and project names
   - Prevents path traversal attacks

2. **Template Security Layer**
   - Validates template syntax and patterns
   - Enforces function whitelisting
   - Implements resource limits

3. **Generation Security Layer**
   - Validates output paths
   - Enforces file count and size limits
   - Prevents dangerous file operations

4. **Code Analysis Layer**
   - Static analysis of generated code
   - Dependency vulnerability scanning
   - Configuration security validation

### Threat Model

#### Template Injection Attacks
**Threat**: Malicious templates executing dangerous code
**Mitigation**: Template validation, function whitelisting, syntax checking

#### Path Traversal Attacks  
**Threat**: Templates accessing unauthorized file system locations
**Mitigation**: Path validation, relative path enforcement, directory restrictions

#### Resource Exhaustion Attacks
**Threat**: Templates consuming excessive resources
**Mitigation**: Size limits, execution timeouts, complexity restrictions

#### Dependency Vulnerabilities
**Threat**: Generated projects with vulnerable dependencies
**Mitigation**: Dependency scanning, version pinning, security updates

#### Configuration Injection
**Threat**: Malicious input through project configuration
**Mitigation**: Input sanitization, validation, type checking

## Security Updates

Security updates are released as patch versions (e.g., 1.0.1) and include:
- CVE fixes for dependencies
- Template security improvements
- Input validation enhancements
- New security scanning features

Subscribe to releases on GitHub to stay informed about security updates.

## Security Resources

- [OWASP Go Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Go_SCP_Cheat_Sheet.html)
- [Go Security Policy](https://golang.org/security)
- [CIS Go Security Benchmarks](https://www.cisecurity.org/benchmark/go)
- [NIST Secure Software Development Framework](https://csrc.nist.gov/Projects/ssdf)

## Contact

For security-related questions or concerns:
- Email: security@go-starter.dev
- Security Team: [@security-team](https://github.com/orgs/go-starter/teams/security-team)

---

*This security policy is regularly reviewed and updated. Last updated: 2024-01-01*