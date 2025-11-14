package rules

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Rule represents a WAF detection rule
type Rule struct {
	ID         string          `json:"id" yaml:"id"`
	Name       string          `json:"name" yaml:"name"`
	Severity   int             `json:"severity" yaml:"severity"`
	Phase      string          `json:"phase" yaml:"phase"`
	Conditions []MatchCondition `json:"conditions" yaml:"conditions"`
	Actions    []Action        `json:"actions" yaml:"actions"`
	Tags       []string        `json:"tags" yaml:"tags"`
	Enabled    bool            `json:"enabled" yaml:"enabled"`
}

// MatchCondition defines a condition to match against request data
type MatchCondition struct {
	Target   string `json:"target" yaml:"target"`
	Operator string `json:"operator" yaml:"operator"`
	Value    string `json:"value" yaml:"value"`
}

// Action defines an action to take when a rule matches
type Action struct {
	Type  string      `json:"type" yaml:"type"`
	Param interface{} `json:"param" yaml:"param"`
}

// LoadRules loads rules from YAML files
func LoadRules(filePaths []string) ([]Rule, error) {
	var allRules []Rule

	for _, filePath := range filePaths {
		rules, err := loadRulesFromFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to load rules from %s: %w", filePath, err)
		}
		allRules = append(allRules, rules...)
	}

	// Filter only enabled rules
	enabledRules := make([]Rule, 0)
	for _, rule := range allRules {
		if rule.Enabled {
			enabledRules = append(enabledRules, rule)
		}
	}

	return enabledRules, nil
}

// loadRulesFromFile loads rules from a single YAML file
func loadRulesFromFile(filePath string) ([]Rule, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read rules file: %w", err)
	}

	var rules []Rule
	if err := yaml.Unmarshal(data, &rules); err != nil {
		return nil, fmt.Errorf("failed to parse rules file: %w", err)
	}

	return rules, nil
}

// Match checks if a condition matches the given value
func (c *MatchCondition) Match(value string) (bool, error) {
	switch c.Operator {
	case "equals":
		return value == c.Value, nil
	case "contains":
		return strings.Contains(strings.ToLower(value), strings.ToLower(c.Value)), nil
	case "regex":
		re, err := regexp.Compile(c.Value)
		if err != nil {
			return false, fmt.Errorf("invalid regex pattern: %w", err)
		}
		return re.MatchString(value), nil
	case "starts_with":
		return strings.HasPrefix(strings.ToLower(value), strings.ToLower(c.Value)), nil
	case "ends_with":
		return strings.HasSuffix(strings.ToLower(value), strings.ToLower(c.Value)), nil
	default:
		return false, fmt.Errorf("unknown operator: %s", c.Operator)
	}
}

