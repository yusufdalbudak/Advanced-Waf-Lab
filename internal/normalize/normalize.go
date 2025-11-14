package normalize

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// NormalizedRequest represents a normalized HTTP request
type NormalizedRequest struct {
	Path        string
	OriginalPath string // Original path before normalization (for detection)
	Query       map[string][]string
	Body        string
	Method      string
	Headers     map[string]string
}

// Request normalizes an HTTP request
func Request(r *http.Request, logBody bool) (*NormalizedRequest, error) {
	norm := &NormalizedRequest{
		Method:  r.Method,
		Headers: make(map[string]string),
		Query:   make(map[string][]string),
	}

	// Store original path for detection
	norm.OriginalPath = r.URL.Path
	// Normalize path
	norm.Path = normalizePath(r.URL.Path)

	// Normalize query parameters
	norm.Query = normalizeQuery(r.URL.Query())

	// Normalize headers (lowercase keys)
	for k, v := range r.Header {
		if len(v) > 0 {
			norm.Headers[strings.ToLower(k)] = v[0]
		}
	}

	// Read body if enabled
	if logBody && r.Body != nil {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		// Restore body for downstream processing
		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		norm.Body = string(bodyBytes)
	}

	return norm, nil
}

// normalizePath cleans and normalizes the URL path
func normalizePath(rawPath string) string {
	// URL decode
	decoded, err := url.PathUnescape(rawPath)
	if err != nil {
		// If decoding fails, use original
		decoded = rawPath
	}

	// Clean path (collapse multiple slashes, resolve . and ..)
	cleaned := path.Clean(decoded)

	// Ensure it starts with /
	if !strings.HasPrefix(cleaned, "/") {
		cleaned = "/" + cleaned
	}

	return cleaned
}

// normalizeQuery normalizes query parameters
func normalizeQuery(rawQuery url.Values) map[string][]string {
	normalized := make(map[string][]string)

	for key, values := range rawQuery {
		// URL decode key
		decodedKey, err := url.QueryUnescape(key)
		if err != nil {
			decodedKey = key
		}

		// URL decode values
		decodedValues := make([]string, len(values))
		for i, v := range values {
			decoded, err := url.QueryUnescape(v)
			if err != nil {
				decoded = v
			}
			decodedValues[i] = decoded
		}

		normalized[decodedKey] = decodedValues
	}

	return normalized
}

// GetQueryString returns a single query parameter value (first if multiple)
func (n *NormalizedRequest) GetQueryString(key string) string {
	if values, ok := n.Query[key]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}

// GetQueryValues returns all values for a query parameter
func (n *NormalizedRequest) GetQueryValues(key string) []string {
	if values, ok := n.Query[key]; ok {
		return values
	}
	return nil
}

// GetHeader returns a header value
func (n *NormalizedRequest) GetHeader(key string) string {
	return n.Headers[strings.ToLower(key)]
}

