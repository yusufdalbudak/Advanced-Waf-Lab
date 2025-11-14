# How to View Attack Logs

## 📊 View Logs via API

### Get All Logs
```bash
curl http://localhost:8080/logs
```

### Get Only Blocked Requests
```bash
curl "http://localhost:8080/logs?filter=blocked"
```

### Get Only Allowed Requests
```bash
curl "http://localhost:8080/logs?filter=allowed"
```

### Limit Number of Logs
```bash
curl "http://localhost:8080/logs?limit=50"
```

### Combined Filters
```bash
curl "http://localhost:8080/logs?filter=blocked&limit=20"
```

## 📝 Log File

Logs are written to: `waf.log` (in the project root)

### View Logs in Terminal
```bash
# View all logs
cat waf.log

# View only blocked attacks
cat waf.log | grep '"action":"block"'

# View recent logs (last 20 lines)
tail -20 waf.log

# Follow logs in real-time
tail -f waf.log

# Pretty print JSON logs
cat waf.log | jq '.'

# Filter by attack type
cat waf.log | jq 'select(.attack_type != null)'

# Filter by severity
cat waf.log | jq 'select(.severity == "HIGH")'
```

## 🔍 Log Structure

Each log entry contains:
```json
{
  "timestamp": "2025-11-13T21:00:00Z",
  "source_ip": "127.0.0.1",
  "method": "GET",
  "path": "/api/users",
  "status": 403,
  "severity": "HIGH",
  "attack_type": "SQL Injection",
  "decision": {
    "action": "block",
    "reason": "Anomaly score 18 exceeds threshold 10",
    "score": 18,
    "matched_rules": ["SQLI-001", "CI-002"]
  },
  "query_string": "id=1 OR 1=1",
  "user_agent": "Mozilla/5.0...",
  "request_id": "1234567890-/api/users"
}
```

## 🎯 Attack Types Detected

- **SQL Injection** - Detected from `sqli` or `injection` tags
- **XSS** - Detected from `xss` tag
- **Path Traversal** - Detected from `path-traversal` or `lfi` tags
- **Command Injection** - Detected from `command-injection` or `rce` tags
- **File Inclusion** - Detected from `file-inclusion` or `rfi` tags
- **Header Injection** - Detected from `header-injection` tag

## 📈 Severity Levels

- **HIGH** - Request was blocked
- **MEDIUM** - Request matched rules but was allowed
- **LOW** - Normal request, no rules matched

## 🧪 Test Attacks and View Logs

1. **Start WAF:**
   ```bash
   go run ./cmd/wafd
   ```

2. **Make attack requests:**
   ```bash
   # SQL Injection
   curl "http://localhost:8080/users?id=1%20OR%201=1"
   
   # XSS
   curl "http://localhost:8080/search?q=<script>alert('xss')</script>"
   
   # Command Injection
   curl "http://localhost:8080/command?cmd=;%20ls%20-la"
   ```

3. **View logs:**
   ```bash
   # Via API
   curl "http://localhost:8080/logs?filter=blocked"
   
   # Or from file
   tail -f waf.log | jq '.'
   ```

## 💡 Tips

- Use `jq` for pretty JSON formatting
- Use `grep` to filter specific patterns
- Use `tail -f` to watch logs in real-time
- Check the `severity` and `attack_type` fields for quick filtering

