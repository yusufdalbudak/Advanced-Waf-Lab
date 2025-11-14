package decision

import (
	"fmt"

	"github.com/waf-draft/waf/internal/config"
	"github.com/waf-draft/waf/internal/detection"
	"github.com/waf-draft/waf/internal/detection/rules"
)

// Decision represents the WAF decision for a request
type Decision struct {
	Action       string   `json:"action"`
	Reason       string   `json:"reason"`
	Score        int      `json:"score"`
	MatchedRules []string `json:"matched_rules"`
}

// Decide makes a decision based on anomaly score and configuration
func Decide(score *detection.AnomalyScore, matchedRules []rules.Rule, cfg *config.Config) Decision {
	decision := Decision{
		Score:        score.Total,
		MatchedRules: make([]string, 0, len(matchedRules)),
	}

	// Collect matched rule IDs
	for _, rule := range matchedRules {
		decision.MatchedRules = append(decision.MatchedRules, rule.ID)
	}

	// Make decision based on threshold
	if score.Total >= cfg.Security.AnomalyThreshold {
		decision.Action = "block"
		decision.Reason = fmt.Sprintf("Anomaly score %d exceeds threshold %d", score.Total, cfg.Security.AnomalyThreshold)
	} else {
		decision.Action = "allow"
		decision.Reason = "Request passed WAF checks"
	}

	return decision
}

