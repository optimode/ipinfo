# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.1] - 2026-03-11

### Removed

- cache hit log messages from stderr for cleaner output

## [0.3.0] - 2026-03-11

### Added

- `--version` flag with version, git commit, and build time injected at build time
- remove duplicate empty IP check in run function

## [0.2.0] - 2026-03-11

### Added

- `-s` flag as shorthand for `--format summary`
- combine all input sources (cli args, `--file`, stdin) instead of mutually exclusive

### Fixed

- GitHub Actions workflow triggered twice on push (now only on tags and PRs)

## [0.1.0] - 2026-03-11

### Added

- IP geolocation lookup via ip-api.com Pro API
- output formats: table, summary, json, csv
- input from CLI args, `--file`, or stdin pipe
- /24 subnet cache with `--no-cache` to disable
- configurable concurrency with `-c` flag
- config via `/etc/ipinfo/config.yaml` and `IPINFO_` env vars
- cross-platform builds (linux/darwin, amd64/arm64)
- GitHub Actions CI/CD with artifact upload and release
- one-liner install script (`deployments/install.sh`)
