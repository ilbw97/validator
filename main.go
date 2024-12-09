package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type TestData struct {
	Host           string `json:"host"`
	RequiredScheme bool   `json:"required_scheme"`
	Path           string `json:"path"`
}

var mux sync.Mutex

func IsSubdomain(input string) (bool, error) {
	// Parse the input as a URL
	parsedURL, err := url.Parse(input)
	if err != nil {
		return false, err
	}

	// Extract the host from the parsed URL
	host := parsedURL.Host

	// Split the host into parts
	parts := strings.Split(host, ".")

	// A root domain typically has two parts (e.g., example.com)
	// A subdomain has more than two parts (e.g., sub.example.com)
	return len(parts) > 2, nil
}

func URIValidator(s string, requiredScheme bool) bool {
	if s == "" {
		return false
	}

	parsedURI, err := url.Parse(s)
	if err != nil {
		fmt.Printf("Error parsing URL: %v\n", err)
		return false // Invalid URI format
	}

	fmt.Printf("%s's parsedURI : %s\n", s, parsedURI)

	// Check for valid scheme
	if requiredScheme {
		switch parsedURI.Scheme {
		case "http", "https":
		default:
			fmt.Printf("%v is invalid scheme \n", parsedURI.Scheme)
			return false
		}
	}

	// fmt.Printf("parsedURI.Host : %s\n", parsedURI.Host)

	// Validate host
	if parsedURI.Host != "" {
		host := parsedURI.Hostname()
		fmt.Printf("parsedURI.Host : %s\n", parsedURI.Host)
		fmt.Printf("parsedURI.Hostname : %s\n", host)

		// Check if host contains a dot and does not start or end with a dot
		if !strings.Contains(host, ".") || strings.HasPrefix(host, ".") || strings.HasSuffix(host, ".") {
			fmt.Printf("%s has problem with dot\n", host)
			return false
		}

		// Check if host is a valid domain name or IP address
		if net.ParseIP(host) == nil {
			if len(host) > 253 || len(strings.Split(host, ".")) < 2 {
				fmt.Printf("%s is invalid domain name, or too long. host length : %v\n", host, len(host))
				return false
			}
			for _, label := range strings.Split(host, ".") {
				if len(label) > 63 || !isAlphaNumeric(label) {
					fmt.Printf("%s is invalid domain name, or too long. label length : %v\n", label, len(label))
					return false
				}
			}
		} else {
			return false
		}

		// Validate port number if present
		if parsedURI.Port() != "" {
			port, err := strconv.Atoi(parsedURI.Port())
			if err != nil || port < 1 || port > 65535 {
				return false
			}
		}
	} else {
		return false
	}

	// Validate path (basic check for illegal characters)
	if strings.Contains(parsedURI.Path, " ") {
		return false
	}

	return true
}

// Helper function to check if a string is alphanumeric (including hyphen)
func isAlphaNumeric(s string) bool {
	for _, c := range s {
		if !(('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') || c == '-') {
			return false
		}
	}
	return true
}

func PathValidator(s string) bool {
	if len(s) == 0 {
		return false
	}
	u, err := url.Parse(s)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("parsedUrl : %v, path : %v, scheme : %v, host : %v, rawQuery : %v\n", u, u.Path, u.Scheme, u.Host, u.RawQuery)
	return err == nil && u.Path != "" && u.Scheme == "" && u.Host == "" && u.RawQuery == ""
}

func DomainValidator(domain string) bool {

	// The domain length should be greater than 1 and less than 253 characters
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	// Combined regex pattern to check domain validity
	domainPattern := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z]{2,})+$`)
	if !domainPattern.MatchString(domain) {
		fmt.Printf("Invalid domain format: %s\n", domain)
		return false
	} else {
		fmt.Println("Valid domain format")
	}

	return true
}

func main() {

	mux.Lock()

	fname := "./test.json"
	jsonData, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println(err)
		return
	}

	testDatas := []TestData{}
	// JSON 파일을 읽어오기
	err = json.Unmarshal(jsonData, &testDatas)
	if err != nil {
		fmt.Println(err)
		return
	}

	mux.Unlock()

	for _, testData := range testDatas {
		fmt.Println("-------------------------------------------------------------------------")
		if testData.Host != "" {
			if domainRes := DomainValidator(testData.Host); domainRes {
				fmt.Printf("DomainValidator said %s is valid domain \n\n", testData.Host)
			} else {
				fmt.Printf("DomainValidator said %s is invalid domain \n\n", testData.Host)
			}

			res := URIValidator(testData.Host, testData.RequiredScheme)
			if res {
				fmt.Printf("%v is valid url \n", testData.Host)
				isSub, err := IsSubdomain(testData.Host)
				if err != nil {
					fmt.Println(err)
				}
				if isSub {
					fmt.Printf("%v is sub domain \n", testData.Host)
				} else {
					fmt.Printf("%v is root domain \n", testData.Host)
				}
			} else {
				fmt.Printf("%v is invalid url \n", testData.Host)
			}
		}
		if testData.Path != "" {
			if isValid := PathValidator(testData.Path); isValid {
				fmt.Printf("%v is valid path \n", testData.Host)
			} else {
				fmt.Printf("%v is invalid path \n", testData.Host)
			}
		}
		fmt.Println("-------------------------------------------------------------------------")
	}
}
