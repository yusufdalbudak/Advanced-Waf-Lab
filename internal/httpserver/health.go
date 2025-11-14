package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/waf-draft/waf/internal/telemetry"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// ServeHTTP implements http.Handler for health checks
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(telemetry.GetMetrics().StartTime).Seconds(),
	}
	
	json.NewEncoder(w).Encode(response)
}

// MetricsHandler handles metrics requests
type MetricsHandler struct{}

// ServeHTTP implements http.Handler for metrics
func (m *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	stats := telemetry.GetMetrics().GetStats()
	json.NewEncoder(w).Encode(stats)
}

