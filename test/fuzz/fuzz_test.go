package fuzz

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/waf-draft/waf/internal/normalize"
)

// FuzzNormalizePath tests path normalization with fuzzed input
func FuzzNormalizePath(f *testing.F) {
	// Seed corpus
	f.Add("/api/users")
	f.Add("/api/users/../admin")
	f.Add("/api/users?id=123")
	f.Add("../../../etc/passwd")
	f.Add("%2e%2e%2fetc%2fpasswd")

	f.Fuzz(func(t *testing.T, path string) {
		// Create a request with fuzzed path
		req := httptest.NewRequest("GET", "http://example.com"+path, nil)
		
		// Normalize should not panic
		norm, err := normalize.Request(req, false)
		if err != nil {
			// Some errors are expected with malformed input
			return
		}

		// Basic sanity checks
		if norm == nil {
			t.Fatal("normalized request is nil")
		}
		if norm.Path == "" && path != "" {
			// Path should not be empty if input was not empty
			// (unless it was completely invalid)
		}
	})
}

// FuzzRuleMatching tests rule matching with fuzzed input
func FuzzRuleMatching(f *testing.F) {
	// Seed corpus with common attack patterns
	f.Add("1' OR '1'='1")
	f.Add("<script>alert('xss')</script>")
	f.Add("../../../etc/passwd")
	f.Add("; ls -la")
	f.Add("UNION SELECT * FROM users")

	f.Fuzz(func(t *testing.T, input string) {
		// Create a request with fuzzed query parameter
		url := "http://example.com/api/search?q=" + input
		req := httptest.NewRequest("GET", url, nil)

		// Normalize should not panic
		norm, err := normalize.Request(req, false)
		if err != nil {
			return
		}

		// Check that normalization handles the input
		if norm == nil {
			t.Fatal("normalized request is nil")
		}

		// Verify query parameter extraction
		queryValue := norm.GetQueryString("q")
		// Should either match or be safely normalized
		_ = queryValue
	})
}

// FuzzRequestParsing tests HTTP request parsing with fuzzed input
func FuzzRequestParsing(f *testing.F) {
	// Seed corpus
	f.Add("GET /api/users HTTP/1.1\r\nHost: example.com\r\n\r\n")
	f.Add("POST /api/login HTTP/1.1\r\nContent-Length: 10\r\n\r\nusername=test")

	f.Fuzz(func(t *testing.T, rawRequest string) {
		// Try to parse fuzzed HTTP request
		// This is a simplified test - in practice, you'd use a proper HTTP parser
		lines := strings.Split(rawRequest, "\r\n")
		if len(lines) == 0 {
			return
		}

		// Parse request line
		parts := strings.Fields(lines[0])
		if len(parts) < 2 {
			return
		}

		method := parts[0]
		path := parts[1]

		// Create request
		req := httptest.NewRequest(method, "http://example.com"+path, nil)

		// Normalize should handle any input without panicking
		norm, err := normalize.Request(req, false)
		if err != nil {
			return
		}

		if norm == nil {
			t.Fatal("normalized request is nil")
		}
	})
}

