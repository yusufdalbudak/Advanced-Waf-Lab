package httpserver

import (
	"net/http"
	"strings"
)

// Router handles routing for WAF endpoints
type Router struct {
	wafHandler     *WAFHandler
	healthHandler  *HealthHandler
	metricsHandler *MetricsHandler
	logsHandler    *LogsHandler
}

// NewRouter creates a new router
func NewRouter(wafHandler *WAFHandler, logFile string) *Router {
	return &Router{
		wafHandler:     wafHandler,
		healthHandler:  &HealthHandler{},
		metricsHandler: &MetricsHandler{},
		logsHandler:    NewLogsHandler(logFile),
	}
}

// ServeHTTP implements http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	// Handle management endpoints (bypass WAF)
	if strings.HasPrefix(path, "/health") {
		r.healthHandler.ServeHTTP(w, req)
		return
	}

	if strings.HasPrefix(path, "/metrics") {
		r.metricsHandler.ServeHTTP(w, req)
		return
	}

	if strings.HasPrefix(path, "/logs") {
		r.logsHandler.ServeHTTP(w, req)
		return
	}

	// All other requests go through WAF
	r.wafHandler.ServeHTTP(w, req)
}

