# Security Policy

## Supported versions

Only the latest published image (`justsky/edgetpu-exporter:latest`) and the `latest` branch receive security fixes.

## Reporting a vulnerability

**Do not open a public GitHub issue for security vulnerabilities.**

Please report security issues by emailing the maintainer directly via the contact listed on the [GitHub profile](https://github.com/Just5KY). Include:

- A description of the vulnerability
- Steps to reproduce
- Potential impact
- Any suggested fix (optional)

You can expect an acknowledgement within **72 hours** and a resolution or status update within **7 days**.

## Dependency scanning

Dependencies are monitored automatically via:

- [Dependabot](.github/dependabot.yml) — daily updates for Go modules, Docker base images, and GitHub Actions
- [`govulncheck`](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) — run in CI on every push and pull request
