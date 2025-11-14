package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

const dashboardHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>WAF Dashboard</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: #1a1a1a;
            color: #e0e0e0;
            margin: 0;
            padding: 20px;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        h1 {
            color: #4CAF50;
            border-bottom: 2px solid #4CAF50;
            padding-bottom: 10px;
        }
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin: 20px 0;
        }
        .stat-card {
            background: #2a2a2a;
            padding: 20px;
            border-radius: 8px;
            border-left: 4px solid #4CAF50;
        }
        .stat-value {
            font-size: 2em;
            font-weight: bold;
            color: #4CAF50;
        }
        .stat-label {
            color: #aaa;
            margin-top: 5px;
        }
        .rules {
            background: #2a2a2a;
            padding: 20px;
            border-radius: 8px;
            margin-top: 20px;
        }
        .rule-item {
            display: flex;
            justify-content: space-between;
            padding: 10px;
            border-bottom: 1px solid #3a3a3a;
        }
        .rule-item:last-child {
            border-bottom: none;
        }
        .refresh-btn {
            background: #4CAF50;
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 5px;
            cursor: pointer;
            font-size: 16px;
            margin: 10px 0;
        }
        .refresh-btn:hover {
            background: #45a049;
        }
        .status {
            display: inline-block;
            padding: 5px 10px;
            border-radius: 3px;
            font-size: 0.9em;
        }
        .status.healthy {
            background: #4CAF50;
            color: white;
        }
        .attacks-section {
            background: #2a2a2a;
            padding: 20px;
            border-radius: 8px;
            margin-top: 20px;
        }
        .attacks-section h2 {
            color: #4CAF50;
            margin-bottom: 15px;
        }
        .filter-buttons {
            display: flex;
            gap: 10px;
            margin-bottom: 15px;
        }
        .filter-btn {
            padding: 8px 16px;
            background: #3a3a3a;
            color: #e0e0e0;
            border: 1px solid #4CAF50;
            border-radius: 5px;
            cursor: pointer;
        }
        .filter-btn:hover {
            background: #4a4a4a;
        }
        .filter-btn.active {
            background: #4CAF50;
            color: white;
        }
        .attack-logs {
            max-height: 500px;
            overflow-y: auto;
        }
        .logs-header {
            display: grid;
            grid-template-columns: 100px 120px 200px 1fr 60px 100px;
            gap: 10px;
            padding: 10px;
            background: #1a1a1a;
            font-weight: bold;
            color: #4CAF50;
            border-bottom: 2px solid #4CAF50;
            position: sticky;
            top: 0;
        }
        .log-entry {
            display: grid;
            grid-template-columns: 100px 120px 200px 1fr 60px 100px;
            gap: 10px;
            padding: 10px;
            border-bottom: 1px solid #3a3a3a;
            cursor: pointer;
            transition: background 0.2s;
        }
        .log-entry:hover {
            background: #3a3a3a;
        }
        .log-entry.blocked {
            border-left: 4px solid #f44336;
        }
        .log-entry.allowed {
            border-left: 4px solid #4CAF50;
        }
        .log-type {
            font-weight: bold;
        }
        .log-type.high {
            color: #f44336;
        }
        .log-type.medium {
            color: #ff9800;
        }
        .log-type.low {
            color: #4CAF50;
        }
        .log-status.block {
            color: #f44336;
            font-weight: bold;
        }
        .log-status.allow {
            color: #4CAF50;
        }
        .log-score {
            text-align: center;
            font-weight: bold;
        }
    </style>
    <script>
        function refreshData() {
            fetch('/api/metrics')
                .then(r => r.json())
                .then(data => {
                    document.getElementById('total-requests').textContent = data.total_requests || 0;
                    document.getElementById('blocked-requests').textContent = data.blocked_requests || 0;
                    document.getElementById('allowed-requests').textContent = data.allowed_requests || 0;
                    document.getElementById('block-rate').textContent = ((data.block_rate || 0) * 100).toFixed(1) + '%';
                    document.getElementById('avg-latency').textContent = (data.avg_latency_ms || 0).toFixed(2) + ' ms';
                    document.getElementById('uptime').textContent = Math.floor(data.uptime_seconds || 0) + 's';
                    
                    const rulesDiv = document.getElementById('rules');
                    if (data.rule_matches && Object.keys(data.rule_matches).length > 0) {
                        rulesDiv.innerHTML = '<h3>Rule Matches</h3>';
                        for (const [rule, count] of Object.entries(data.rule_matches)) {
                            const div = document.createElement('div');
                            div.className = 'rule-item';
                            div.innerHTML = '<span>' + rule + '</span><span>' + count + '</span>';
                            rulesDiv.appendChild(div);
                        }
                    } else {
                        rulesDiv.innerHTML = '<p>No rule matches yet</p>';
                    }
                });
            
            fetch('/api/health')
                .then(r => r.json())
                .then(data => {
                    document.getElementById('status').textContent = data.status || 'unknown';
                    document.getElementById('status').className = 'status ' + (data.status || 'unknown');
                });
        }
        
        let currentFilter = 'blocked';
        
        function filterLogs(filter) {
            currentFilter = filter;
            document.querySelectorAll('.filter-btn').forEach(btn => btn.classList.remove('active'));
            event.target.classList.add('active');
            loadAttackLogs();
        }
        
        function loadAttackLogs() {
            fetch('/api/logs?filter=' + currentFilter + '&limit=20')
                .then(r => r.json())
                .then(data => {
                    const logsDiv = document.getElementById('attack-logs');
                    if (data.logs && data.logs.length > 0) {
                        logsDiv.innerHTML = '<div class="logs-header"><span>Time</span><span>IP</span><span>Attack Type</span><span>Path</span><span>Score</span><span>Status</span></div>';
                        data.logs.forEach(log => {
                            const div = document.createElement('div');
                            div.className = 'log-entry ' + (log.decision.action === 'block' ? 'blocked' : 'allowed');
                            
                            const time = new Date(log.timestamp).toLocaleTimeString();
                            const attackType = log.attack_type || 'N/A';
                            const severity = log.severity || 'LOW';
                            const score = log.decision.score || 0;
                            const status = log.decision.action === 'block' ? '🚫 BLOCKED' : '✅ ALLOWED';
                            
                            div.innerHTML = '<span class="log-time">' + time + '</span>' +
                                '<span class="log-ip">' + log.source_ip.split(':')[0] + '</span>' +
                                '<span class="log-type ' + severity.toLowerCase() + '">' + attackType + '</span>' +
                                '<span class="log-path">' + log.path + '</span>' +
                                '<span class="log-score">' + score + '</span>' +
                                '<span class="log-status ' + log.decision.action + '">' + status + '</span>';
                            
                            // Add click to show details
                            div.onclick = () => showLogDetails(log);
                            logsDiv.appendChild(div);
                        });
                    } else {
                        logsDiv.innerHTML = '<p>No logs found</p>';
                    }
                })
                .catch(err => {
                    document.getElementById('attack-logs').innerHTML = '<p>Error loading logs</p>';
                });
        }
        
        function showLogDetails(log) {
            const details = '<strong>Request Details:</strong><br>' +
                'Time: ' + new Date(log.timestamp).toLocaleString() + '<br>' +
                'IP: ' + log.source_ip + '<br>' +
                'Method: ' + log.method + '<br>' +
                'Path: ' + log.path + '<br>' +
                (log.query_string ? 'Query: ' + log.query_string + '<br>' : '') +
                '<br>' +
                '<strong>Attack Info:</strong><br>' +
                'Type: ' + (log.attack_type || 'None') + '<br>' +
                'Severity: ' + (log.severity || 'LOW') + '<br>' +
                'Score: ' + log.decision.score + '<br>' +
                'Status: ' + log.decision.action.toUpperCase() + '<br>' +
                '<br>' +
                '<strong>Matched Rules:</strong><br>' +
                (log.decision.matched_rules.join(', ') || 'None') + '<br>' +
                '<br>' +
                '<strong>Reason:</strong><br>' +
                log.decision.reason;
            alert(details);
        }
        
        // Refresh every 2 seconds
        setInterval(refreshData, 2000);
        setInterval(loadAttackLogs, 2000);
        // Initial load
        refreshData();
        loadAttackLogs();
    </script>
</head>
<body>
    <div class="container">
        <h1>🛡️ WAF Dashboard</h1>
        <button class="refresh-btn" onclick="refreshData()">Refresh</button>
        
        <div class="stats">
            <div class="stat-card">
                <div class="stat-value" id="total-requests">0</div>
                <div class="stat-label">Total Requests</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="blocked-requests">0</div>
                <div class="stat-label">Blocked Requests</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="allowed-requests">0</div>
                <div class="stat-label">Allowed Requests</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="block-rate">0%</div>
                <div class="stat-label">Block Rate</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="avg-latency">0 ms</div>
                <div class="stat-label">Avg Latency</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="uptime">0s</div>
                <div class="stat-label">Uptime</div>
            </div>
        </div>
        
        <div class="stat-card">
            <strong>Status:</strong> <span class="status healthy" id="status">healthy</span>
        </div>
        
        <div class="rules" id="rules">
            <p>Loading rule matches...</p>
        </div>
        
        <div class="attacks-section">
            <h2>🛡️ Recent Attack Logs</h2>
            <div class="filter-buttons">
                <button class="filter-btn active" onclick="filterLogs('blocked')">Blocked Attacks</button>
                <button class="filter-btn" onclick="filterLogs('all')">All Requests</button>
                <button class="filter-btn" onclick="filterLogs('allowed')">Allowed</button>
            </div>
            <div id="attack-logs" class="attack-logs">
                <p>Loading attack logs...</p>
            </div>
        </div>
    </div>
</body>
</html>
`

func main() {
	wafURL := "http://localhost:8080"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("dashboard").Parse(dashboardHTML))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/api/metrics", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(wafURL + "/metrics")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.Copy(w, resp.Body)
	})

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(wafURL + "/health")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.Copy(w, resp.Body)
	})

	http.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		filter := r.URL.Query().Get("filter")
		limit := r.URL.Query().Get("limit")
		url := wafURL + "/logs"
		if filter != "" {
			url += "?filter=" + filter
			if limit != "" {
				url += "&limit=" + limit
			}
		} else if limit != "" {
			url += "?limit=" + limit
		}

		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		io.Copy(w, resp.Body)
	})

	port := ":8082"
	fmt.Printf("WAF Dashboard starting on http://localhost%s\n", port)
	fmt.Printf("Connecting to WAF at %s\n", wafURL)
	log.Fatal(http.ListenAndServe(port, nil))
}
