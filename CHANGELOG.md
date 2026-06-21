# Changelog

All notable changes are documented here.  
Format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) · Versioning follows [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] — 2026-06-21

First release under active maintenance. Forked from [adaptant-labs/edgetpu-exporter](https://github.com/adaptant-labs/edgetpu-exporter), which was last updated in January 2021.

### Added
- `PORT` environment variable for port configuration (flag still overrides)
- HTTP server timeouts: read 30s, write 30s, idle 60s (DoS hardening)
- Non-root `exporter` user (uid 1000) in Docker image
- GitHub Actions CI workflow — build, vet, `govulncheck` on every push and PR
- Structured issue templates: bug report, feature request
- Pull request template with verification checklist
- `CODEOWNERS` — auto-assigns reviewer on PRs
- `SECURITY.md` — responsible disclosure policy
- `CONTRIBUTING.md` — build steps, conventions, commit style
- `CHANGELOG.md`
- GitHub release workflow triggered on `v*` tags

### Changed
- Go module path: `github.com/adaptant-labs/edgetpu-exporter` → `github.com/just5ky/edgetpu-exporter`
- Go minimum version: `1.13` → `1.25`; toolchain pinned to `go1.26.4`
- Docker builder: `golang` (unpinned) → `golang:1.26.4-alpine`
- Docker final stage: `scratch` → `alpine:3.22` (required for user management)
- Docker build flags: added `-trimpath -ldflags="-s -w"` for reproducible stripped binaries
- Kubernetes DaemonSet image: `adaptant/edgetpu-exporter:latest` → `justsky/edgetpu-exporter:latest`
- Dependency upgrades:

  | Module | From | To |
  |--------|------|----|
  | `prometheus/client_golang` | v1.21.1 | v1.23.2 |
  | `prometheus/common` | v0.62.0 | v0.69.0 |
  | `prometheus/procfs` | v0.15.1 | v0.20.1 |
  | `prometheus/client_model` | v0.6.1 | v0.6.2 |
  | `klauspost/compress` | v1.17.11 | v1.18.6 |
  | `golang.org/x/sys` | v0.28.0 | v0.46.0 |
  | `google.golang.org/protobuf` | v1.36.1 | v1.36.11 |

### Removed
- `.travis.yml` — replaced by GitHub Actions
- Unused `github.com/json-iterator/go` explicit dependency
- `GO111MODULE` env var — unnecessary since Go 1.16
- `io/ioutil` usage — replaced with `os.ReadFile` (deprecated since Go 1.16)

### Fixed
- C-style `for i := 0; i < len(...)` loops replaced with idiomatic `range`
- Broken `[node-feature-discovery]` link in documentation
- Wrong Docker Hub and issue tracker links pointing to upstream `adaptant-labs`

### Security
- Resolved 18 Go standard library CVEs by upgrading toolchain from `go1.24.4` → `go1.26.4`
  (affects `crypto/tls`, `crypto/x509`, `net`, `net/url`, `net/textproto`, `os`, `encoding/asn1`, `encoding/pem`)
