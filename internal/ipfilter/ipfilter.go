package ipfilter

import (
	"net"
	"sync"
)

// IPFilter manages IP whitelist and blacklist
type IPFilter struct {
	whitelist map[string]bool
	blacklist map[string]bool
	mu        sync.RWMutex
}

// NewIPFilter creates a new IP filter
func NewIPFilter() *IPFilter {
	return &IPFilter{
		whitelist: make(map[string]bool),
		blacklist: make(map[string]bool),
	}
}

// AddToWhitelist adds an IP or CIDR to the whitelist
func (f *IPFilter) AddToWhitelist(ipOrCIDR string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Check if it's a CIDR
	if _, ipnet, err := net.ParseCIDR(ipOrCIDR); err == nil {
		// For CIDR, we'll store the network
		f.whitelist[ipnet.String()] = true
		return nil
	}

	// Single IP
	if ip := net.ParseIP(ipOrCIDR); ip != nil {
		f.whitelist[ip.String()] = true
		return nil
	}

	return nil
}

// AddToBlacklist adds an IP or CIDR to the blacklist
func (f *IPFilter) AddToBlacklist(ipOrCIDR string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// Check if it's a CIDR
	if _, ipnet, err := net.ParseCIDR(ipOrCIDR); err == nil {
		f.blacklist[ipnet.String()] = true
		return nil
	}

	// Single IP
	if ip := net.ParseIP(ipOrCIDR); ip != nil {
		f.blacklist[ip.String()] = true
		return nil
	}

	return nil
}

// IsWhitelisted checks if an IP is whitelisted
func (f *IPFilter) IsWhitelisted(ipStr string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// Check exact match
	if f.whitelist[ip.String()] {
		return true
	}

	// Check CIDR matches
	for cidrStr := range f.whitelist {
		if _, ipnet, err := net.ParseCIDR(cidrStr); err == nil {
			if ipnet.Contains(ip) {
				return true
			}
		}
	}

	return false
}

// IsBlacklisted checks if an IP is blacklisted
func (f *IPFilter) IsBlacklisted(ipStr string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// Check exact match
	if f.blacklist[ip.String()] {
		return true
	}

	// Check CIDR matches
	for cidrStr := range f.blacklist {
		if _, ipnet, err := net.ParseCIDR(cidrStr); err == nil {
			if ipnet.Contains(ip) {
				return true
			}
		}
	}

	return false
}

// RemoveFromWhitelist removes an IP from whitelist
func (f *IPFilter) RemoveFromWhitelist(ipOrCIDR string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.whitelist, ipOrCIDR)
}

// RemoveFromBlacklist removes an IP from blacklist
func (f *IPFilter) RemoveFromBlacklist(ipOrCIDR string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.blacklist, ipOrCIDR)
}

// GetStats returns filter statistics
func (f *IPFilter) GetStats() map[string]interface{} {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return map[string]interface{}{
		"whitelist_count": len(f.whitelist),
		"blacklist_count": len(f.blacklist),
	}
}

