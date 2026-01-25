package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents the API client
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ListServers retrieves all servers for the authenticated user
func (c *Client) ListServers(page int) (*ServersResponse, error) {
	var resp ServersResponse
	path := "servers"
	if page > 0 {
		path = fmt.Sprintf("servers?page=%d", page)
	}
	err := c.do("GET", path, nil, &resp)
	return &resp, err
}

// CreateServer creates a new server
func (c *Client) CreateServer(req *CreateServerRequest) (*Server, error) {
	var resp ServerResponse
	err := c.do("POST", "servers", req, &resp)
	return &resp.Server, err
}

// DeleteServer deletes a server by ID
func (c *Client) DeleteServer(id int) error {
	return c.do("DELETE", fmt.Sprintf("servers/%d", id), nil, nil)
}

// ListProviders retrieves all available cloud providers
func (c *Client) ListProviders(page int) (*ProvidersResponse, error) {
	var resp ProvidersResponse
	path := "providers"
	if page > 0 {
		path = fmt.Sprintf("providers?page=%d", page)
	}
	err := c.do("GET", path, nil, &resp)
	return &resp, err
}

// ListRegions retrieves regions for a provider
func (c *Client) ListRegions(providerID, page int) (*RegionsResponse, error) {
	var resp RegionsResponse
	path := fmt.Sprintf("regions?provider_id=%d", providerID)
	if page > 0 {
		path = fmt.Sprintf("regions?provider_id=%d&page=%d", providerID, page)
	}
	err := c.do("GET", path, nil, &resp)
	return &resp, err
}

// ListSizes retrieves available server sizes
func (c *Client) ListSizes(regionID, providerID, page int) (*SizesResponse, error) {
	var resp SizesResponse
	path := fmt.Sprintf("sizes?region_id=%d&provider_id=%d", regionID, providerID)
	if page > 0 {
		path = fmt.Sprintf("sizes?region_id=%d&provider_id=%d&page=%d", regionID, providerID, page)
	}
	err := c.do("GET", path, nil, &resp)
	return &resp, err
}

// ListImages retrieves available OS images
func (c *Client) ListImages(providerID int, architecture string, page int) (*ImagesResponse, error) {
	var resp ImagesResponse
	path := fmt.Sprintf("images?provider_id=%d&architecture=%s&type=STANDARD", providerID, architecture)
	if page > 0 {
		path = fmt.Sprintf("images?provider_id=%d&architecture=%s&type=STANDARD&page=%d", providerID, architecture, page)
	}
	err := c.do("GET", path, nil, &resp)
	return &resp, err
}

// ListSSHKeys retrieves the user's SSH keys
func (c *Client) ListSSHKeys(page int) (*KeysResponse, error) {
	var resp KeysResponse
	path := "keys"
	if page > 0 {
		path = fmt.Sprintf("keys?page=%d", page)
	}
	err := c.do("GET", path, nil, &resp)
	return &resp, err
}

// CreateSSHKey creates a new SSH key
func (c *Client) CreateSSHKey(req *CreateSSHKeyRequest) (*SSHKey, error) {
	// Wrap the request in a "key" object as expected by Rails
	wrappedReq := map[string]any{
		"key": map[string]any{
			"label": req.Label,
			"key":   req.Key,
		},
		"provider_id": req.ProviderID,
	}

	var resp KeyResponse
	err := c.do("POST", "keys", wrappedReq, &resp)
	return &resp.Key, err
}

// GetUser retrieves the user information including balance and server limit
func (c *Client) GetUser() (*UserResponse, error) {
	var resp UserResponse
	err := c.do("GET", "user", nil, &resp)
	return &resp, err
}

// do performs an HTTP request with authentication
func (c *Client) do(method, path string, body any, result any) error {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleErrorResponse(resp)
	}

	if result == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

// handleErrorResponse parses error response from the API
func (c *Client) handleErrorResponse(resp *http.Response) error {
	var errResp ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Errors:     []string{fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status)},
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Errors:     errResp.Errors,
	}
}
