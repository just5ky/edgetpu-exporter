# Contributing

Contributions are welcome. This is a small project — keep changes focused.

## Getting started

```bash
git clone https://github.com/Just5KY/edgetpu-exporter.git
cd edgetpu-exporter
go build ./...
go vet ./...
```

Requirements: **Go 1.25+**

## Before submitting a PR

1. **Build** — `go build -trimpath -ldflags="-s -w" ./...`
2. **Vet** — `go vet ./...`
3. **Vuln scan** — `govulncheck ./...` (install: `go install golang.org/x/vuln/cmd/govulncheck@latest`)
4. **Tidy** — `go mod tidy` if you changed dependencies
5. Update `README.md` if you changed flags, env vars, metrics, or behaviour

## What belongs here

- Bug fixes for device discovery or metric collection
- Support for new EdgeTPU hardware variants (new USB VID/PID, new sysfs paths)
- Additional metrics exposed via sysfs
- CI / build improvements

## What doesn't belong here

- Unrelated refactors
- New external dependencies without strong justification
- Features unrelated to EdgeTPU monitoring

## Code style

Follow standard Go conventions (`gofmt`, idiomatic error handling). No comments that just restate what the code does.

## Commit messages

Use plain, imperative present-tense subject lines, e.g.:

```
Add temperature support for USB devices
Fix nil panic when sysfs path is missing
Bump prometheus/client_golang to v1.23.2
```

## Reporting bugs or requesting features

Use the [issue tracker](https://github.com/Just5KY/edgetpu-exporter/issues) with the provided templates.
