# GoVHost Timeout Reference Guide

## Overview

GoVHost implements multiple timeout layers to handle various network conditions and prevent hanging on unresponsive servers.

## Timeout Layers

### 1. Connection Timeout (30 seconds - Fixed)
**What it does**: Maximum time allowed to establish a TCP connection to the target IP.

**When it triggers**:
- Target IP is unreachable
- Firewall blocking connection
- Network routing issues
- Host is down

**Example**:
```bash
./govhost -ip 192.168.1.100 -domain example.com -v
# Output if connection fails:
# ⏱ Timeout: https://example.com (IP: 192.168.1.100)
```

### 2. TLS Handshake Timeout (10 seconds - Fixed)
**What it does**: Maximum time for SSL/TLS negotiation on HTTPS connections.

**When it triggers**:
- SSL/TLS misconfiguration
- Certificate issues
- Slow cryptographic operations

### 3. Request Timeout (Configurable - Default: 10 seconds)
**What it does**: Total time for complete HTTP request/response cycle.

**Configuration**: Use `-timeout` flag

**Recommendations by scenario**:

| Scenario | Recommended Timeout | Command Example |
|----------|-------------------|-----------------|
| Local network | 3-5 seconds | `./govhost -ip 192.168.1.0/24 -domain test.com -timeout 3` |
| Internet targets | 10 seconds (default) | `./govhost -ip 1.2.3.4 -domain example.com` |
| Slow networks | 15-20 seconds | `./govhost -ip 1.2.3.4 -domain example.com -timeout 15` |

### 4. Idle Connection Timeout (30 seconds - Fixed)
**What it does**: How long to keep idle connections in the connection pool.

## Real-World Examples

### Fast Local Network Scan
```bash
# Scan entire /24 subnet quickly
./govhost -ip 192.168.1.0/24 -domains domains.txt \
          -timeout 3 -threads 20 -v
```

### Internet Target with Retries
```bash
# Longer timeout for remote servers
./govhost -ip 203.0.113.10 -domain example.com \
          -timeout 15 -v
```

### Large-Scale Scan
```bash
# Scan multiple subnets with optimized settings
./govhost -ip 10.0.0.0/16 -domains domains.txt \
          -wordlist subs.txt \
          -timeout 5 -threads 50 \
          -f json -o results.json
```

## Timeout Behavior

### Verbose Mode (`-v`)
When verbose mode is enabled, timeout events are displayed:

```
Scanning 3 IP(s) with 2 domain(s)
Total requests: 12 (http+https per domain per IP)
Checking https://example.com (IP: 192.0.2.1)
  ⏱ Timeout: https://example.com (IP: 192.0.2.1)
Checking http://example.com (IP: 192.0.2.1)
  ⏱ Timeout: http://example.com (IP: 192.0.2.1)
```

### Silent Mode (Default)
Without verbose mode, timeouts are silently skipped (no output).

## Performance Impact

### Timeout vs Scan Speed

| Setting | Speed | Risk | Best For |
|---------|-------|------|----------|
| Short timeout (3-5s) | Fast | May miss slow servers | Local networks, quick scans |
| Medium timeout (10s) | Balanced | Good coverage | General use (default) |
| Long timeout (15-20s) | Slow | Complete coverage | Slow networks, thorough scans |

### Thread Count Interaction

Higher thread count + shorter timeout = faster scans:

```bash
# Fast scan (may miss some hosts)
./govhost -ip 192.168.0.0/16 -domain test.com -timeout 3 -threads 50

# Thorough scan (slower but more complete)
./govhost -ip 192.168.0.0/16 -domain test.com -timeout 15 -threads 10
```

## Troubleshooting

### Issue: Too many timeouts
**Solution**: Increase timeout value
```bash
./govhost -ip TARGET -domain example.com -timeout 20
```

### Issue: Scan takes too long
**Solution**: Decrease timeout and increase threads
```bash
./govhost -ip TARGET -domain example.com -timeout 5 -threads 20
```

### Issue: Missing legitimate hosts
**Solution**: Increase timeout, check network conditions
```bash
./govhost -ip TARGET -domain example.com -timeout 15 -v
```

## Technical Details

The timeout configuration in code:

```go
client := &http.Client{
    Timeout: time.Duration(*requestTimeout) * time.Second,
    Transport: &http.Transport{
        TLSClientConfig:     &tls.Config{InsecureSkipVerify: *ignoreCert},
        DisableKeepAlives:   true,
        MaxIdleConns:        100,
        IdleConnTimeout:     30 * time.Second,
        TLSHandshakeTimeout: 10 * time.Second,
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
    },
}
```

## Summary

- **Connection timeout**: 30s (fixed) - TCP connection establishment
- **TLS timeout**: 10s (fixed) - SSL/TLS handshake
- **Request timeout**: Configurable (default 10s) - Full HTTP cycle
- **Idle timeout**: 30s (fixed) - Connection pooling
- **Verbose mode**: Shows timeout events with ⏱ indicator

