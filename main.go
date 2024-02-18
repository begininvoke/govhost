package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Define a struct to hold the request result, including a field for errors
type RequestResult struct {
	Address    string `json:"address,omitempty"`
	IP         string `json:"ip"`
	StatusCode int    `json:"status_code"`
	Error      string `json:"error,omitempty"` // Include error message if any, omit from JSON if empty
}

func main() {
	hostname := flag.String("hostname", "", "The hostname to check")
	ipFilePath := flag.String("ipfile", "", "Path to the file containing IP addresses")
	threads := flag.Int("threads", 5, "Number of concurrent threads")
	outputJSON := flag.Bool("oj", false, "Output results in JSON format")
	requestTimeout := flag.Int("timeout", 10, "HTTP request timeout in seconds")
	protocol := flag.String("protocol", "http", "Protocol to use for requests (http or https)")
	match := flag.String("match", "", "Comma-separated list of status codes to include")
	notMatch := flag.String("notmatch", "", "Comma-separated list of status codes to exclude")
	ignoreCert := flag.Bool("ignoreCert", false, "Ignore SSL certificate verification errors")
	flag.Parse()

	matchCodes := parseStatusCodes(*match)
	notMatchCodes := parseStatusCodes(*notMatch)

	if *hostname == "" || *ipFilePath == "" || (*protocol != "http" && *protocol != "https") {
		fmt.Println("Hostname, IP file path, and valid protocol (http or https) are required.")
		return
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

	ipList, err := readIPsFromFile(*ipFilePath)
	if err != nil {
		fmt.Printf("Error reading IP addresses from file: %v\n", err)
		return
	}

	for _, ip := range ipList {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			semaphore <- struct{}{}

			result := RequestResult{IP: ip}
			requestURL := fmt.Sprintf("%s://%s", *protocol, ip)
			req, err := http.NewRequest("GET", requestURL, nil)
			if err != nil {
				result.Error = fmt.Sprintf("Failed to create request: %v", err)
			} else {
				req.Host = *hostname
				resp, err := client.Do(req)
				if err != nil {
					result.StatusCode = 0
					result.Error = fmt.Sprintf("Request failed: %v", err)
				} else {
					defer resp.Body.Close()
					result.StatusCode = resp.StatusCode
					if !statusCodeMatches(result.StatusCode, matchCodes, notMatchCodes) {
						<-semaphore
						return
					}
				}
			}

			resultsMutex.Lock()
			results = append(results, result)
			resultsMutex.Unlock()

			<-semaphore
		}(ip)
	}

	wg.Wait()

	if *outputJSON {
		jsonResults, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Printf("Failed to encode results to JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonResults))
	} else {
		for _, result := range results {
			if result.Error != "" {
				fmt.Printf("IP %s encountered an error: %s\n", result.IP, result.Error)
			} else {
				fmt.Printf("IP %s responded with status code: %d\n", result.IP, result.StatusCode)
			}
		}
	}
}

func statusCodeMatches(code int, match, notMatch []int) bool {
	if len(match) > 0 && !inSlice(code, match) {
		return false
	}
	if inSlice(code, notMatch) {
		return false
	}
	return true
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

func readIPsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var ips []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := scanner.Text()
		if ip != "" {
			ips = append(ips, ip)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ips, nil
}
