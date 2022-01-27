# Openvpn Processor
[![CI](https://github.com/vpnbeast/openvpn-processor/workflows/CI/badge.svg?event=push)](https://github.com/vpnbeast/openvpn-processor/actions?query=workflow%3ACI)
[![Docker pulls](https://img.shields.io/docker/pulls/vpnbeast/openvpn-processor)](https://hub.docker.com/r/vpnbeast/openvpn-processor/)
[![Go Report Card](https://goreportcard.com/badge/github.com/vpnbeast/openvpn-processor)](https://goreportcard.com/report/github.com/vpnbeast/openvpn-processor)
[![codecov](https://codecov.io/gh/vpnbeast/openvpn-processor/branch/master/graph/badge.svg)](https://codecov.io/gh/vpnbeast/openvpn-processor)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=vpnbeast_openvpn-processor&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=vpnbeast_openvpn-processor)
[![Go version](https://img.shields.io/github/go-mod/go-version/vpnbeast/openvpn-processor)](https://github.com/vpnbeast/openvpn-processor)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This is a scheduled application which fetches usable Openvpn servers from **VPNGATE_URL** environment variable and
then inserts into [vpnbeast-mysql](https://github.com/vpnbeast/vpnbeast-mysql) database.

## Prerequisites
openvpn-processor requires [vpnbeast/config-service](https://github.com/vpnbeast/config-service) to fetch configuration. Configurations
are stored at [vpnbeast/properties](https://github.com/vpnbeast/properties).

## Configuration
This project fetches the configuration from [config-service](https://github.com/vpnbeast/config-service).
But you can still override them with environment variables:
```
DB_URL
DB_DRIVER
TICKER_INTERVAL_MIN
DB_MAX_OPEN_CONN
DB_MAX_IDLE_CONN
DB_CONN_MAX_LIFETIME_MIN
HEALTH_CHECK_MAX_TIMEOUT_MIN
DIAL_TCP_TIMEOUT_SECONDS
METRICS_PORT
METRICS_ENDPOINT
WRITE_TIMEOUT_SECONDS
READ_TIMEOUT_SECONDS
HEALTH_PORT
HEALTH_ENDPOINT
```

## Development
This project requires below tools while developing:
- [Golang 1.17](https://golang.org/doc/go1.17)
- [pre-commit](https://pre-commit.com/)
- [golangci-lint](https://golangci-lint.run/usage/install/) - required by [pre-commit](https://pre-commit.com/)

## License
Apache License 2.0
