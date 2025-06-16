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

func TestIsSubdomain(t *testing.T) {
	input := []struct {
		Domain string
		IsSub bool
	}{
		{"https://test.byeoungwoolee.com", true},
		{"https://byeoungwoolee.com", false},
		{"byeoungwoolee.com", false},
		{"test.byeoungwoolee.com", true},
	}

	for _, data := range input {
		isSub, err := IsSubdomainWithoutScheme(data.Domain)
		if err != nil {
			t.Error(err)
		}
		t.Logf("isSub : %v, data.IsSub : %v", isSub, data.IsSub)
		if isSub != data.IsSub {
			t.Errorf("check failed %v. Expected %v, got %v",data.Domain, data.IsSub, isSub)
		}

		// isSub, err = IsSubdomainWithScheme(data.Domain)
		// if err != nil {
		// 	t.Error(err)
		// }
		// if isSub != data.IsSub {
		// 	t.Errorf("check failed %v. Expected %v, got %v",data.Domain, data.IsSub, isSub)
		// }
	}
	
}


func TestDomainValidator(t *testing.T) {
	input := []struct {
		Domain string
		IsSub bool
	}{
		{"https://test.byeoungwoolee.com", false},
		{"https://byeoungwoolee.com", false},
		{"byeoungwoolee.com", false},
		{"test.byeoungwoolee.com", true},
	}

	for _, data := range input {
		valid := DomainValidator(data.Domain, data.IsSub)
		if valid != data.IsSub {
			t.Errorf("check failed %v. Expected %v, got %v",data.Domain, data.IsSub, valid)
		}
	}
	
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
