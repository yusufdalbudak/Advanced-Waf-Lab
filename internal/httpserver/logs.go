package httpserver

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

// LogsHandler handles log viewing requests
type LogsHandler struct {
	logFile string
}

// NewLogsHandler creates a new logs handler
func NewLogsHandler(logFile string) *LogsHandler {
	return &LogsHandler{logFile: logFile}
}

// ServeHTTP implements http.Handler for logs
func (h *LogsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	// Get query parameters
	limit := 100 // default
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l := parseInt(limitStr); l > 0 && l <= 1000 {
			limit = l
		}
	}
	
	filter := r.URL.Query().Get("filter") // "blocked", "allowed", "all"
	if filter == "" {
		filter = "all"
	}
	
	// Read log file
	logs, err := h.readLogs(limit, filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}
	
	response := map[string]interface{}{
		"logs":   logs,
		"count":  len(logs),
		"filter": filter,
	}
	
	json.NewEncoder(w).Encode(response)
}

// readLogs reads and parses log entries
func (h *LogsHandler) readLogs(limit int, filter string) ([]map[string]interface{}, error) {
	if h.logFile == "" || h.logFile == "stdout" {
		// Can't read from stdout, return empty
		return []map[string]interface{}{}, nil
	}
	
	data, err := os.ReadFile(h.logFile)
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(string(data), "\n")
	var logs []map[string]interface{}
	
	// Read from end (most recent first)
	start := len(lines) - 1
	if start < 0 {
		start = 0
	}
	
	count := 0
	for i := start; i >= 0 && count < limit; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		
		var logEntry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &logEntry); err != nil {
			continue
		}
		
		// Apply filter
		if filter != "all" {
			decision, ok := logEntry["decision"].(map[string]interface{})
			if !ok {
				continue
			}
			action, ok := decision["action"].(string)
			if !ok {
				continue
			}
			
			if filter == "blocked" && action != "block" {
				continue
			}
			if filter == "allowed" && action != "allow" {
				continue
			}
		}
		
		logs = append(logs, logEntry)
		count++
	}
	
	return logs, nil
}

func parseInt(s string) int {
	var result int
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			return 0
		}
	}
	return result
}

