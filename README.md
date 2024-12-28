# OpenVPN Exporter

[![license](https://img.shields.io/github/license/natrontech/openvpn-exporter)](https://github.com/natrontech/openvpn-exporter/blob/main/LICENSE)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/natrontech/openvpn-exporter/badge)](https://securityscorecards.dev/viewer/?uri=github.com/natrontech/openvpn-exporter)
[![release](https://img.shields.io/github/v/release/natrontech/openvpn-exporter)](https://github.com/natrontech/openvpn-exporter/releases)
[![go-version](https://img.shields.io/github/go-mod/go-version/natrontech/openvpn-exporter)](https://github.com/natrontech/openvpn-exporter/blob/main/go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/natrontech/openvpn-exporter)](https://goreportcard.com/report/github.com/natrontech/openvpn-exporter)
[![SLSA 3](https://slsa.dev/images/gh-badge-level3.svg)](https://slsa.dev)

---

Export [OpenVPN Community](https://openvpn.net/community-downloads/) statistics to [Prometheus](https://prometheus.io/).

Metrics are retrieved using the OpenVPN Status file.

## Exported Metrics

| Metric          | Meaning                                            | Labels                         |
| --------------- | -------------------------------------------------- | ------------------------------ |
| openvpn_up      | Was the last query of OpenVPN Exporter successful? |                                |
| openvpn_version | Version of OpenVPN Exporter                        | `version`, `repoid`, `release` |

## Flags / Environment Variables

```bash
$ ./openvpn-exporter -help
```

You can use the following flags to configure the exporter. All flags can also be set using environment variables. Environment variables take precedence over flags.

| Flag                         | Environment Variable     | Description                                          | Default                       |
| ---------------------------- | ------------------------ | ---------------------------------------------------- | ----------------------------- |
| `openvpn.loglevl`            | `OPENVPN_LOGLEVEL`       | Log level (debug, info)                              | `info`                        |
| `openvpn.status-file`        | `OPENVPN_STATUS_FILE`    | Path to OpenVPN status file                          | `/var/log/openvpn-status.log` |
| `openvpn.metrics-path`       | `OPENVPN_METRICS_PATH`   | Path under which to expose metrics                   | `/metrics`                    |
| `openvpn.web.listen-address` | `OPENVPN_LISTEN_ADDRESS` | Address to listen on for web interface and telemetry | `:9999`                       |
|                              |                          |                                                      |                               |
## Release

Each release of the application includes Go-binary archives, checksums file, SBOMs and container images. 

The release workflow creates provenance for its builds using the [SLSA standard](https://slsa.dev), which conforms to the [Level 3 specification](https://slsa.dev/spec/v1.0/levels#build-l3). Each artifact can be verified using the `slsa-verifier` or `cosign` tool (see [Release verification](SECURITY.md#release-verification)).
