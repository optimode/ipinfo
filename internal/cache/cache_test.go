package cache

import (
	"testing"

	"github.com/optimode/ipinfo/internal/api"
)

func TestKey(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want string
	}{
		{"valid ipv4", "8.8.8.8", "8.8.8"},
		{"another ipv4", "192.168.1.100", "192.168.1"},
		{"zeros", "0.0.0.0", "0.0.0"},
		{"invalid ip", "notanip", ""},
		{"empty string", "", ""},
		{"ipv6", "::1", ""},
		{"ipv6 full", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Key(tt.ip)
			if got != tt.want {
				t.Errorf("Key(%q) = %q, want %q", tt.ip, got, tt.want)
			}
		})
	}
}

func TestSubnetCache_GetSet(t *testing.T) {
	c := New()
	resp := &api.Response{
		Query:       "1.2.3.4",
		Status:      "success",
		CountryCode: "US",
		RegionName:  "California",
		City:        "Los Angeles",
		ISP:         "TestISP",
	}

	key := "1.2.3"
	c.Set(key, resp)

	// Same IP should return the response
	got := c.Get("1.2.3.4", key)
	if got == nil {
		t.Fatal("expected cached response, got nil")
	}
	if got.Query != "1.2.3.4" {
		t.Errorf("Query = %q, want %q", got.Query, "1.2.3.4")
	}
	if got.CountryCode != "US" {
		t.Errorf("CountryCode = %q, want %q", got.CountryCode, "US")
	}

	// Different IP in same /24 should return with substituted Query
	got2 := c.Get("1.2.3.99", key)
	if got2 == nil {
		t.Fatal("expected cached response for /24 neighbor, got nil")
	}
	if got2.Query != "1.2.3.99" {
		t.Errorf("Query = %q, want %q", got2.Query, "1.2.3.99")
	}
	if got2.CountryCode != "US" {
		t.Errorf("CountryCode = %q, want %q", got2.CountryCode, "US")
	}
}

func TestSubnetCache_GetMiss(t *testing.T) {
	c := New()
	got := c.Get("1.2.3.4", "1.2.3")
	if got != nil {
		t.Errorf("expected nil for cache miss, got %+v", got)
	}
}

func TestSubnetCache_EmptyKey(t *testing.T) {
	c := New()
	resp := &api.Response{Query: "test"}

	// Set with empty key should be a no-op
	c.Set("", resp)

	// Get with empty key should return nil
	got := c.Get("1.2.3.4", "")
	if got != nil {
		t.Errorf("expected nil for empty key, got %+v", got)
	}
}

func TestSubnetCache_DoesNotMutateOriginal(t *testing.T) {
	c := New()
	resp := &api.Response{
		Query:       "1.2.3.4",
		Status:      "success",
		CountryCode: "US",
	}

	c.Set("1.2.3", resp)
	got := c.Get("1.2.3.99", "1.2.3")

	// Original should be unchanged
	if resp.Query != "1.2.3.4" {
		t.Errorf("original Query mutated to %q", resp.Query)
	}
	if got.Query != "1.2.3.99" {
		t.Errorf("cached Query = %q, want %q", got.Query, "1.2.3.99")
	}
}
