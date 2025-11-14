package httpserver

import (
	"net/http"
	"time"

	"github.com/waf-draft/waf/internal/config"
	"github.com/waf-draft/waf/internal/decision"
	"github.com/waf-draft/waf/internal/detection"
	"github.com/waf-draft/waf/internal/detection/rules"
	"github.com/waf-draft/waf/internal/logging"
	"github.com/waf-draft/waf/internal/mitigation"
	"github.com/waf-draft/waf/internal/normalize"
	"github.com/waf-draft/waf/internal/telemetry"
)

// WAFHandler wraps the WAF processing logic
type WAFHandler struct {
	cfg        *config.Config
	rules      []rules.Rule
	logger     *logging.Logger
	proxy      http.Handler
}

// NewWAFHandler creates a new WAF handler
func NewWAFHandler(cfg *config.Config, rules []rules.Rule, logger *logging.Logger, proxy http.Handler) *WAFHandler {
	return &WAFHandler{
		cfg:    cfg,
		rules:  rules,
		logger: logger,
		proxy:  proxy,
	}
}

// ServeHTTP implements http.Handler
func (h *WAFHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	metrics := telemetry.GetMetrics()
	
	// Track total requests
	metrics.IncrementTotalRequests()

	// Normalize request
	norm, err := normalize.Request(r, h.cfg.Security.LogRequestBody)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Evaluate request against rules
	score, matchedRules, err := detection.EvaluateRequest(r, norm, h.rules)
	if err != nil {
		// Log error but continue
	}

	// Track rule matches
	for _, rule := range matchedRules {
		metrics.IncrementRuleMatch(rule.ID)
	}

	// Make decision
	dec := decision.Decide(score, matchedRules, h.cfg)

	// Determine status code
	statusCode := http.StatusOK
	if dec.Action == "block" {
		statusCode = http.StatusForbidden
		metrics.IncrementBlockedRequests()
	} else {
		metrics.IncrementAllowedRequests()
	}

	// Track latency
	latency := time.Since(start)
	metrics.AddLatency(latency.Nanoseconds())

	// Log request
	h.logger.LogRequest(r, norm, dec, matchedRules, statusCode)

	// Apply mitigation
	if dec.Action == "block" {
		mitigation.ApplyDecision(dec, w, r, nil)
	} else {
		// Forward to upstream
		h.proxy.ServeHTTP(w, r)
	}
}

