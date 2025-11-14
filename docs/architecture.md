# WAF Architecture

## Overview

This document describes the architecture of the WAF (Web Application Firewall) system. The WAF is a reverse proxy that sits between clients and backend applications, inspecting and filtering HTTP requests based on configurable security rules.

## System Architecture

```
┌─────────┐
│ Client  │
└────┬────┘
     │ HTTP Request
     ▼
┌─────────────────────────────────────┐
│         WAF Daemon (wafd)           │
│  ┌───────────────────────────────┐  │
│  │   HTTP Server (Port 8080)     │  │
│  └───────────┬───────────────────┘  │
│              │                       │
│  ┌───────────▼───────────────────┐  │
│  │      WAF Handler              │  │
│  │  ┌─────────────────────────┐  │  │
│  │  │ 1. Normalize Request    │  │  │
│  │  └──────────┬──────────────┘  │  │
│  │             │                  │  │
│  │  ┌──────────▼──────────────┐  │  │
│  │  │ 2. Detection Engine     │  │  │
│  │  │    - Rule Evaluation    │  │  │
│  │  │    - Anomaly Scoring    │  │  │
│  │  └──────────┬──────────────┘  │  │
│  │             │                  │  │
│  │  ┌──────────▼──────────────┐  │  │
│  │  │ 3. Decision Engine      │  │  │
│  │  │    - Threshold Check    │  │  │
│  │  └──────────┬──────────────┘  │  │
│  │             │                  │  │
│  │  ┌──────────▼──────────────┐  │  │
│  │  │ 4. Logging              │  │  │
│  │  └──────────┬──────────────┘  │  │
│  │             │                  │  │
│  │  ┌──────────▼──────────────┐  │  │
│  │  │ 5. Mitigation           │  │  │
│  │  │    - Block (403)        │  │  │
│  │  │    - Allow (Proxy)      │  │  │
│  │  └─────────────────────────┘  │  │
│  └───────────────────────────────┘  │
└───────────┬─────────────────────────┘
            │
            │ Allowed Request
            ▼
┌─────────────────────┐
│  Backend Server     │
│  (Port 8081)        │
└─────────────────────┘
```

## Module Structure

### 1. Configuration (`internal/config`)

- **Purpose**: Load and manage WAF configuration from YAML files
- **Key Components**:
  - `Config`: Main configuration structure
  - `ServerConfig`: HTTP server settings (listen address, timeouts, upstream URL)
  - `SecurityConfig`: Security settings (anomaly threshold, body logging)
  - `LoggingConfig`: Logging settings (level, output destination)
  - `RulesConfig`: Rule file paths

### 2. Normalization (`internal/normalize`)

- **Purpose**: Normalize and sanitize incoming HTTP requests
- **Key Functions**:
  - URL path decoding and cleaning
  - Query parameter normalization
  - Header normalization (lowercase keys)
  - Optional request body reading
- **Output**: `NormalizedRequest` struct with cleaned data

### 3. Detection Engine (`internal/detection`)

- **Purpose**: Evaluate requests against security rules
- **Components**:
  - **Rules** (`internal/detection/rules`): Rule loading and matching logic
    - Rule definition (ID, name, severity, conditions, actions, tags)
    - Match conditions (target, operator, value)
    - Operators: `equals`, `contains`, `regex`, `starts_with`, `ends_with`
  - **Anomaly Scoring** (`internal/detection/anomaly`): Track anomaly scores
    - Total score accumulation
    - Per-tag score tracking
  - **Engine** (`internal/detection/engine`): Main evaluation logic
    - Iterate through enabled rules
    - Evaluate conditions against normalized request
    - Accumulate scores based on matched rules

### 4. Decision Engine (`internal/decision`)

- **Purpose**: Make allow/block decisions based on anomaly scores
- **Logic**:
  - Compare total anomaly score against configured threshold
  - If score >= threshold → `block`
  - If score < threshold → `allow`
- **Output**: `Decision` struct with action, reason, score, and matched rule IDs

### 5. Mitigation (`internal/mitigation`)

- **Purpose**: Apply security decisions
- **Actions**:
  - **Block**: Return HTTP 403 Forbidden with JSON error response
  - **Allow**: Forward request to upstream via reverse proxy
- **Components**:
  - Reverse proxy setup for upstream forwarding
  - Block response generation

### 6. Logging (`internal/logging`)

- **Purpose**: Structured JSON logging of all requests
- **Log Events**:
  - Timestamp, source IP, method, path
  - HTTP status code
  - WAF decision (action, reason, score, matched rules)
  - Request ID, user agent, query string
- **Output**: JSON lines to stdout or file

### 7. HTTP Server (`internal/httpserver`)

- **Purpose**: HTTP server and request handling
- **Components**:
  - `Server`: HTTP server with configurable timeouts
  - `WAFHandler`: Main request handler implementing WAF pipeline
  - `Proxy`: Reverse proxy for upstream forwarding

## Request Lifecycle

1. **Request Reception**: HTTP request arrives at WAF server
2. **Normalization**: Request is normalized (URL decode, path cleaning, etc.)
3. **Rule Evaluation**: Detection engine evaluates request against all enabled rules
4. **Anomaly Scoring**: Matched rules contribute to anomaly score
5. **Decision Making**: Decision engine compares score to threshold
6. **Logging**: Request and decision are logged in JSON format
7. **Mitigation**:
   - If blocked: Return 403 Forbidden
   - If allowed: Forward to upstream via reverse proxy

## Rule Evaluation Flow

```
For each enabled rule:
  For each condition in rule:
    Extract target value (path, query, header, body, etc.)
    Apply operator (regex, contains, equals, etc.)
    If condition matches:
      Continue to next condition
    Else:
      Skip this rule (AND logic)
  
  If all conditions matched:
    Add rule to matched rules list
    Process actions (add_score, etc.)
    Accumulate anomaly score
```

## Configuration Flow

1. Load main config from `configs/waf.yaml`
2. Load rules from files specified in config
3. Filter to only enabled rules
4. Initialize logger, proxy, and server
5. Start HTTP server

## Extensibility Points

- **New Rule Types**: Add to `internal/detection/rules`
- **New Operators**: Extend `MatchCondition.Match()` method
- **New Actions**: Extend action processing in detection engine
- **New Log Fields**: Extend `LogEvent` struct
- **Custom Mitigation**: Extend `mitigation` package

## Performance Considerations

- Rules are evaluated sequentially (can be parallelized in future)
- Request body reading is optional (configurable)
- Reverse proxy uses Go's standard library for efficiency
- JSON logging is non-blocking (buffered writes)

## Security Considerations

- All rules are evaluated before decision (no short-circuit)
- Anomaly scores are cumulative (multiple rules can trigger)
- Request normalization prevents evasion techniques
- Structured logging enables security analysis

