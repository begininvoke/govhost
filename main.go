package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Define a struct to hold the request result, including a field for errors
type RequestResult struct {
	Domain     string `json:"domain,omitempty"`
	IP         string `json:"ip"`
	Protocol   string `json:"protocol"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error,omitempty"`
}

func printUsage() {
	fmt.Println("govhost - Virtual Host Discovery Tool")
	fmt.Println("")
	fmt.Println("DESCRIPTION:")
	fmt.Println("  Scans IP addresses for virtual hosts by testing domains against HTTP/HTTPS")
	fmt.Println("  and checking for specific status codes. Supports single IPs, IP ranges,")
	fmt.Println("  CIDR notation, and both single domains or domain list files.")
	fmt.Println("")
	fmt.Println("USAGE:")
	fmt.Println("  govhost -ip <IP/RANGE/CIDR> -domain <DOMAIN> [OPTIONS]")
	fmt.Println("  govhost -ip <IP/RANGE/CIDR> -domains <FILE> [OPTIONS]")
	fmt.Println("")
	fmt.Println("REQUIRED ARGUMENTS:")
	fmt.Println("  -ip string")
	fmt.Println("        Target IP, IP range, or CIDR to scan")
	fmt.Println("        Examples: 192.168.1.1, 192.168.1.1-192.168.1.10, 192.168.1.0/24")
	fmt.Println("  -domain string OR -domains string")
	fmt.Println("        Single domain (-domain) or path to domain list file (-domains)")
	fmt.Println("")
	fmt.Println("OPTIONAL ARGUMENTS:")
	fmt.Println("  -wordlist string")
	fmt.Println("        Path to wordlist file for subdomain enumeration")
	fmt.Println("        Combines words with domains (e.g., www.example.com)")
	fmt.Println("  -threads int")
	fmt.Println("        Number of concurrent threads (default: 5)")
	fmt.Println("  -timeout int")
	fmt.Println("        HTTP request timeout in seconds (default: 10)")
	fmt.Println("  -match string")
	fmt.Println("        Comma-separated status codes to include (default: 200)")
	fmt.Println("        Example: -match 200,301,302,403,404")
	fmt.Println("  -ignoreCert")
	fmt.Println("        Ignore SSL certificate verification errors (default: true)")
	fmt.Println("  -f string")
	fmt.Println("        Output format: json, csv, or text (default: text)")
	fmt.Println("  -o string")
	fmt.Println("        Output file path (default: stdout)")
	fmt.Println("  -v")
	fmt.Println("        Verbose mode - show all requests and checks")
	fmt.Println("")
	fmt.Println("EXAMPLES:")
	fmt.Println("  # Single IP with single domain")
	fmt.Println("  govhost -ip 192.168.1.100 -domain example.com")
	fmt.Println("")
	fmt.Println("  # Single IP with domain list file")
	fmt.Println("  govhost -ip 192.168.1.100 -domains domains.txt")
	fmt.Println("")
	fmt.Println("  # IP range with domain list")
	fmt.Println("  govhost -ip 192.168.1.1-192.168.1.50 -domains domains.txt")
	fmt.Println("")
	fmt.Println("  # CIDR notation with wordlist")
	fmt.Println("  govhost -ip 10.0.0.0/24 -domain example.com -wordlist subdomains.txt")
	fmt.Println("")
	fmt.Println("  # Advanced scan with custom settings")
	fmt.Println("  govhost -ip 172.16.0.0/24 -domains domains.txt -wordlist subs.txt \\")
	fmt.Println("          -threads 20 -timeout 5 -match 200,301,302 -f json -o results.json -v")
	fmt.Println("")
	fmt.Println("IP FORMAT OPTIONS:")
	fmt.Println("  Single IP:    192.168.1.100")
	fmt.Println("  IP Range:     192.168.1.1-192.168.1.254")
	fmt.Println("  CIDR:         192.168.1.0/24")
	fmt.Println("")
	fmt.Println("DOMAIN FILE FORMAT:")
	fmt.Println("  One domain per line:")
	fmt.Println("  example.com")
	fmt.Println("  test.example.com")
	fmt.Println("  admin.example.com")
	fmt.Println("")
	fmt.Println("WORDLIST FILE FORMAT:")
	fmt.Println("  One subdomain per line:")
	fmt.Println("  www")
	fmt.Println("  admin")
	fmt.Println("  api")
	fmt.Println("  dev")
	fmt.Println("  staging")
}

func main() {
	// Check for help flag first
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help") {
		printUsage()
		return
	}

	ip := flag.String("ip", "", "Target IP, IP range, or CIDR (e.g., 192.168.1.1, 192.168.1.1-10, 192.168.1.0/24)")
	domain := flag.String("domain", "", "Single domain to test")
	domainFile := flag.String("domains", "", "Path to file containing list of domains")
	wordlistFile := flag.String("wordlist", "", "Path to wordlist file for subdomain enumeration")
	threads := flag.Int("threads", 5, "Number of concurrent threads")
	requestTimeout := flag.Int("timeout", 10, "HTTP request timeout in seconds")
	match := flag.String("match", "200", "Comma-separated list of status codes to include")
	ignoreCert := flag.Bool("ignoreCert", true, "Ignore SSL certificate verification errors")
	format := flag.String("f", "text", "Output format (json, csv, or text)")
	output := flag.String("o", "", "Output file path")
	verbose := flag.Bool("v", false, "Show all requests and checks")
	flag.Parse()

	if *ip == "" {
		fmt.Println("ERROR: IP address (-ip) is required.")
		fmt.Println("")
		fmt.Println("Use -h or --help for usage information.")
		return
	}

	if *domain == "" && *domainFile == "" {
		fmt.Println("ERROR: Either -domain or -domains is required.")
		fmt.Println("")
		fmt.Println("Use -h or --help for usage information.")
		return
	}

	if *domain != "" && *domainFile != "" {
		fmt.Println("ERROR: Cannot use both -domain and -domains together. Choose one.")
		fmt.Println("")
		fmt.Println("Use -h or --help for usage information.")
		return
	}

	matchCodes := parseStatusCodes(*match)

	// Parse IP addresses
	ips, err := parseIPInput(*ip)
	if err != nil {
		fmt.Printf("ERROR: Invalid IP format: %v\n", err)
		fmt.Println("")
		fmt.Println("Supported formats:")
		fmt.Println("  - Single IP: 192.168.1.1")
		fmt.Println("  - IP Range: 192.168.1.1-192.168.1.10")
		fmt.Println("  - CIDR: 192.168.1.0/24")
		return
	}

	// Get domains list
	var domains []string
	if *domain != "" {
		domains = []string{*domain}
	} else {
		domains, err = readDomainsFromFile(*domainFile)
		if err != nil {
			fmt.Printf("ERROR: Reading domains from file: %v\n", err)
			return
		}
	}

	// Read wordlist if provided and combine with domains
	if *wordlistFile != "" {
		wordlist, err := readWordlistFromFile(*wordlistFile)
		if err != nil {
			fmt.Printf("ERROR: Reading wordlist from file: %v\n", err)
			return
		}
		domains = combineWordlistWithDomains(domains, wordlist)
	}

	if *verbose {
		fmt.Printf("Scanning %d IP(s) with %d domain(s)\n", len(ips), len(domains))
		fmt.Printf("Total requests: %d (http+https per domain per IP)\n", len(ips)*len(domains)*2)
	}

	client := &http.Client{Timeout: time.Duration(*requestTimeout) * time.Second}
	if *ignoreCert {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, *threads)
	var results []RequestResult
	var resultsMutex sync.Mutex

	// Iterate over all IPs and domains
	for _, targetIP := range ips {
		for _, targetDomain := range domains {
			for _, protocol := range []string{"http", "https"} {
				wg.Add(1)
				go func(ipAddr, domain, proto string) {
					defer wg.Done()
					semaphore <- struct{}{}

					result := RequestResult{
						Domain:   domain,
						IP:       ipAddr,
						Protocol: proto,
					}

					requestURL := fmt.Sprintf("%s://%s", proto, ipAddr)
					req, err := http.NewRequest("GET", requestURL, nil)
					if err != nil {
						result.Error = fmt.Sprintf("Failed to create request: %v", err)
					} else {
						req.Host = domain
						if *verbose {
							fmt.Printf("Checking %s://%s (IP: %s)\n", proto, domain, ipAddr)
						}

						resp, err := client.Do(req)
						if err != nil {
							// Skip failed requests entirely
							<-semaphore
							return
						}
						defer resp.Body.Close()

						result.StatusCode = resp.StatusCode
						if !statusCodeMatches(result.StatusCode, matchCodes) {
							<-semaphore
							return
						}

						// Only append successful results
						resultsMutex.Lock()
						results = append(results, RequestResult{
							Domain:     domain,
							IP:         ipAddr,
							Protocol:   proto,
							StatusCode: resp.StatusCode,
						})
						resultsMutex.Unlock()
					}

					<-semaphore
				}(targetIP, targetDomain, protocol)
			}
		}
	}

	wg.Wait()

	// Output handling based on format
	var outputString string
	switch *format {
	case "json":
		jsonData, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Printf("Error encoding JSON: %v\n", err)
			return
		}
		outputString = string(jsonData)
	case "csv":
		var csvLines []string
		csvLines = append(csvLines, "domain,ip,protocol,status_code,error")
		for _, r := range results {
			csvLines = append(csvLines, fmt.Sprintf("%s,%s,%s,%d,%s",
				r.Domain, r.IP, r.Protocol, r.StatusCode, r.Error))
		}
		outputString = strings.Join(csvLines, "\n")
	default:
		var lines []string
		for _, r := range results {
			if r.Error != "" {
				lines = append(lines, fmt.Sprintf("%s://%s (IP: %s) - Error: %s",
					r.Protocol, r.Domain, r.IP, r.Error))
			} else {
				lines = append(lines, fmt.Sprintf("%s://%s (IP: %s) - Status: %d",
					r.Protocol, r.Domain, r.IP, r.StatusCode))
			}
		}
		outputString = strings.Join(lines, "\n")
	}

	if *output != "" {
		// Create directory if it doesn't exist
		dir := filepath.Dir(*output)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("Error creating output directory: %v\n", err)
				return
			}
		}

		err := os.WriteFile(*output, []byte(outputString), 0644)
		if err != nil {
			fmt.Printf("Error writing to output file: %v\n", err)
			return
		}
	} else {
		fmt.Println(outputString)
	}

}

func statusCodeMatches(code int, match []int) bool {
	return inSlice(code, match)
}

func inSlice(val int, slice []int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func parseStatusCodes(s string) []int {
	var codes []int
	if s == "" {
		return codes
	}
	for _, part := range strings.Split(s, ",") {
		code, err := strconv.Atoi(part)
		if err == nil {
			codes = append(codes, code)
		}
	}
	return codes
}

func readDomainsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain != "" {
			domains = append(domains, domain)
		}
	}
	return domains, scanner.Err()
}

func readWordlistFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var wordlist []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			wordlist = append(wordlist, word)
		}
	}
	return wordlist, scanner.Err()
}

func combineWordlistWithDomains(domains []string, wordlist []string) []string {
	var combinedDomains []string
	for _, domain := range domains {
		for _, word := range wordlist {
			combinedDomains = append(combinedDomains, fmt.Sprintf("%s.%s", word, domain))
		}
	}
	return combinedDomains
}

// parseIPInput parses IP input supporting single IP, IP range, or CIDR notation
func parseIPInput(input string) ([]string, error) {
	input = strings.TrimSpace(input)

	// Check if it's CIDR notation
	if strings.Contains(input, "/") {
		return parseCIDR(input)
	}

	// Check if it's IP range
	if strings.Contains(input, "-") {
		return parseIPRange(input)
	}

	// Single IP address
	ip := net.ParseIP(input)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", input)
	}
	return []string{input}, nil
}

// parseCIDR parses CIDR notation and returns all IPs in the range
func parseCIDR(cidr string) ([]string, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR notation: %v", err)
	}

	var ips []string
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
		ips = append(ips, ip.String())
	}

	// Remove network and broadcast addresses for /24 and smaller
	if len(ips) > 2 {
		return ips[1 : len(ips)-1], nil
	}
	return ips, nil
}

// parseIPRange parses IP range like 192.168.1.1-192.168.1.10 or 192.168.1.1-10
func parseIPRange(ipRange string) ([]string, error) {
	parts := strings.Split(ipRange, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid IP range format")
	}

	startIP := net.ParseIP(strings.TrimSpace(parts[0]))
	if startIP == nil {
		return nil, fmt.Errorf("invalid start IP in range")
	}

	endPart := strings.TrimSpace(parts[1])
	var endIP net.IP

	// Check if end is just a number (e.g., 192.168.1.1-10)
	if !strings.Contains(endPart, ".") {
		// Extract the base IP and replace last octet
		ipParts := strings.Split(parts[0], ".")
		if len(ipParts) != 4 {
			return nil, fmt.Errorf("invalid IP format")
		}
		endIP = net.ParseIP(fmt.Sprintf("%s.%s.%s.%s", ipParts[0], ipParts[1], ipParts[2], endPart))
	} else {
		endIP = net.ParseIP(endPart)
	}

	if endIP == nil {
		return nil, fmt.Errorf("invalid end IP in range")
	}

	var ips []string
	for ip := startIP; !ip.Equal(endIP); incrementIP(ip) {
		ips = append(ips, ip.String())
		// Safety check to prevent infinite loops
		if len(ips) > 65536 {
			return nil, fmt.Errorf("IP range too large (max 65536 IPs)")
		}
	}
	ips = append(ips, endIP.String())

	return ips, nil
}

// incrementIP increments an IP address by 1
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
