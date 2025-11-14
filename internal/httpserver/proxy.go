package httpserver

import (
	"net/http"
	"net/http/httputil"

	"github.com/waf-draft/waf/internal/mitigation"
)

// NewProxy creates a reverse proxy handler
func NewProxy(upstreamURL string) (http.Handler, error) {
	proxy, err := mitigation.NewReverseProxy(upstreamURL)
	if err != nil {
		return nil, err
	}

	return proxy, nil
}

// ProxyHandler wraps the reverse proxy
type ProxyHandler struct {
	proxy *httputil.ReverseProxy
}

// ServeHTTP implements http.Handler for the proxy
func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

