# GoVHost Redirect Handling Guide

## Overview

HTTP redirects (301, 302, 307, 308) are common responses from web servers that tell the client to request a different URL. By default, GoVHost follows these redirects to reach the final destination, but you can disable this with the `-noRedirect` flag.

## Redirect Behavior

### Default Behavior (Follows Redirects)

Without any special flags, GoVHost follows HTTP redirects automatically:

```bash
$ ./govhost -ip 1.1.1.1 -domain cloudflare.com -match "200,301,302"
✓ FOUND: http://cloudflare.com (IP: 1.1.1.1) - Status: 200
```

**What happens**:
1. Request sent to `http://cloudflare.com` at IP 1.1.1.1
2. Server responds with 301 redirect to `https://cloudflare.com`
3. GoVHost automatically follows the redirect
4. Final response is 200 OK
5. **Status displayed: 200** (the final destination)

### With `-noRedirect` Flag

When you disable automatic redirect following:

```bash
$ ./govhost -ip 1.1.1.1 -domain cloudflare.com -noRedirect -match "200,301,302"
✓ FOUND: http://cloudflare.com (IP: 1.1.1.1) - Status: 301
```

**What happens**:
1. Request sent to `http://cloudflare.com` at IP 1.1.1.1
2. Server responds with 301 redirect
3. GoVHost **stops here** and doesn't follow the redirect
4. **Status displayed: 301** (the redirect itself)

## When to Use `-noRedirect`

### 1. Finding Redirect-Only Virtual Hosts

Some virtual hosts are configured to only redirect traffic:

```bash
# Find all hosts that redirect
./govhost -ip 192.168.1.0/24 -domains domains.txt -noRedirect -match 301,302
```

**Use case**: Discovering misconfigured or redirect-only subdomains

### 2. Mapping Redirect Chains

See which hosts redirect where:

```bash
# First: Find redirects
./govhost -ip 192.168.1.100 -domains domains.txt -noRedirect -match 301,302 -o redirects.csv

# Then: Analyze redirect patterns
cat redirects.csv | grep 301
```

**Use case**: Understanding site architecture and redirect patterns

### 3. Detecting Redirect Loops

Identify problematic redirect configurations:

```bash
# Look for suspicious redirect patterns
./govhost -ip 192.168.1.0/24 -domain example.com -wordlist subs.txt \
          -noRedirect -match 301,302,307,308 -v
```

**Use case**: Finding misconfigurations in load balancers or proxies

### 4. Performance Optimization

Following redirects takes extra time and network requests:

```bash
# Faster scans by not following redirects
./govhost -ip 192.168.0.0/16 -domains domains.txt -noRedirect \
          -match 200,301,302 -threads 50 -timeout 3
```

**Use case**: Large-scale scans where speed is important

## Common Redirect Status Codes

| Code | Name | Description |
|------|------|-------------|
| 301 | Moved Permanently | Resource has permanently moved to a new URL |
| 302 | Found (Moved Temporarily) | Resource is temporarily at a different URL |
| 303 | See Other | Response can be found at another URL (GET only) |
| 307 | Temporary Redirect | Like 302 but preserves HTTP method |
| 308 | Permanent Redirect | Like 301 but preserves HTTP method |

## Practical Examples

### Example 1: Comprehensive Discovery

Find all response types including redirects:

```bash
./govhost -ip 192.168.1.100 -domains domains.txt \
          -match "200,201,301,302,401,403,404" \
          -f json -o results.json
```

This follows redirects by default, so you'll see final status codes.

### Example 2: Redirect-Only Scan

Focus specifically on redirect responses:

```bash
./govhost -ip 192.168.1.0/24 -domain example.com \
          -wordlist subdomains.txt \
          -noRedirect -match 301,302,307,308 \
          -threads 20 -o redirects.txt
```

### Example 3: HTTP to HTTPS Detection

Find which hosts redirect from HTTP to HTTPS:

```bash
# Step 1: Check redirects
./govhost -ip 192.168.1.100 -domains domains.txt -noRedirect -match 301,302

# Look for patterns like:
# http://example.com → 301 (likely redirecting to https://)
```

### Example 4: Comparing Behaviors

See the difference between redirect and final destination:

```bash
echo "=== With Redirects (Default) ==="
./govhost -ip 1.1.1.1 -domain cloudflare.com -match 200,301,302

echo ""
echo "=== Without Redirects ==="
./govhost -ip 1.1.1.1 -domain cloudflare.com -noRedirect -match 200,301,302
```

## Technical Details

### How It Works

When `-noRedirect` is enabled, GoVHost sets a custom `CheckRedirect` function:

```go
client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
    return http.ErrUseLastResponse
}
```

This tells the HTTP client to return the redirect response instead of following it.

### Redirect Chain Limits

By default (without `-noRedirect`), Go's HTTP client follows up to 10 redirects. After that, it returns an error. This prevents infinite redirect loops.

### Performance Impact

**Following redirects (default)**:
- Slower: Each redirect adds another HTTP request
- More network traffic: Multiple round trips
- Shows final destination: Useful for understanding where requests end up

**Not following redirects (`-noRedirect`)**:
- Faster: Only one HTTP request per test
- Less network traffic: Single round trip
- Shows redirect status: Useful for understanding server configuration

## Output Formats with Redirects

### Text Output

```bash
$ ./govhost -ip 1.1.1.1 -domain cloudflare.com -noRedirect -match 301
✓ FOUND: http://cloudflare.com (IP: 1.1.1.1) - Status: 301
```

### JSON Output

```bash
$ ./govhost -ip 1.1.1.1 -domain cloudflare.com -noRedirect -match 301 -f json
[
  {
    "domain": "cloudflare.com",
    "ip": "1.1.1.1",
    "protocol": "http",
    "status_code": 301
  }
]
```

### CSV Output

```bash
$ ./govhost -ip 1.1.1.1 -domain cloudflare.com -noRedirect -match 301 -f csv
domain,ip,protocol,status_code,error
cloudflare.com,1.1.1.1,http,301,
```

## Best Practices

### 1. Start Broad, Then Focus

```bash
# First: See what's there (including final destinations)
./govhost -ip 192.168.1.100 -domains domains.txt -match "200,301,302"

# Then: Focus on redirects if interesting patterns emerge
./govhost -ip 192.168.1.100 -domains domains.txt -noRedirect -match 301,302
```

### 2. Combine with Status Code Filtering

```bash
# Only see 301 (permanent) redirects
./govhost -ip 192.168.1.0/24 -domains domains.txt -noRedirect -match 301

# Only see 302 (temporary) redirects  
./govhost -ip 192.168.1.0/24 -domains domains.txt -noRedirect -match 302
```

### 3. Use Verbose Mode for Debugging

```bash
# See what's happening with each redirect
./govhost -ip 192.168.1.100 -domain example.com -noRedirect -match 301,302 -v
```

## Troubleshooting

### Problem: Not seeing any redirects

**Solution**: Make sure you're including redirect status codes in `-match`:

```bash
# Wrong - won't show redirects
./govhost -ip 1.1.1.1 -domain example.com -noRedirect -match 200

# Correct - includes redirect codes
./govhost -ip 1.1.1.1 -domain example.com -noRedirect -match 200,301,302
```

### Problem: Getting different results than expected

**Solution**: Try both modes and compare:

```bash
# See redirect
./govhost -ip IP -domain example.com -noRedirect -match 301,302 -v

# See final destination
./govhost -ip IP -domain example.com -match 200 -v
```

### Problem: Some hosts not responding

**Solution**: Increase timeout when checking for redirects:

```bash
./govhost -ip 192.168.1.0/24 -domains domains.txt -noRedirect \
          -match 301,302 -timeout 15
```

## Summary

| Scenario | Use Default | Use `-noRedirect` |
|----------|------------|-------------------|
| General discovery | ✅ | |
| Fast scanning | | ✅ |
| Find redirect-only hosts | | ✅ |
| Map redirect chains | | ✅ |
| Detect redirect loops | | ✅ |
| See final destinations | ✅ | |
| Comprehensive coverage | ✅ | |
| Performance optimization | | ✅ |

**Default**: New default match codes are `200,301,302` to catch both content and redirects regardless of mode.

**Remember**: The `-noRedirect` flag gives you more control and insight into server behavior, but following redirects (default) often makes more sense for discovering actual virtual hosts.

