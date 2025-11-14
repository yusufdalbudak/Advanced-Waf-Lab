package logging

import (
	"time"

	"github.com/waf-draft/waf/internal/decision"
)

// LogEvent represents a structured log event
type LogEvent struct {
	Timestamp   time.Time        `json:"timestamp"`
	SourceIP    string           `json:"source_ip"`
	Method      string           `json:"method"`
	Path        string           `json:"path"`
	Status      int              `json:"status"`
	Decision    decision.Decision `json:"decision"`
	RequestID   string           `json:"request_id"`
	UserAgent   string           `json:"user_agent"`
	QueryString string           `json:"query_string,omitempty"`
	Severity    string           `json:"severity,omitempty"`
	AttackType  string           `json:"attack_type,omitempty"`
}

