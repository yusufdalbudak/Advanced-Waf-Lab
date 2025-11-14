package telemetry

import (
	"sync"
	"sync/atomic"
	"time"
)

// Metrics tracks WAF performance and security metrics
type Metrics struct {
	TotalRequests     int64
	BlockedRequests   int64
	AllowedRequests   int64
	TotalLatency      int64 // nanoseconds
	RuleMatches       sync.Map // map[string]int64
	StartTime         time.Time
}

var globalMetrics = &Metrics{
	StartTime: time.Now(),
}

// GetMetrics returns the global metrics instance
func GetMetrics() *Metrics {
	return globalMetrics
}

// IncrementTotalRequests increments the total request counter
func (m *Metrics) IncrementTotalRequests() {
	atomic.AddInt64(&m.TotalRequests, 1)
}

// IncrementBlockedRequests increments the blocked request counter
func (m *Metrics) IncrementBlockedRequests() {
	atomic.AddInt64(&m.BlockedRequests, 1)
}

// IncrementAllowedRequests increments the allowed request counter
func (m *Metrics) IncrementAllowedRequests() {
	atomic.AddInt64(&m.AllowedRequests, 1)
}

// AddLatency adds latency to the total
func (m *Metrics) AddLatency(nanoseconds int64) {
	atomic.AddInt64(&m.TotalLatency, nanoseconds)
}

// IncrementRuleMatch increments the match count for a specific rule
func (m *Metrics) IncrementRuleMatch(ruleID string) {
	for {
		value, _ := m.RuleMatches.LoadOrStore(ruleID, int64(0))
		oldCount := value.(int64)
		if m.RuleMatches.CompareAndSwap(ruleID, oldCount, oldCount+1) {
			break
		}
	}
}

// GetStats returns current statistics
func (m *Metrics) GetStats() map[string]interface{} {
	total := atomic.LoadInt64(&m.TotalRequests)
	blocked := atomic.LoadInt64(&m.BlockedRequests)
	allowed := atomic.LoadInt64(&m.AllowedRequests)
	totalLatency := atomic.LoadInt64(&m.TotalLatency)

	stats := map[string]interface{}{
		"total_requests":   total,
		"blocked_requests": blocked,
		"allowed_requests": allowed,
		"uptime_seconds":   time.Since(m.StartTime).Seconds(),
	}

	if total > 0 {
		stats["block_rate"] = float64(blocked) / float64(total)
		stats["avg_latency_ms"] = float64(totalLatency) / float64(total) / 1e6
	}

	// Collect rule match statistics
	ruleStats := make(map[string]int64)
	m.RuleMatches.Range(func(key, value interface{}) bool {
		ruleStats[key.(string)] = value.(int64)
		return true
	})
	if len(ruleStats) > 0 {
		stats["rule_matches"] = ruleStats
	}

	return stats
}

// Reset resets all metrics (useful for testing)
func (m *Metrics) Reset() {
	atomic.StoreInt64(&m.TotalRequests, 0)
	atomic.StoreInt64(&m.BlockedRequests, 0)
	atomic.StoreInt64(&m.AllowedRequests, 0)
	atomic.StoreInt64(&m.TotalLatency, 0)
	m.RuleMatches.Range(func(key, value interface{}) bool {
		m.RuleMatches.Delete(key)
		return true
	})
	m.StartTime = time.Now()
}

