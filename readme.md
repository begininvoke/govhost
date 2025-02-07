# GoVHost - Virtual Host Scanner

GoVHost is a fast and efficient virtual host scanner written in Go that helps security researchers and penetration testers discover virtual hosts on IP addresses.

## Features

- Multi-threaded scanning for improved performance
- Support for both HTTP and HTTPS protocols
- Customizable status code matching/filtering
- Multiple output formats (JSON, CSV, text)
- SSL certificate verification bypass option
- Configurable request timeout
- Directory auto-creation for output files

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
govhost -ip <target_ip> -domains <domains_file> [options]
```

### Basic Example
```bash
govhost -ip 192.168.1.100 -domains domains.txt
```

### Advanced Example
```bash
govhost -ip 192.168.1.100 -domains domains.txt \
  -threads 10 \
  -timeout 5 \
  -match "200,302" \
  -ignoreCert \
  -f json \
  -o results/output.json \
  -v
```

### Output Behavior
- Only successful requests with matching status codes are included in the output
- Failed requests and non-matching status codes are filtered out
- If no matches are found:
  - CSV output will contain only the header
  - JSON output will be an empty array
  - Text output will be empty

### Options
| Flag         | Description                                      | Default |
|--------------|--------------------------------------------------|---------|
| `-ip`        | Target IP address to scan (required)            |         |
| `-domains`   | Path to file containing domains (required)       |         |
| `-threads`   | Number of concurrent threads                     | 5       |
| `-timeout`   | HTTP request timeout in seconds                  | 10      |
| `-match`     | Comma-separated status codes to include          | "200"   |
| `-ignoreCert`| Ignore SSL certificate verification errors       | false   |
| `-f`         | Output format (json, csv, or text)               | "text"  |
| `-o`         | Output file path                                 | ""      |
| `-v`         | Verbose mode (show all requests)                 | false   |

## Contributing
Contributions are welcome! Please open an issue or submit a pull request.

## License
[MIT License](LICENSE)
