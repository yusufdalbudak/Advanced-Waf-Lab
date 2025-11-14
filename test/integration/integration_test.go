package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/waf-draft/waf/internal/config"
	"github.com/waf-draft/waf/internal/detection/rules"
	"github.com/waf-draft/waf/internal/httpserver"
	"github.com/waf-draft/waf/internal/logging"
)

// createTestUpstreamServer creates a mock upstream server for testing
func createTestUpstreamServer(t *testing.T) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
			"path":   r.URL.Path,
		})
	})
	return httptest.NewServer(handler)
}

// createTestWAFServer creates a WAF server for testing
func createTestWAFServer(t *testing.T, upstreamURL string) (*httptest.Server, error) {
	// Create minimal config
	cfg := &config.Config{
		Server: config.ServerConfig{
			ListenAddress:      ":0", // Let system choose port
			UpstreamURL:        upstreamURL,
			ReadTimeoutSeconds: 10,
			WriteTimeoutSeconds: 10,
			IdleTimeoutSeconds: 60,
		},
		Security: config.SecurityConfig{
			AnomalyThreshold: 10,
			LogRequestBody:   false,
		},
		Logging: config.LoggingConfig{
			Level:  "info",
			Output: "stdout",
		},
		Rules: config.RulesConfig{
			Files: []string{"../../configs/ruleset.yaml"},
		},
	}

	// Load rules
	ruleSet, err := rules.LoadRules(cfg.Rules.Files)
	if err != nil {
		return nil, fmt.Errorf("failed to load rules: %w", err)
	}

	// Initialize logger
	logger, err := logging.NewLogger(cfg.Logging.Output)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	// Create reverse proxy
	proxy, err := httpserver.NewProxy(cfg.Server.UpstreamURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create reverse proxy: %w", err)
	}

	// Create WAF handler
	handler := httpserver.NewWAFHandler(cfg, ruleSet, logger, proxy)

	// Create test server
	server := httptest.NewServer(handler)
	return server, nil
}

func TestBenignRequest(t *testing.T) {
	// Start upstream server
	upstream := createTestUpstreamServer(t)
	defer upstream.Close()

	// Start WAF server
	wafServer, err := createTestWAFServer(t, upstream.URL)
	if err != nil {
		t.Fatalf("Failed to create WAF server: %v", err)
	}
	defer wafServer.Close()

	// Make benign request
	resp, err := http.Get(wafServer.URL + "/api/users?id=123")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Should be allowed (200 OK)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify response from upstream
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %v", result["status"])
	}
}

func TestSQLInjectionBlock(t *testing.T) {
	// Start upstream server
	upstream := createTestUpstreamServer(t)
	defer upstream.Close()

	// Start WAF server
	wafServer, err := createTestWAFServer(t, upstream.URL)
	if err != nil {
		t.Fatalf("Failed to create WAF server: %v", err)
	}
	defer wafServer.Close()

	// Make SQL injection request (URL encode the query parameter)
	sqliURL := wafServer.URL + "/api/users?id=1%20OR%201=1"
	resp, err := http.Get(sqliURL)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Should be blocked (403 Forbidden)
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", resp.StatusCode)
	}

	// Verify error response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var errorResp map[string]interface{}
	if err := json.Unmarshal(body, &errorResp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errorResp["error"] != "Forbidden" {
		t.Errorf("Expected error 'Forbidden', got %v", errorResp["error"])
	}
}

func TestXSSBlock(t *testing.T) {
	// Start upstream server
	upstream := createTestUpstreamServer(t)
	defer upstream.Close()

	// Start WAF server
	wafServer, err := createTestWAFServer(t, upstream.URL)
	if err != nil {
		t.Fatalf("Failed to create WAF server: %v", err)
	}
	defer wafServer.Close()

	// Make XSS request
	xssURL := wafServer.URL + "/api/search?q=<script>alert('xss')</script>"
	resp, err := http.Get(xssURL)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Should be blocked (403 Forbidden)
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", resp.StatusCode)
	}
}

func TestPathTraversalBlock(t *testing.T) {
	// Start upstream server
	upstream := createTestUpstreamServer(t)
	defer upstream.Close()

	// Start WAF server
	wafServer, err := createTestWAFServer(t, upstream.URL)
	if err != nil {
		t.Fatalf("Failed to create WAF server: %v", err)
	}
	defer wafServer.Close()

	// Make path traversal request
	ptURL := wafServer.URL + "/api/files/../../../etc/passwd"
	resp, err := http.Get(ptURL)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Should be blocked (403 Forbidden)
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", resp.StatusCode)
	}
}

func TestMultipleRules(t *testing.T) {
	// Start upstream server
	upstream := createTestUpstreamServer(t)
	defer upstream.Close()

	// Start WAF server
	wafServer, err := createTestWAFServer(t, upstream.URL)
	if err != nil {
		t.Fatalf("Failed to create WAF server: %v", err)
	}
	defer wafServer.Close()

	// Make request that triggers multiple rules (URL encode)
	maliciousURL := wafServer.URL + "/api/search?q=1%27%20UNION%20SELECT%20*%20FROM%20users--"
	resp, err := http.Get(maliciousURL)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Should be blocked (403 Forbidden)
	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", resp.StatusCode)
	}
}

func TestPOSTRequest(t *testing.T) {
	// Start upstream server
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer upstream.Close()

	// Start WAF server
	wafServer, err := createTestWAFServer(t, upstream.URL)
	if err != nil {
		t.Fatalf("Failed to create WAF server: %v", err)
	}
	defer wafServer.Close()

	// Make benign POST request
	body := bytes.NewBufferString(`{"name": "test"}`)
	resp, err := http.Post(wafServer.URL+"/api/users", "application/json", body)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Should be allowed (200 OK)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

