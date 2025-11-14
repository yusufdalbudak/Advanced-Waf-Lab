package mitigation

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/waf-draft/waf/internal/decision"
)

// ApplyDecision applies the WAF decision to the request
func ApplyDecision(dec decision.Decision, w http.ResponseWriter, r *http.Request, proxy *httputil.ReverseProxy) {
	if dec.Action == "block" {
		blockRequest(w, dec)
	} else {
		// Forward request to upstream
		proxy.ServeHTTP(w, r)
	}
}

// blockRequest returns a 403 Forbidden response
func blockRequest(w http.ResponseWriter, dec decision.Decision) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)

	response := map[string]interface{}{
		"error":   "Forbidden",
		"message": "Request blocked by WAF",
		"reason":  dec.Reason,
	}

	json.NewEncoder(w).Encode(response)
}

// NewReverseProxy creates a new reverse proxy for the upstream server
func NewReverseProxy(upstreamURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(upstreamURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Customize the proxy director
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = target.Host
	}

	return proxy, nil
}

