package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	baseURL = "https://pro.ip-api.com/json/%s?fields=%s&key=%s"
	fields  = "query,status,message,countryCode,regionName,city,isp,proxy,hosting,mobile"
)

// Response represents the ip-api.com JSON response.
type Response struct {
	Query       string `json:"query"`
	Status      string `json:"status"`
	Message     string `json:"message,omitempty"`
	CountryCode string `json:"countryCode"`
	RegionName  string `json:"regionName"`
	City        string `json:"city"`
	ISP         string `json:"isp"`
	Proxy       bool   `json:"proxy"`
	Hosting     bool   `json:"hosting"`
	Mobile      bool   `json:"mobile"`
}

// Client is an ip-api.com HTTP client.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// New creates a new Client.
func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Lookup queries the ip-api.com Pro API for the given IP address.
func (c *Client) Lookup(ip string) (*Response, error) {
	url := fmt.Sprintf(baseURL, ip, fields, c.apiKey)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http error for %s: %w", ip, err)
	}
	defer resp.Body.Close()

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("json decode error for %s: %w", ip, err)
	}

	if result.Status != "success" {
		return &result, fmt.Errorf("api error for %s: %s", ip, result.Message)
	}

	return &result, nil
}
