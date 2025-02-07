package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
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

func main() {
	ip := flag.String("ip", "", "The IP address to check")
	domainFile := flag.String("domains", "", "Path to the file containing domains")
	threads := flag.Int("threads", 5, "Number of concurrent threads")
	requestTimeout := flag.Int("timeout", 10, "HTTP request timeout in seconds")
	match := flag.String("match", "200", "Comma-separated list of status codes to include")
	ignoreCert := flag.Bool("ignoreCert", true, "Ignore SSL certificate verification errors")
	format := flag.String("f", "text", "Output format (json, csv, or text)")
	output := flag.String("o", "", "Output file path")
	verbose := flag.Bool("v", false, "Show all requests and checks")
	flag.Parse()

	if *ip == "" || *domainFile == "" {
		fmt.Println("IP address and domain file path are required.")
		return
	}

	matchCodes := parseStatusCodes(*match)

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

	domains, err := readDomainsFromFile(*domainFile)
	if err != nil {
		fmt.Printf("Error reading domains from file: %v\n", err)
		return
	}

	for _, domain := range domains {
		for _, protocol := range []string{"http", "https"} {
			wg.Add(1)
			go func(domain, protocol string) {
				defer wg.Done()
				semaphore <- struct{}{}

				result := RequestResult{
					Domain:   domain,
					IP:       *ip,
					Protocol: protocol,
				}

				requestURL := fmt.Sprintf("%s://%s", protocol, *ip)
				req, err := http.NewRequest("GET", requestURL, nil)
				if err != nil {
					result.Error = fmt.Sprintf("Failed to create request: %v", err)
				} else {
					req.Host = domain
					if *verbose {
						fmt.Printf("Checking %s://%s (IP: %s)\n", protocol, domain, *ip)
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
						IP:         *ip,
						Protocol:   protocol,
						StatusCode: resp.StatusCode,
					})
					resultsMutex.Unlock()
				}

				<-semaphore
			}(domain, protocol)
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
