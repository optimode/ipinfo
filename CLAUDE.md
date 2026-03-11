# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

```bash
go build -o ipinfo .                # build binary
go run . 8.8.8.8                    # run directly
go vet ./...                        # lint
go test ./...                       # run all tests (none yet)
```

Requires an API key: set `IPINFO_API_KEY` env var, `--api-key` flag, or `/etc/ipinfo/config.yaml`.

## Architecture

Go CLI tool that queries the ip-api.com Pro API for IP geolocation. Uses cobra for CLI and viper for config (flag > env > config file).

- `main.go` — entrypoint, calls `cmd.Execute()`
- `cmd/root.go` — cobra command definition, flag/config setup, worker pool orchestration. The `run` function collects IPs, spins up a concurrent worker pool, and calls `processIP` per IP
- `internal/api/` — HTTP client for `pro.ip-api.com/json/{ip}`. `Response` struct is the shared data type used across all packages
- `internal/cache/` — in-memory /24 subnet cache (thread-safe). On hit, returns cached response with `.Query` replaced by actual IP
- `internal/format/` — output printer supporting table, summary, json, csv. Thread-safe via mutex for concurrent writes
- `internal/input/` — reads IPs from `io.Reader`, skips comments (`#`) and blank lines

## Key Design Decisions

- `api.Response` is the central data type flowing through cache → format
- The format `Printer` and subnet `SubnetCache` are both mutex-protected for concurrent access from the worker pool
- Config precedence: CLI flag > env var (`IPINFO_PREFIX`) > `/etc/ipinfo/config.yaml` > defaults
- Exit codes: 0 = success, 1 = partial errors, 2 = usage/config error

## Versioning

Semantic Versioning, Conventional Commits, Keep a Changelog, GitHub Flow (main + rc tags).
