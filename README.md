# edgetpu-exporter

[![CI](https://github.com/Just5KY/edgetpu-exporter/actions/workflows/ci.yml/badge.svg)](https://github.com/Just5KY/edgetpu-exporter/actions/workflows/ci.yml)
[![Docker Build](https://github.com/Just5KY/edgetpu-exporter/actions/workflows/docker.yml/badge.svg)](https://github.com/Just5KY/edgetpu-exporter/actions/workflows/docker.yml)
[![Docker Pulls](https://img.shields.io/docker/pulls/justsky/edgetpu-exporter)](https://hub.docker.com/r/justsky/edgetpu-exporter)
[![Docker Image Size](https://img.shields.io/docker/image-size/justsky/edgetpu-exporter/latest)](https://hub.docker.com/r/justsky/edgetpu-exporter)
[![Go Version](https://img.shields.io/badge/go-1.25+-00ADD8?logo=go)](go.mod)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)

Prometheus exporter for [Google Coral / EdgeTPU](https://coral.ai) accelerator metrics. Exposes device count and per-device temperature via a lightweight HTTP endpoint compatible with any Prometheus scrape setup.

> **Fork notice** — This is an actively maintained fork of [adaptant-labs/edgetpu-exporter](https://github.com/adaptant-labs/edgetpu-exporter), which has been unmaintained since January 2021. See [what's changed](#whats-changed-from-upstream).

---

## Contents

- [Requirements](#requirements)
- [Quick start](#quick-start)
- [Configuration](#configuration)
- [Metrics](#metrics)
- [Docker](#docker)
- [Kubernetes](#kubernetes)
- [What's changed from upstream](#whats-changed-from-upstream)
- [Contributing](#contributing)
- [Security](#security)
- [License](#license)

---

## Requirements

- Linux with sysfs (`/sys`) access
- Google Coral EdgeTPU — PCIe (Apex class), USB Accelerator, or Dev Board
- **Docker** (recommended) _or_ Go 1.25+ to build from source

---

## Quick start

**Docker (recommended):**

```bash
docker run -d \
  --name edgetpu-exporter \
  -p 8080:8080 \
  -v /sys:/host-sys:ro \
  justsky/edgetpu-exporter:latest \
  -sysfs /host-sys
```

**Kubernetes DaemonSet:**

```bash
kubectl apply -f https://raw.githubusercontent.com/Just5KY/edgetpu-exporter/latest/edgetpu-daemonset.yaml
```

**From source:**

```bash
git clone https://github.com/Just5KY/edgetpu-exporter.git
cd edgetpu-exporter
go build -o edgetpu-exporter .
./edgetpu-exporter
```

---

## Configuration

`edgetpu-exporter` runs without any required configuration. All settings are optional.

```
Usage: edgetpu-exporter [flags]

Environment variables:
  PORT    Port to listen on (default: 8080)

Flags:
  -port int
        Port to listen on — overrides $PORT (default 8080)
  -sysfs string
        Mountpoint of sysfs to scan (default "/sys")
```

The `-port` flag takes precedence over the `PORT` environment variable.

```bash
# Via env var
PORT=9090 ./edgetpu-exporter

# Via flag
./edgetpu-exporter -port 9090

# Custom sysfs (e.g. in a container with a host bind-mount)
./edgetpu-exporter -sysfs /host-sys
```

---

## Metrics

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `edgetpu_num_devices` | Gauge | — | Total number of EdgeTPU devices detected |
| `edgetpu_temperature_celsius` | Gauge | `name` | Per-device temperature in °C |

> Temperature is only available for PCIe-attached (Apex class) devices. USB-attached Coral accelerators do not expose a temperature sensor via sysfs.

**Example scrape output:**

```
# HELP edgetpu_num_devices Number of EdgeTPU devices
# TYPE edgetpu_num_devices gauge
edgetpu_num_devices 1
# HELP edgetpu_temperature_celsius EdgeTPU device temperature in Celsius
# TYPE edgetpu_temperature_celsius gauge
edgetpu_temperature_celsius{name="apex_0"} 49.3
```

**Prometheus scrape config:**

```yaml
scrape_configs:
  - job_name: edgetpu
    static_configs:
      - targets: ['<node-ip>:8080']
```

---

## Docker

Multi-arch images (`linux/amd64`, `linux/arm64`) are published to Docker Hub on every push to `latest` and on version tags.

```bash
# Latest
docker pull justsky/edgetpu-exporter:latest

# Specific version
docker pull justsky/edgetpu-exporter:v1.0.0
```

Images run as a **non-root user** (`exporter`, uid 1000) on `alpine:3.22`.

---

## Kubernetes

The included [`edgetpu-daemonset.yaml`](edgetpu-daemonset.yaml) deploys a `DaemonSet` with node affinity targeting any node labelled by the sources below.

```bash
kubectl apply -f https://raw.githubusercontent.com/Just5KY/edgetpu-exporter/latest/edgetpu-daemonset.yaml
```

**Supported node labels:**

| Label | Source | Device |
|-------|--------|--------|
| `kkohtaka.org/edgetpu` | [EdgeTPU Device Plugin][edgetpu-device-plugin] | USB |
| `feature.node.kubernetes.io/usb-fe_1a6e_089a.present` | [node-feature-discovery] | USB |
| `feature.node.kubernetes.io/pci-0880_1ac1.present` | [node-feature-discovery] | PCIe / Dev Board |
| `beta.devicetree.org/fsl-imx8mq-phanbell` | [k8s-dt-node-labeller] | Coral Dev Board |

The DaemonSet mounts `/sys` from the host as read-only at `/host-sys` and passes `-sysfs /host-sys`.

---

## What's changed from upstream

This fork picks up where [adaptant-labs/edgetpu-exporter](https://github.com/adaptant-labs/edgetpu-exporter) left off. Key changes since the last upstream commit (January 2021):

- **Security** — 18 Go stdlib CVEs resolved; toolchain upgraded from go1.24 → go1.26.4
- **HTTP hardening** — read / write / idle timeouts added (prevents slow-client DoS)
- **Non-root container** — final image runs as uid 1000, not root
- **Pinned base image** — `golang:1.26.4-alpine` + `alpine:3.22` (no more unpinned `golang` tag)
- **Modern Go** — `io/ioutil` removed, idiomatic `range` loops, go 1.25 minimum
- **All deps current** — `prometheus/client_golang` v1.23.2, full dep tree updated
- **PORT env var** — twelve-factor app friendly port configuration
- **CI pipeline** — GitHub Actions with `govulncheck` on every PR; Travis CI removed
- **Dependabot** — daily updates for Go modules, Docker images, and Actions

Full details in [CHANGELOG.md](CHANGELOG.md).

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## Security

See [SECURITY.md](SECURITY.md) for the vulnerability disclosure policy.  
Dependencies are scanned daily via Dependabot and on every CI run via `govulncheck`.

---

## License

Apache 2.0 — see [LICENSE](LICENSE).  
Original work © [Adaptant Solutions AG](https://github.com/adaptant-labs). Fork maintained by [@Just5KY](https://github.com/Just5KY).

[node-feature-discovery]: https://github.com/kubernetes-sigs/node-feature-discovery
[edgetpu-device-plugin]: https://github.com/kkohtaka/edgetpu-device-plugin
[k8s-dt-node-labeller]: https://github.com/adaptant-labs/k8s-dt-node-labeller
