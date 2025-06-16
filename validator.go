package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type trafficTypeReqBlock struct {
	Type  string `json:"type" validate:"oneof=traffic bps visit threat cache"`
	Stime int64  `json:"stime" validate:"required,min=1"`
	Etime int64  `json:"etime" validate:"required,gtfield=stime"`
}

// func UseValidator(req interface{}) bool {
func UseValidator() bool {
	jsonData := []byte(`{
        "stime": 1721192000,
        "etime": 1733420270,
        "type": "traffic"
    }`)

	var req trafficTypeReqBlock
	err := json.Unmarshal(jsonData, &req)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return false
	}

	fmt.Printf("Parsed Struct: %+v\n", req)

	v := validator.New()
	err = v.Struct(req)
	if err != nil {
		fmt.Println("Validation Error:", err)
		return false
	}

	return true
}

func CheckBase64(input string) bool {
	_, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		fmt.Println("Invalid Base64:", err)
		return false
	}
	// // Check if the input string matches the base64 pattern
	// return regexp.MustCompile(`^[a-zA-Z0-9+/]*={0,2}$`).MatchString(input)
	return true
}

// without scheme
func IsSubdomainWithoutScheme(input string) (bool, error) {
	// Split the input into parts by "."
	parts := strings.Split(input, ".")
	fmt.Printf("parts : %v,len : %v\n", parts, len(parts))

	// A root domain typically has two parts (e.g., example.com)
	// A subdomain has more than two parts (e.g., sub.example.com)
	return len(parts) > 2, nil
}

// with scheme
func IsSubdomainWithScheme(input string) (bool, error) {
	// Ensure input is a valid URL by adding a scheme if missing
	if !strings.HasPrefix(input, "http://") && !strings.HasPrefix(input, "https://") {
		input = "http://" + input
	}

	// Parse the input as a URL
	parsedURL, err := url.Parse(input)
	if err != nil {
		return false, err
	}

	// Extract the host from the parsed URL
	host := parsedURL.Host

	// Handle cases where port is included (e.g., "example.com:8080")
	host = strings.Split(host, ":")[0]

	// Split the host into parts
	parts := strings.Split(host, ".")

	// A root domain typically has two parts (e.g., example.com)
	// A subdomain has more than two parts (e.g., sub.example.com)
	return len(parts) > 2, nil
}


// func IsSubdomain(input string) (bool, error) {
// 	// Parse the input as a URL
// 	parsedURL, err := url.Parse(input)
// 	if err != nil {
// 		return false, err
// 	}

// 	parsedURL.String()
// 	// Extract the host from the parsed URL
// 	host := parsedURL.Host

// 	// Split the host into parts
// 	parts := strings.Split(host, ".")

// 	// A root domain typically has two parts (e.g., example.com)
// 	// A subdomain has more than two parts (e.g., sub.example.com)
// 	return len(parts) > 2, nil
// }

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

// func DomainValidator(domain string) bool {

// 	// The domain length should be greater than 1 and less than 253 characters
// 	if len(domain) == 0 || len(domain) > 253 {
// 		return false
// 	}

// 	// Combined regex pattern to check domain validity
// 	domainPattern := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z]{2,})+$`)
// 	schemePattern := regexp.MustCompile(`^https?://`)
// 	if !domainPattern.MatchString(domain) || !schemePattern.MatchString(domain) {
// 		fmt.Printf("Invalid domain format: %s\n", domain)
// 		return false
// 	} else {
// 		fmt.Println("Valid domain format")
// 	}

// 	return true
// }

func DomainValidator(domain string, isSubDomain bool) bool {

	fmt.Printf("domain : %s, isSubDomain : %v\n", domain, isSubDomain)
	// The domain length should be greater than 1 and less than 253 characters
	if len(domain) == 0 || len(domain) > 253 {
		fmt.Printf("Invalid domain length: %d\n", len(domain))
		return false
	}

	// Combined regex pattern to check domain validity
	domainPattern := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z]{2,})+$`)
	if !domainPattern.MatchString(domain) {
		fmt.Printf("Invalid domain format: %s\n", domain)
		return false
	}

	var domainData string
	if strings.HasPrefix(domain, "http://") || strings.HasPrefix(domain, "https://") {
		parsedUrl, err := url.Parse(domain)
		if err != nil{
			fmt.Printf("Error parsing URL: %v\n", err)
			return false
		}
		domainData = parsedUrl.Host
	} else {
		domainData = domain
	}

	fmt.Printf("\n\ndomainData : %s\n", domainData)

	isSub, err := IsSubdomainWithScheme(domainData)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if !isSub {
		fmt.Printf("%s is not sub domain\n", domain)
		return false
	}else{
		fmt.Printf("%s is sub domain\n", domain)
		return true
	}

}