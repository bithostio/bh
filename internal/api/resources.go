package api

import "fmt"

// ListProviders retrieves all available cloud providers
func (c *Client) ListProviders() ([]Provider, error) {
	var resp ProvidersResponse
	err := c.do("GET", "/providers", nil, &resp)
	return resp.Providers, err
}

// ListRegions retrieves regions for a provider
func (c *Client) ListRegions(providerID int) ([]Region, error) {
	var resp RegionsResponse
	path := fmt.Sprintf("/regions?provider_id=%d", providerID)
	err := c.do("GET", path, nil, &resp)
	return resp.Regions, err
}

// ListSizes retrieves available server sizes
func (c *Client) ListSizes(regionID, providerID int) ([]Size, error) {
	var resp SizesResponse
	path := fmt.Sprintf("/sizes?region_id=%d&provider_id=%d", regionID, providerID)
	err := c.do("GET", path, nil, &resp)
	return resp.Sizes, err
}

// ListImages retrieves available OS images
func (c *Client) ListImages(providerID int, architecture string) ([]Image, error) {
	var resp ImagesResponse
	path := fmt.Sprintf("/images?provider_id=%d&architecture=%s&type=STANDARD", providerID, architecture)
	err := c.do("GET", path, nil, &resp)
	return resp.Images, err
}

// ListSSHKeys retrieves the user's SSH keys
func (c *Client) ListSSHKeys() ([]SSHKey, error) {
	var resp KeysResponse
	err := c.do("GET", "/keys", nil, &resp)
	return resp.Keys, err
}
