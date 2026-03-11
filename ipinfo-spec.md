# ipinfo – IP geolocation CLI tool specification

## Overview

Go-based CLI tool that queries the ip-api.com Pro API for IP address geolocation and metadata.
Replaces the existing bash script with better performance, /24 subnet caching, and concurrent lookups.

## Binary name

`ipinfo`

## Dependencies

- `github.com/spf13/cobra` – CLI
- `github.com/spf13/viper` – config (API key)
- Standard library only for HTTP, JSON, formatting

## Configuration

Lookup order (flag > env > config file):

| Setting | Flag | Env var | Config key |
|---------|------|---------|------------|
| API key | `--api-key` | `IPINFO_API_KEY` | `api_key` |
| Concurrency | `--concurrency` / `-c` | `IPINFO_CONCURRENCY` | `concurrency` |

Config file location: `/etc/ipinfo/config.yaml`

Example config:
```yaml
api_key: "your_key_here"
concurrency: 5
```

## API

Endpoint: `https://pro.ip-api.com/json/{ip}?fields={fields}&key={api_key}`

Fields requested:
`query,status,message,countryCode,regionName,city,isp,proxy,hosting,mobile`

## CLI flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--format` | `-f` | `table` | Output format: `table`, `summary`, `json`, `csv` |
| `--file` | | - | Input file, one IP per line |
| `--api-key` | | - | API key (overrides config/env) |
| `--concurrency` | `-c` | `5` | Parallel requests |
| `--no-cache` | | false | Disable /24 subnet cache |
| `--help` | `-h` | - | Help |

## Input sources (priority order)

1. CLI arguments: `ipinfo 1.2.3.4 5.6.7.8`
2. `--file` flag: `ipinfo --file ips.txt`
3. stdin/pipe: `cat ips.txt | ipinfo`

Empty lines and lines starting with `#` are skipped.
Whitespace and `\r` are trimmed.

## /24 subnet cache

- Key: first 3 octets (e.g. `213.96.49`)
- On cache hit: return cached JSON with `.query` replaced by actual IP
- Cache is in-memory, per-run only
- Disable with `--no-cache`
- Cache hits logged to stderr: `(cached: 1.2.3.4 → 1.2.3.x/24)`

## Output formats

### table (default)
Pipe-separated, with header row and separator row. Suitable for pasting into Claude or markdown renderers.

```
| IP | Country | Region | City | ISP | Proxy | Hosting | Mobile |
|----|---------|--------|------|-----|-------|---------|--------|
| 8.8.8.8 | US | California | Mountain View | Google LLC | false | true | false |
```

### summary
Tab-separated, one line per IP. Suitable for terminal review.

```
8.8.8.8	US	California	Mountain View	Google LLC	proxy=false	hosting=true	mobile=false
```

### json
Raw JSON from API, one object per line (NDJSON).

```json
{"query":"8.8.8.8","countryCode":"US",...}
```

### csv
Comma-separated with header row. Suitable for Google Sheets import.

```
ip,country,region,city,isp,proxy,hosting,mobile
8.8.8.8,US,California,Mountain View,Google LLC,false,true,false
```

## Error handling

- API errors: print to stderr, continue with next IP
- Network errors: print to stderr, continue
- Invalid IP format: print to stderr, skip
- Non-success API status: print error row in chosen format

## Exit codes

- `0` – all IPs processed successfully
- `1` – one or more errors occurred
- `2` – usage/config error

## Project structure

```
ipinfo/
├── main.go
├── cmd/
│   └── root.go
├── internal/
│   ├── api/
│   │   └── ipapi.go       # HTTP client, structs
│   ├── cache/
│   │   └── cache.go       # /24 subnet cache
│   ├── format/
│   │   └── format.go      # table/summary/json/csv formatters
│   └── input/
│       └── input.go       # args/file/stdin reader
├── go.mod
├── go.sum
└── README.md
```

## Versioning

- Semantic Versioning (semver)
- Conventional Commits
- Keep a Changelog
- GitHub Flow (main + rc tags)
