# WAF (Web Application Firewall) - Core MVP

A production-inspired, modular Web Application Firewall (WAF) implemented in Go. This WAF acts as a reverse proxy, inspecting and filtering HTTP requests based on configurable security rules before forwarding them to backend applications.

## Features

- **Reverse Proxy**: Forwards legitimate requests to upstream backend servers
- **Rule Engine**: Signature-based detection with regex, contains, and equals operators
- **Anomaly Scoring**: Each matched rule contributes to an anomaly score
- **Decision Engine**: Blocks requests when anomaly score exceeds threshold
- **Request Normalization**: URL decoding, path cleaning, query parameter normalization
- **Structured Logging**: JSON-formatted logs for all requests and decisions
- **YAML Configuration**: Easy-to-manage configuration files
- **Modular Architecture**: Clean, extensible codebase

## Quick Start

### Prerequisites

- Go 1.21 or later
- A backend application to protect (optional, for testing)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yusufdalbudak/Advanced-Waf-Lab
cd WAF-DRAFT
```

2. Install dependencies:
```bash
go mod download
```

3. Run the WAF:
```bash
go run ./cmd/wafd
```

The WAF will start on port 8080 (configurable) and forward requests to `http://localhost:8081` by default.

## Configuration

### Main Configuration (`configs/waf.yaml`)

```yaml
server:
  listen_address: ":8080"          # WAF listening address
  upstream_url: "http://localhost:8081"  # Backend server URL
  read_timeout_seconds: 10
  write_timeout_seconds: 10
  idle_timeout_seconds: 60

security:
  anomaly_threshold: 10            # Score threshold for blocking
  log_request_bodies: false        # Log request bodies (privacy consideration)

logging:
  level: "info"
  output: "stdout"                 # "stdout" or file path

rules:
  files:
    - "configs/ruleset.yaml"       # Rule files to load
```

### Rules Configuration (`configs/ruleset.yaml`)

Rules define conditions that, when matched, contribute to the anomaly score:

```yaml
- id: "SQLI-001"
  name: "Basic SQL Injection"
  severity: 10
  phase: "request"
  enabled: true
  tags: ["sqli", "injection"]
  conditions:
    - target: "query"
      operator: "regex"
      value: "(?i)(\\bor\\b\\s+1\\s*=\\s*1|union\\s+select)"
  actions:
    - type: "add_score"
      param: 10
```

#### Rule Fields

- **id**: Unique rule identifier
- **name**: Human-readable rule name
- **severity**: Base severity score (0-100)
- **phase**: Request phase ("request" for now)
- **enabled**: Whether the rule is active
- **tags**: Categories for the rule (e.g., "sqli", "xss")
- **conditions**: List of match conditions (all must match - AND logic)
- **actions**: Actions to take when rule matches

#### Condition Targets

- `path`: URL path
- `query`: All query parameters
- `header`: All headers
- `body`: Request body
- `method`: HTTP method

#### Condition Operators

- `equals`: Exact match
- `contains`: Substring match (case-insensitive)
- `regex`: Regular expression match
- `starts_with`: Prefix match
- `ends_with`: Suffix match

## How It Works

### Request Lifecycle

1. **Request Reception**: HTTP request arrives at WAF server
2. **Normalization**: Request is normalized (URL decode, path cleaning, etc.)
3. **Rule Evaluation**: Detection engine evaluates request against all enabled rules
4. **Anomaly Scoring**: Matched rules contribute to anomaly score
5. **Decision Making**: Decision engine compares score to threshold
6. **Logging**: Request and decision are logged in JSON format
7. **Mitigation**:
   - If blocked: Return 403 Forbidden
   - If allowed: Forward to upstream via reverse proxy

### Anomaly Threshold

The `anomaly_threshold` setting determines when requests are blocked:

- If `anomaly_score >= threshold` → **BLOCK** (403 Forbidden)
- If `anomaly_score < threshold` → **ALLOW** (forward to upstream)

Each rule can contribute to the score via the `add_score` action. Multiple rules can match, and their scores are cumulative.

### Example

A request with SQL injection (`?id=1 OR 1=1`) might match:
- SQLI-001: +10 points
- SQLI-003: +8 points
- **Total: 18 points**

If threshold is 10, this request is **blocked**.

## Usage Examples

### Starting the WAF

```bash
# Use default config
go run ./cmd/wafd

# Use custom config
go run ./cmd/wafd -config /path/to/config.yaml
```

### Testing with curl

```bash
# Legitimate request (should be allowed)
curl http://localhost:8080/api/users?id=123

# SQL injection (should be blocked)
curl http://localhost:8080/api/users?id=1%20OR%201=1

# XSS attempt (should be blocked)
curl "http://localhost:8080/api/search?q=<script>alert('xss')</script>"
```

### Log Output

The WAF logs all requests in JSON format:

```json
{
  "timestamp": "2024-01-15T10:30:45Z",
  "source_ip": "127.0.0.1",
  "method": "GET",
  "path": "/api/users",
  "status": 403,
  "decision": {
    "action": "block",
    "reason": "Anomaly score 18 exceeds threshold 10",
    "score": 18,
    "matched_rules": ["SQLI-001", "SQLI-003"]
  },
  "request_id": "1705315845000000000-/api/users",
  "user_agent": "curl/7.68.0",
  "query_string": "id=1 OR 1=1"
}
```

## Adding Rules

To add a new security rule:

1. Edit `configs/ruleset.yaml`
2. Add a new rule entry:

```yaml
- id: "NEW-001"
  name: "My New Rule"
  severity: 8
  phase: "request"
  enabled: true
  tags: ["custom"]
  conditions:
    - target: "query"
      operator: "contains"
      value: "malicious-pattern"
  actions:
    - type: "add_score"
      param: 8
```

3. Restart the WAF (rules are loaded at startup)

## Testing

### Run Integration Tests

```bash
go test ./test/integration/...
```

### Run Fuzz Tests

```bash
go test -fuzz=FuzzNormalizePath ./test/fuzz/...
go test -fuzz=FuzzRuleMatching ./test/fuzz/...
```

### Run All Tests

```bash
go test ./...
```

## Project Structure

```
waf/
├── cmd/
│   └── wafd/              # Main entrypoint
│       └── main.go
├── internal/
│   ├── config/            # Configuration loading
│   ├── normalize/         # Request normalization
│   ├── detection/         # Rule engine and detection
│   │   ├── rules/         # Rule definitions
│   │   ├── engine.go      # Detection engine
│   │   └── anomaly.go     # Anomaly scoring
│   ├── decision/          # Decision engine
│   ├── mitigation/        # Block/allow actions
│   ├── logging/           # Structured logging
│   ├── httpserver/        # HTTP server and handlers
│   └── telemetry/         # Metrics collection
├── configs/
│   ├── waf.yaml           # Main configuration
│   └── ruleset.yaml       # Security rules
├── test/
│   ├── integration/       # Integration tests
│   └── fuzz/              # Fuzz tests
├── docs/
│   ├── architecture.md    # System architecture
│   ├── threat-model.md    # Security threat model
│   └── waf-benchmark-plan.md  # Benchmarking plan
└── README.md
```

## Documentation

- [Architecture](docs/architecture.md): Detailed system architecture and design
- [Threat Model](docs/threat-model.md): Security analysis using STRIDE framework
- [Benchmark Plan](docs/waf-benchmark-plan.md): Testing and benchmarking strategy

## Security Considerations

- **Configuration Security**: Protect configuration files with appropriate permissions
- **Rule Updates**: Regularly update rules to detect new attack patterns
- **Logging**: Be mindful of logging sensitive data (PII, credentials)
- **Performance**: Monitor performance impact, especially with complex regex rules
- **False Positives**: Tune anomaly threshold to balance security and usability

## Limitations

This is an MVP implementation. Future enhancements could include:

- Rate limiting
- IP whitelisting/blacklisting
- Request body size limits
- Advanced evasion detection
- Machine learning-based anomaly detection
- Real-time rule updates
- Distributed deployment
- Metrics dashboard

## Contributing

This is a draft/prototype implementation. Contributions and improvements are welcome!


## Acknowledgments

- Inspired by ModSecurity and OWASP ModSecurity Core Rule Set (CRS)
- Built following OWASP WAF guidelines

