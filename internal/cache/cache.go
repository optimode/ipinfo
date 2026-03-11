package cache

import (
	"net"
	"strings"
	"sync"

	"github.com/optimode/ipinfo/internal/api"
)

// SubnetCache caches API responses keyed by /24 subnet.
type SubnetCache struct {
	mu    sync.RWMutex
	store map[string]*api.Response
}

// New creates a new SubnetCache.
func New() *SubnetCache {
	return &SubnetCache{
		store: make(map[string]*api.Response),
	}
}

// Key returns the /24 subnet key for an IP (first 3 octets).
// Returns empty string if not a valid IPv4.
func Key(ip string) string {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return ""
	}
	parsed = parsed.To4()
	if parsed == nil {
		return "" // IPv6 – no subnet caching
	}
	parts := strings.Split(ip, ".")
	if len(parts) < 3 {
		return ""
	}
	return strings.Join(parts[:3], ".")
}

// Get returns a cached response for the given /24 key, or nil if not found.
// The returned response has Query replaced with the requested IP.
func (c *SubnetCache) Get(ip, subnetKey string) *api.Response {
	if subnetKey == "" {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if r, ok := c.store[subnetKey]; ok {
		// Return a copy with the actual IP as Query
		copy := *r
		copy.Query = ip
		return &copy
	}
	return nil
}

// Set stores a response under the given /24 key.
func (c *SubnetCache) Set(subnetKey string, r *api.Response) {
	if subnetKey == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[subnetKey] = r
}
