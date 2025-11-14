# Test Website for WAF Testing

A vulnerable test web application designed to test Web Application Firewall (WAF) rules.

## Features

- **Home Page** - Overview and navigation
- **Search Page** - Test XSS and injection attacks
- **Users Page** - Test SQL injection in user queries
- **Login Page** - Test authentication bypass attempts
- **Files Page** - Test path traversal attacks

## Running the Test Website

### Option 1: Direct Access (No WAF)
```bash
cd /Users/yusufdalbudak/Documents/github/WAF-DRAFT
./test-website
```
Access at: `http://localhost:8081`

### Option 2: Through WAF (Recommended)
1. Start the test website:
```bash
./test-website
```

2. Start the WAF:
```bash
go run ./cmd/wafd
```

3. Access through WAF at: `http://localhost:8080`

## Test Scenarios

### SQL Injection
- **Users Page**: Try `?id=1 OR 1=1`
- **Users Page**: Try `?id=1' UNION SELECT * FROM users--`
- **Login**: Try username `admin' OR '1'='1`

### XSS (Cross-Site Scripting)
- **Search Page**: Try `?q=<script>alert('XSS')</script>`
- **Search Page**: Try `?q=<img src=x onerror=alert('XSS')>`
- **Search Page**: Try `?q=javascript:alert('XSS')`

### Path Traversal
- **Files Page**: Try `?file=../../../etc/passwd`
- **Files Page**: Try `?file=..\..\..\windows\system32\config\sam`

### Command Injection
- **Search Page**: Try `?q=; ls -la`
- **Search Page**: Try `?q=| cat /etc/passwd`

## Pages

- `/` - Home page
- `/search` - Search functionality
- `/users` - User listing and search
- `/login` - Login form
- `/files` - File access simulation

## API Endpoints

- `GET /api/users?id=1` - Get user by ID
- `GET /api/search?q=query` - Search API
- `GET /api/login?username=admin&password=pass` - Login API

## Warning

⚠️ **This website is intentionally vulnerable for testing purposes only!**
Do not use in production or expose to the internet.

