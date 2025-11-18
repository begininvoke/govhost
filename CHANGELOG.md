# Changelog

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

