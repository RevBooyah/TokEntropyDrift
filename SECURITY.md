# Security Policy

## Supported Versions

Use this section to tell people about which versions of your project are currently being supported with security updates.

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability within TokEntropyDrift, please send an email to [your-email@example.com]. All security vulnerabilities will be promptly addressed.

Please include the following information in your report:

- A description of the vulnerability
- Steps to reproduce the issue
- Potential impact
- Suggested fix (if available)

## Security Best Practices

When using TokEntropyDrift:

1. **Keep dependencies updated**: Regularly update your Go dependencies
2. **Review configuration**: Ensure your `ted.config.yaml` doesn't expose sensitive information
3. **Validate input**: Always validate input files before processing
4. **Monitor output**: Review analysis results for any unexpected behavior

## Disclosure Policy

- Security vulnerabilities will be disclosed via GitHub Security Advisories
- Patches will be released as soon as possible
- CVE numbers will be requested for significant vulnerabilities 