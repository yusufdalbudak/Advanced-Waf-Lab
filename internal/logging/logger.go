package logging

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/waf-draft/waf/internal/decision"
	"github.com/waf-draft/waf/internal/detection/rules"
	"github.com/waf-draft/waf/internal/normalize"
)

// Logger handles structured JSON logging
type Logger struct {
	output *os.File
}

// NewLogger creates a new logger instance
func NewLogger(output string) (*Logger, error) {
	var file *os.File
	var err error

	if output == "stdout" || output == "" {
		file = os.Stdout
	} else {
		file, err = os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
	}

	return &Logger{output: file}, nil
}

// LogRequest logs a request event
func (l *Logger) LogRequest(req *http.Request, norm *normalize.NormalizedRequest, dec decision.Decision, matchedRules []rules.Rule, statusCode int) {
	event := LogEvent{
		Timestamp: time.Now().UTC(),
		SourceIP:  getSourceIP(req),
		Method:    req.Method,
		Path:      norm.Path,
		Status:    statusCode,
		Decision:  dec,
		RequestID: getRequestID(req),
		UserAgent: req.UserAgent(),
	}
	
	// Add severity level based on decision
	if dec.Action == "block" {
		event.Severity = "HIGH"
	} else if len(matchedRules) > 0 {
		event.Severity = "MEDIUM"
	} else {
		event.Severity = "LOW"
	}

	// Add query string if present
	if len(norm.Query) > 0 {
		queryStr := ""
		for k, vals := range norm.Query {
			for _, v := range vals {
				if queryStr != "" {
					queryStr += "&"
				}
				queryStr += k + "=" + v
			}
		}
		event.QueryString = queryStr
	}
	
	// Detect attack type from matched rules
	if len(matchedRules) > 0 {
		attackTypes := make(map[string]bool)
		for _, rule := range matchedRules {
			for _, tag := range rule.Tags {
				switch tag {
				case "sqli", "injection":
					attackTypes["SQL Injection"] = true
				case "xss":
					attackTypes["XSS"] = true
				case "path-traversal", "lfi":
					attackTypes["Path Traversal"] = true
				case "command-injection", "rce":
					attackTypes["Command Injection"] = true
				case "file-inclusion", "rfi":
					attackTypes["File Inclusion"] = true
				case "header-injection":
					attackTypes["Header Injection"] = true
				}
			}
		}
		
		types := make([]string, 0, len(attackTypes))
		for at := range attackTypes {
			types = append(types, at)
		}
		if len(types) > 0 {
			event.AttackType = strings.Join(types, ", ")
		}
	}

	l.writeJSON(event)
}

// writeJSON writes a JSON-encoded log event
func (l *Logger) writeJSON(event LogEvent) {
	encoder := json.NewEncoder(l.output)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(event); err != nil {
		// Fallback to stderr if logging fails
		fmt.Fprintf(os.Stderr, "Failed to write log: %v\n", err)
	}
}

// getSourceIP extracts the client IP from the request
func getSourceIP(req *http.Request) string {
	// Check X-Forwarded-For header
	if xff := req.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	// Check X-Real-IP header
	if xri := req.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Fallback to RemoteAddr
	return req.RemoteAddr
}

// getRequestID extracts or generates a request ID
func getRequestID(req *http.Request) string {
	// Check if request ID is already set
	if id := req.Header.Get("X-Request-ID"); id != "" {
		return id
	}
	// Generate a simple ID based on timestamp and path
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), req.URL.Path)
}

// Close closes the logger output file
func (l *Logger) Close() error {
	if l.output != os.Stdout && l.output != os.Stderr {
		return l.output.Close()
	}
	return nil
}

