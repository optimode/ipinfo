# ipinfo

IP geolocation CLI tool using the [ip-api.com](https://ip-api.com) Pro API.

## Installation

```bash
go install github.com/optimode/ipinfo@latest
```

Or build from source:

```bash
git clone ...
cd ipinfo
go build -o ipinfo .
sudo mv ipinfo /usr/local/bin/
```

## Configuration

Create `/etc/ipinfo/config.yaml`:

```yaml
api_key: "your_key_here"
concurrency: 5
```

Or use environment variables:

```bash
export IPINFO_API_KEY="your_key_here"
export IPINFO_CONCURRENCY=10
```

## Usage

```bash
# Single IP
ipinfo 8.8.8.8

# Multiple IPs
ipinfo 8.8.8.8 1.1.1.1 9.9.9.9

# From file
ipinfo --file ips.txt

# From pipe
cat ips.txt | ipinfo
grep -oP 'rip=\K[0-9.]+' /var/log/dovecot.log | sort -u | ipinfo

# Output formats
ipinfo -f table 8.8.8.8       # default, pipe-separated markdown table
ipinfo -f summary 8.8.8.8     # tab-separated single line
ipinfo -f json 8.8.8.8        # raw JSON (NDJSON)
ipinfo -f csv 8.8.8.8         # CSV with header

# Concurrency
ipinfo -c 10 --file large_list.txt

# Disable /24 subnet cache
ipinfo --no-cache 8.8.8.8
```

## Output examples

### table (default)
```
| IP | Country | Region | City | ISP | Proxy | Hosting | Mobile |
|----|---------|--------|------|-----|-------|---------|--------|
| 8.8.8.8 | US | California | Mountain View | Google LLC | false | true | false |
```

### summary
```
8.8.8.8	US	California	Mountain View	Google LLC	proxy=false	hosting=true	mobile=false
```

### csv
```
ip,country,region,city,isp,proxy,hosting,mobile
8.8.8.8,US,California,Mountain View,Google LLC,false,true,false
```

## /24 Subnet cache

By default, if multiple IPs fall within the same /24 subnet, only the first one
is looked up via the API. The rest return cached data with the actual IP substituted.
Cache hits are logged to stderr.

Disable with `--no-cache`.

## Changelog

### v0.1.0
- Initial release
- table, summary, json, csv output formats
- CLI args, --file, stdin input
- /24 subnet cache
- Configurable concurrency
- /etc/ipinfo/config.yaml + env var support
