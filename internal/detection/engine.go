package detection

import (
	"net/http"
	"strings"

	"github.com/waf-draft/waf/internal/detection/rules"
	"github.com/waf-draft/waf/internal/normalize"
)

// EvaluateRequest evaluates a request against all rules and returns anomaly score and matched rules
func EvaluateRequest(req *http.Request, norm *normalize.NormalizedRequest, ruleSet []rules.Rule) (*AnomalyScore, []rules.Rule, error) {
	score := NewAnomalyScore()
	var matchedRules []rules.Rule

	for _, rule := range ruleSet {
		// Skip rules that don't match the current phase
		if rule.Phase != "request" && rule.Phase != "" {
			continue
		}

		matched, err := evaluateRule(req, norm, rule)
		if err != nil {
			// Log error but continue with other rules
			continue
		}

		if matched {
			matchedRules = append(matchedRules, rule)
			// Process actions
			for _, action := range rule.Actions {
				if action.Type == "add_score" {
					var scoreValue int
					switch v := action.Param.(type) {
					case int:
						scoreValue = v
					case float64:
						scoreValue = int(v)
					default:
						// Use rule severity as fallback
						scoreValue = rule.Severity
					}
					score.Add(scoreValue, rule.Tags)
				}
			}
		}
	}

	return score, matchedRules, nil
}

// evaluateRule checks if a rule matches the request
func evaluateRule(req *http.Request, norm *normalize.NormalizedRequest, rule rules.Rule) (bool, error) {
	// All conditions must match (AND logic)
	for _, condition := range rule.Conditions {
		matched, err := evaluateCondition(req, norm, condition)
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
	}

	return true, nil
}

// evaluateCondition checks if a condition matches
func evaluateCondition(req *http.Request, norm *normalize.NormalizedRequest, condition rules.MatchCondition) (bool, error) {
	var value string

	switch condition.Target {
	case "path":
		// For path traversal detection, check original path
		// For other checks, use normalized path
		if strings.Contains(condition.Value, "..") || strings.Contains(condition.Value, "%2e") {
			value = norm.OriginalPath
		} else {
			value = norm.Path
		}
	case "query":
		// Check all query parameters
		queryStr := ""
		for k, vals := range norm.Query {
			queryStr += k + "=" + strings.Join(vals, ",") + "&"
		}
		value = queryStr
	case "query_param":
		// This would need a specific parameter name, for now check all
		queryStr := ""
		for k, vals := range norm.Query {
			queryStr += k + "=" + strings.Join(vals, ",") + "&"
		}
		value = queryStr
	case "header":
		// Check all headers
		headerStr := ""
		for k, v := range norm.Headers {
			headerStr += k + ":" + v + "\n"
		}
		value = headerStr
	case "body":
		value = norm.Body
	case "method":
		value = norm.Method
	default:
		// Try to match against all fields
		allFields := norm.Path + " " + norm.Body
		for k, vals := range norm.Query {
			allFields += " " + k + "=" + strings.Join(vals, ",")
		}
		for k, v := range norm.Headers {
			allFields += " " + k + ":" + v
		}
		value = allFields
	}

	return (&condition).Match(value)
}

