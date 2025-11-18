# GoVHost - Virtual Host Scanner

GoVHost is a fast and efficient virtual host scanner written in Go that helps security researchers and penetration testers discover virtual hosts on IP addresses.

## Features

- **Multi-IP scanning** - Support for single IPs, IP ranges, and CIDR notation
- **Flexible domain input** - Single domain or domain list file
- **Wordlist-based subdomain enumeration** - Combine wordlists with domains
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
- Only successful requests with matching status codes are included in the output
- Failed requests and non-matching status codes are filtered out
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

## Contributing
Contributions are welcome! Please open an issue or submit a pull request.

## License
[MIT License](LICENSE)
