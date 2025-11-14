package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// PrometheusMetrics holds Prometheus metric collectors
type PrometheusMetrics struct {
	RequestsTotal      *prometheus.CounterVec
	RequestsBlocked    prometheus.Counter
	RequestsAllowed    prometheus.Counter
	RequestDuration    *prometheus.HistogramVec
	AnomalyScore       *prometheus.HistogramVec
	RuleMatches        *prometheus.CounterVec
	ActiveConnections  prometheus.Gauge
}

var promMetrics *PrometheusMetrics

// InitPrometheus initializes Prometheus metrics
func InitPrometheus() *PrometheusMetrics {
	if promMetrics != nil {
		return promMetrics
	}

	promMetrics = &PrometheusMetrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "waf_requests_total",
				Help: "Total number of requests processed",
			},
			[]string{"method", "path", "status"},
		),
		RequestsBlocked: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "waf_requests_blocked_total",
				Help: "Total number of blocked requests",
			},
		),
		RequestsAllowed: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "waf_requests_allowed_total",
				Help: "Total number of allowed requests",
			},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "waf_request_duration_seconds",
				Help:    "Request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		AnomalyScore: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "waf_anomaly_score",
				Help:    "Anomaly score distribution",
				Buckets: []float64{0, 5, 10, 15, 20, 25, 30, 50, 100},
			},
			[]string{"action"},
		),
		RuleMatches: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "waf_rule_matches_total",
				Help: "Total number of rule matches",
			},
			[]string{"rule_id", "rule_tag"},
		),
		ActiveConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "waf_active_connections",
				Help: "Number of active connections",
			},
		),
	}

	return promMetrics
}

// GetPrometheusMetrics returns the Prometheus metrics instance
func GetPrometheusMetrics() *PrometheusMetrics {
	if promMetrics == nil {
		return InitPrometheus()
	}
	return promMetrics
}

