# GoVHost v2.3 - Feature Summary

## Latest Updates (v2.3)

### 🔄 **Redirect Control** (NEW in v2.3)
**What**: Control whether HTTP redirects (301, 302) are automatically followed

**Flag**: `-noRedirect`

**Why**: 
- See actual redirect responses instead of final destinations
- Find redirect loops and misconfigurations
- Map redirect chains
- Discover virtual hosts that only return redirects

**Example**:
```bash
# Default - follows redirects
$ ./govhost -ip 1.1.1.1 -domain cloudflare.com -match 200,301,302
✓ FOUND: http://cloudflare.com (IP: 1.1.1.1) - Status: 200

# With -noRedirect - stops at redirect
$ ./govhost -ip 1.1.1.1 -domain cloudflare.com -noRedirect -match 200,301,302
✓ FOUND: http://cloudflare.com (IP: 1.1.1.1) - Status: 301
```

### 📊 **Better Default Match Codes**
- **Old default**: `200` only
- **New default**: `200,301,302`
- Catches more virtual host configurations automatically

---

## Previous Updates (v2.2)

### 🎯 Real-Time Output
**What**: Virtual hosts are displayed immediately when discovered during scanning

**Why**: No need to wait for entire scan to complete before seeing results

**Example**:
```bash
$ ./govhost -ip 192.168.1.100 -domains domains.txt
Scanning 1 IP(s) with 5 domain(s)...
✓ FOUND: http://admin.example.com (IP: 192.168.1.100) - Status: 200
✓ FOUND: https://api.example.com (IP: 192.168.1.100) - Status: 200
✓ FOUND: http://dev.example.com (IP: 192.168.1.100) - Status: 200
```

### ⚡ Parallel Execution Visible
**What**: See requests being processed concurrently in verbose mode

**Why**: Understand scan progress and confirm parallel testing is working

**Example**:
```bash
$ ./govhost -ip 192.168.1.1-5 -domain test.com -v -threads 20
[3/10] Checking https://test.com (IP: 192.168.1.2)
[1/10] Checking http://test.com (IP: 192.168.1.1)
[5/10] Checking https://test.com (IP: 192.168.1.3)
[2/10] Checking https://test.com (IP: 192.168.1.1)
```
*Note: Out-of-order numbering proves parallel execution*

### 📊 Enhanced Verbose Mode
**What**: Beautiful startup banner with complete scan configuration

**Features**:
- Target IP count
- Domain count  
- Total requests
- Thread count
- Timeout configuration

**Example**:
```bash
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
GoVHost - Virtual Host Discovery
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Target IPs: 254
Domains: 100
Total requests: 50,800 (http+https per domain per IP)
Threads: 50
Timeout: 5s (Connection: 30s)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

### ✅ Scan Summary
**What**: Clear completion message with results count

**Example**:
```bash
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Scan Complete! Found 3 matching virtual host(s)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## All Features (v1.0 - v2.2)

### IP Input Formats
- ✅ Single IP: `192.168.1.100`
- ✅ IP Range (Full): `192.168.1.1-192.168.1.50`
- ✅ IP Range (Short): `192.168.1.1-50`
- ✅ CIDR: `192.168.1.0/24`

### Domain Input
- ✅ Single domain: `-domain example.com`
- ✅ Domain list file: `-domains domains.txt`
- ✅ Wordlist enumeration: `-wordlist subdomains.txt`

### Timeout System (v2.1)
- ✅ Connection timeout: 30 seconds (fixed)
- ✅ TLS handshake: 10 seconds (fixed)
- ✅ Request timeout: Configurable (default 10s)
- ✅ Timeout notifications in verbose mode

### Output & Display (v2.2)
- ✅ Real-time discovery notifications
- ✅ Parallel execution indicators
- ✅ Progress tracking `[current/total]`
- ✅ Scan summary with results count
- ✅ Multiple formats: JSON, CSV, Text

### Performance
- ✅ Multi-threaded (configurable)
- ✅ Parallel HTTP/HTTPS requests
- ✅ Connection pooling
- ✅ Efficient timeout handling

---

## Quick Start Examples

### Basic Scan
```bash
./govhost -ip 192.168.1.100 -domain example.com
```

### Fast Local Network Scan
```bash
./govhost -ip 192.168.1.0/24 -domains domains.txt -threads 30 -timeout 3
```

### Subdomain Enumeration
```bash
./govhost -ip 10.0.0.1 -domain example.com -wordlist subs.txt -threads 20
```

### Large-Scale Scan with Progress
```bash
./govhost -ip 10.0.0.0/16 -domains domains.txt -wordlist subs.txt \
          -threads 50 -timeout 5 -f json -o results.json -v
```

### Internet Target
```bash
./govhost -ip 1.2.3.4 -domains domains.txt -timeout 15 -v
```

---

## Performance Tips

### Thread Count
- **Local networks**: 20-50 threads
- **Internet targets**: 5-20 threads
- **Slow networks**: 5-10 threads

### Timeout Settings
- **Fast networks**: `-timeout 3`
- **Normal**: `-timeout 10` (default)
- **Slow/congested**: `-timeout 15-20`

### Large Scans
```bash
# Optimize for speed
./govhost -ip 192.168.0.0/16 -domains domains.txt \
          -threads 50 -timeout 5 -f json -o results.json

# Monitor progress
tail -f results.json  # In another terminal
```

---

## Output Modes

### Standard Mode (Default)
- Real-time discoveries: `✓ FOUND: ...`
- Scan summary
- Results list

### Verbose Mode (`-v`)
- Startup configuration banner
- Progress tracking
- Request-by-request details
- Timeout notifications
- Scan summary

### Silent Mode (Output to File)
```bash
./govhost -ip 192.168.1.0/24 -domains domains.txt -o results.json
# Real-time finds still show on screen
# Final results save to file
```

---

## Troubleshooting

### Not finding expected hosts?
1. Increase timeout: `-timeout 15`
2. Check with verbose: `-v`
3. Verify IP/domain combination manually

### Scan too slow?
1. Increase threads: `-threads 30`
2. Decrease timeout: `-timeout 5`
3. Check network conditions

### Too many timeouts?
1. Increase timeout: `-timeout 20`
2. Reduce threads: `-threads 5`
3. Check target availability

---

## What's Next?

GoVHost is feature-complete for most virtual host discovery scenarios. Possible future enhancements:

- Custom HTTP headers
- Proxy support
- Rate limiting options
- Resume capability for large scans
- HTML report generation

---

## Support

For issues, feature requests, or contributions:
- Check README.md for full documentation
- Review CHANGELOG.md for version history
- See TIMEOUT_GUIDE.md for timeout details

Happy hunting! 🎯

