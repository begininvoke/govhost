# GoVHost - Virtual Host Scanner

GoVHost is a fast and efficient virtual host scanner written in Go that helps security researchers and penetration testers discover virtual hosts on IP addresses.

## Features

- **Multi-IP scanning** - Support for single IPs, IP ranges, and CIDR notation
- **Flexible domain input** - Single domain or domain list file
- **Wordlist-based subdomain enumeration** - Combine wordlists with domains
- **Real-time output** - See discovered virtual hosts as they are found
- **Parallel testing** - Concurrent requests for fast scanning
- Multi-threaded scanning for improved performance
- Support for both HTTP and HTTPS protocols
- Customizable status code matching/filtering
- Multiple output formats (JSON, CSV, text)
- SSL certificate verification bypass option
- Configurable request timeout
- Directory auto-creation for output files
- **Comprehensive help system** with examples and usage guides

## Installation

1. Make sure you have Go installed (version 1.20 or higher recommended)
2. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/govhost.git
   cd govhost
   ```
3. Build the binary:
   ```bash
   go build -o govhost main.go
   ```
4. Move the binary to your PATH (optional):
   ```bash
   sudo mv govhost /usr/local/bin/
   ```

## Usage

```bash
govhost -ip <IP/RANGE/CIDR> -domain <domain> [options]
govhost -ip <IP/RANGE/CIDR> -domains <file> [options]
```

### Basic Examples

**Single IP with single domain:**
```bash
govhost -ip 192.168.1.100 -domain example.com
```

**Single IP with domain list:**
```bash
govhost -ip 192.168.1.100 -domains domains.txt
```

**IP range with domain list:**
```bash
govhost -ip 192.168.1.1-192.168.1.50 -domains domains.txt
```

**Short IP range notation:**
```bash
govhost -ip 192.168.1.1-50 -domains domains.txt
```

**CIDR notation with single domain:**
```bash
govhost -ip 10.0.0.0/24 -domain example.com
```

**Wordlist-based subdomain enumeration:**
```bash
govhost -ip 192.168.1.100 -domain example.com -wordlist subdomains.txt
```

### Advanced Examples

**Comprehensive scan with all features:**
```bash
govhost -ip 172.16.0.0/24 -domains domains.txt -wordlist subs.txt \
  -threads 20 \
  -timeout 5 \
  -match "200,301,302,403,404" \
  -f json \
  -o results/output.json \
  -v
```

**Multiple IPs with wordlist:**
```bash
govhost -ip 192.168.1.1-10 -domain example.com -wordlist wordlist.txt \
  -match "200,403" \
  -threads 10
```

### Output Examples

**Standard mode (real-time findings):**
```bash
$ ./govhost -ip 192.168.1.100 -domains domains.txt
Scanning 1 IP(s) with 5 domain(s)...
✓ FOUND: http://admin.example.com (IP: 192.168.1.100) - Status: 200
✓ FOUND: https://api.example.com (IP: 192.168.1.100) - Status: 200

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Scan Complete! Found 2 matching virtual host(s)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

http://admin.example.com (IP: 192.168.1.100) - Status: 200
https://api.example.com (IP: 192.168.1.100) - Status: 200
```

**Verbose mode (detailed progress):**
```bash
$ ./govhost -ip 192.168.1.100 -domain example.com -v

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
GoVHost - Virtual Host Discovery
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Target IPs: 1
Domains: 1
Total requests: 2 (http+https per domain per IP)
Threads: 5
Timeout: 10s (Connection: 30s)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

[1/2] Checking https://example.com (IP: 192.168.1.100)
[2/2] Checking http://example.com (IP: 192.168.1.100)
✓ FOUND: http://example.com (IP: 192.168.1.100) - Status: 200

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Scan Complete! Found 1 matching virtual host(s)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### IP Format Options

GoVHost supports multiple IP input formats:

| Format | Example | Description |
|--------|---------|-------------|
| Single IP | `192.168.1.100` | Scan a single IP address |
| IP Range (Full) | `192.168.1.1-192.168.1.50` | Scan a range of IPs |
| IP Range (Short) | `192.168.1.1-50` | Short notation for IP range |
| CIDR | `192.168.1.0/24` | Scan entire subnet (excludes network/broadcast) |

### Wordlist Feature

The `-wordlist` option enables subdomain enumeration by combining words with your domains:

- **Format**: One subdomain per line (e.g., `www`, `admin`, `api`)
- **Combination**: Each word is prepended to each domain
- **Example**: Domain `example.com` + wordlist `[www, admin]` = Tests `www.example.com` and `admin.example.com`

**Sample wordlist file:**
```
www
admin
api
dev
staging
blog
mail
ftp
```

### Output Behavior

#### Real-Time Discovery
- **Instant notifications**: Virtual hosts are displayed immediately when found
- Format: `✓ FOUND: protocol://domain (IP: address) - Status: code`
- Results appear as scanning progresses (parallel execution)

#### Verbose Mode (`-v`)
- Detailed scan information with progress counter `[current/total]`
- Shows startup banner with scan configuration
- Displays timeout events with ⏱ indicator
- Real-time request tracking

#### Final Summary
- Scan completion banner
- Total count of discovered virtual hosts
- Full results list in chosen format

#### Result Filtering
- Only successful requests with matching status codes are included in the output
- Failed requests and non-matching status codes are filtered out
- **Timeout handling**: Servers that don't respond within 30 seconds will timeout
  - Connection timeout: 30 seconds (fixed)
  - Request timeout: Configurable with `-timeout` flag (default: 10 seconds)
  - TLS handshake timeout: 10 seconds (fixed)
  - In verbose mode (`-v`), timeout events are displayed with ⏱ indicator
- If no matches are found:
  - CSV output will contain only the header
  - JSON output will be null or empty array
  - Text output will be empty

### Options
| Flag         | Description                                      | Default |
|--------------|--------------------------------------------------|---------|
| `-ip`        | Target IP, IP range, or CIDR (required)         |         |
| `-domain`    | Single domain to test                            |         |
| `-domains`   | Path to file containing domains                  |         |
| `-wordlist`  | Path to wordlist for subdomain enumeration       |         |
| `-threads`   | Number of concurrent threads                     | 5       |
| `-timeout`   | HTTP request timeout in seconds                  | 10      |
|              | Note: Connection timeout is fixed at 30 seconds  |         |
| `-match`     | Comma-separated status codes to include          | "200"   |
| `-ignoreCert`| Ignore SSL certificate verification errors       | true    |
| `-f`         | Output format (json, csv, or text)               | "text"  |
| `-o`         | Output file path                                 | ""      |
| `-v`         | Verbose mode (show all requests)                 | false   |

**Note**: Either `-domain` or `-domains` is required, but not both.

### Help System

Get comprehensive help with examples:
```bash
govhost -h
# or
govhost --help
# or
govhost help
```

## Technical Details

### Timeout Configuration

The tool implements multiple timeout layers to handle unresponsive servers:

1. **Connection Timeout**: 30 seconds (fixed)
   - Time allowed to establish TCP connection to target
   - Prevents hanging on unreachable or firewalled IPs

2. **TLS Handshake Timeout**: 10 seconds (fixed)
   - Time allowed for SSL/TLS negotiation
   - Ensures HTTPS connections don't hang indefinitely

3. **Request Timeout**: Configurable (default: 10 seconds)
   - Set via `-timeout` flag
   - Total time allowed for complete HTTP request/response cycle
   - Recommended: 5-15 seconds for most use cases

4. **Idle Connection Timeout**: 30 seconds (fixed)
   - Time to keep idle connections in pool

### Performance Tips

- **Threads**: Increase `-threads` for faster scans of large IP ranges
  - Default: 5 (conservative)
  - Recommended: 10-20 for local networks
  - Caution: Too many threads may trigger rate limiting or exhaust resources

- **Timeout**: Adjust based on network conditions
  - Fast local networks: `-timeout 3`
  - Internet targets: `-timeout 10` (default)
  - Slow/congested networks: `-timeout 15`

- **Large Scans**: For /24 or larger subnets
  - Use higher thread count: `-threads 20`
  - Consider shorter timeout: `-timeout 5`
  - Save to file: `-o results.json`
  - Monitor with verbose mode: `-v`

## Contributing
Contributions are welcome! Please open an issue or submit a pull request.

## License
[MIT License](LICENSE)
