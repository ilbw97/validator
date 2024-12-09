package main

import (
	"net/url"
	"regexp"
	"testing"
)

func URIValidatortest(s string) bool {
	if len(s) == 0 {
		return false
	}

	parsedURI, err := url.Parse(s)
	if err != nil {
		return false
	}

	if parsedURI.Host == "" && parsedURI.Path == "" {
		return false
	}

	return true
}

func ValidURL(target string) bool {
	reURL, err := regexp.Compile(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)[a-z0-9]+([\-\.]{1}[a-z0-9-]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`)
	if err != nil {
		return false
	}

	return reURL.MatchString(target)
}

// Benchmark tests for each function
func BenchmarkURIValidator(b *testing.B) {
	testURL := "http://example.com/path/to/resource"
	for i := 0; i < b.N; i++ {
		URIValidatortest(testURL)
	}
}

func BenchmarkValidURL(b *testing.B) {
	testURL := "http://example.com/path/to/resource"
	for i := 0; i < b.N; i++ {
		ValidURL(testURL)
	}
}
