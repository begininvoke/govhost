# Changelog

## Version 2.3 - Redirect Control & Enhanced Defaults

### New Features

#### Redirect Control (`-noRedirect`)
- **Disable automatic redirect following**: See actual 301/302 responses instead of final destination
- **Use case**: Find redirect loops, map redirect chains, discover redirect-only virtual hosts
- **Default behavior**: Follows redirects (standard HTTP client behavior)

**Example**:
```bash
# See redirect responses
./govhost -ip 192.168.1.100 -domain example.com -noRedirect -match 301,302

# Compare with default (follows redirects)
./govhost -ip 192.168.1.100 -domain example.com -match 200,301,302
```

### Improvements

#### Better Default Match Codes
- **Old default**: `200` only
- **New default**: `200,301,302` 
- Reason: Catch more virtual host configurations including redirects
- Can still customize with `-match` flag

### Technical Details

The `-noRedirect` flag sets a custom `CheckRedirect` function that returns `http.ErrUseLastResponse`, preventing the HTTP client from automatically following Location headers.

---

## Version 2.2 - Real-Time Output & Parallel Execution

### New Features

#### Real-Time Virtual Host Discovery
- **Instant feedback**: Virtual hosts are displayed immediately when discovered
- Format: `✓ FOUND: protocol://domain (IP: address) - Status: code`
- No need to wait for scan completion to see results
- Results appear as they are found during parallel execution

#### Enhanced Verbose Mode
- **Progress tracking**: Shows `[current/total]` for each request
- **Startup banner**: Displays scan configuration
  - Target IPs count
  - Domains count
  - Total requests
  - Thread count
  - Timeout settings
- **Request tracking**: See which requests are being processed in real-time

#### Improved Scan Summary
- Beautiful completion banner with separator lines
- Total count of discovered virtual hosts
- Clear visual separation between scan progress and results

#### Parallel Execution Indicators
- Requests execute in parallel (visible in verbose mode)
- Out-of-order completion shows true concurrent processing
- Example: `[2/8]`, `[3/8]`, `[1/8]` indicates parallel execution

### User Experience Improvements
- Clear visual hierarchy with Unicode box drawing characters
- Immediate feedback on discoveries
- Better progress visibility during large scans
- Professional-looking output format

### Examples

```bash
# Standard mode - see finds in real-time
$ ./govhost -ip 192.168.1.100 -domains domains.txt
Scanning 1 IP(s) with 5 domain(s)...
✓ FOUND: http://admin.example.com (IP: 192.168.1.100) - Status: 200
✓ FOUND: https://api.example.com (IP: 192.168.1.100) - Status: 200
```

```bash
# Verbose mode - detailed progress
$ ./govhost -ip 192.168.1.1-5 -domain example.com -v -threads 20
[3/10] Checking https://example.com (IP: 192.168.1.2)
[1/10] Checking http://example.com (IP: 192.168.1.1)
[5/10] Checking https://example.com (IP: 192.168.1.3)
✓ FOUND: http://example.com (IP: 192.168.1.1) - Status: 200
```

---

## Version 2.1 - Enhanced Timeout Handling

### New Features

#### Comprehensive Timeout System
- **Connection Timeout**: 30 seconds (fixed) - Prevents hanging on unreachable IPs
- **TLS Handshake Timeout**: 10 seconds (fixed) - Ensures SSL/TLS doesn't hang
- **Request Timeout**: Configurable via `-timeout` flag (default: 10 seconds)
- **Idle Connection Timeout**: 30 seconds (fixed)
- **Verbose Timeout Reporting**: Shows ⏱ indicator when servers timeout

#### Improved Connection Handling
- Better connection pooling with `MaxIdleConns` set to 100
- `DisableKeepAlives` for cleaner connection management
- Proper dial context with 30-second connection establishment timeout

### Technical Improvements
- Enhanced HTTP transport configuration
- Better error reporting for timeout conditions
- Optimized connection handling for large-scale scans

### Examples

```bash
# Quick scan with short timeout
./govhost -ip 192.168.1.0/24 -domain example.com -timeout 3 -threads 20

# Verbose mode shows timeout events
./govhost -ip 192.0.2.1 -domain test.com -v
# Output: ⏱ Timeout: https://test.com (IP: 192.0.2.1)
```

---

## Version 2.0 - Multi-IP and Enhanced Features

### New Features

#### 1. Multiple IP Input Formats
- **Single IP**: `192.168.1.100`
- **IP Range (Full)**: `192.168.1.1-192.168.1.50`
- **IP Range (Short)**: `192.168.1.1-50`
- **CIDR Notation**: `192.168.1.0/24`

#### 2. Flexible Domain Input
- **Single Domain**: Use `-domain` flag for testing a single domain
- **Domain List File**: Use `-domains` flag for testing multiple domains from a file
- Both options support wordlist combination

#### 3. Wordlist Support
- New `-wordlist` flag for subdomain enumeration
- Combines wordlist entries with domain(s) to generate subdomains
- Example: `www` + `example.com` = `www.example.com`

#### 4. Comprehensive Help System
- Detailed help with `-h`, `--help`, or `help`
- Usage examples for common scenarios
- File format specifications
- IP format documentation

### Improvements

#### Better Error Messages
- Clear error messages for missing required arguments
- Helpful format examples for invalid IP inputs
- Validation for conflicting arguments

#### Enhanced Verbosity
- Shows total IPs and domains being scanned
- Displays total expected requests
- Better progress tracking

### Examples

```bash
# Single IP with single domain
govhost -ip 192.168.1.100 -domain example.com

# IP range with domain list
govhost -ip 192.168.1.1-50 -domains domains.txt

# CIDR with wordlist
govhost -ip 10.0.0.0/24 -domain example.com -wordlist subs.txt

# Comprehensive scan
govhost -ip 172.16.0.0/24 -domains domains.txt -wordlist subs.txt \
        -threads 20 -match 200,301,302 -f json -o results.json -v
```

### Breaking Changes
- None. All previous commands remain compatible.
- The `-domains` flag works as before, now with optional `-domain` alternative.

### Technical Details

#### New Functions
- `parseIPInput()` - Parses IP, range, or CIDR notation
- `parseCIDR()` - Handles CIDR notation with network/broadcast exclusion
- `parseIPRange()` - Handles IP range with full or short notation
- `incrementIP()` - IP address increment helper
- `readWordlistFromFile()` - Reads wordlist from file
- `combineWordlistWithDomains()` - Combines wordlist with domains
- `printUsage()` - Comprehensive help system

#### Safety Features
- IP range size limit (max 65,536 IPs) to prevent memory issues
- CIDR network and broadcast address exclusion
- Proper error handling for invalid IP formats
- Input validation for conflicting flags

---

## Version 1.0 - Initial Release

### Features
- Single IP scanning
- Domain list file support
- Multi-threaded scanning
- HTTP and HTTPS protocol support
- Status code filtering
- Multiple output formats (JSON, CSV, text)
- SSL certificate verification bypass
- Configurable timeouts and threads

