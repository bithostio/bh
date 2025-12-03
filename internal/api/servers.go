package api

import "fmt"

// ListServers retrieves all servers for the authenticated user
func (c *Client) ListServers() ([]Server, error) {
	var resp ServersResponse
	err := c.do("GET", "/servers", nil, &resp)
	return resp.Servers, err
}

// CreateServer creates a new server
func (c *Client) CreateServer(req *CreateServerRequest) (*Server, error) {
	var resp ServerResponse
	err := c.do("POST", "/servers", req, &resp)
	return &resp.Server, err
}

// DeleteServer deletes a server by ID
func (c *Client) DeleteServer(id int) error {
	return c.do("DELETE", fmt.Sprintf("/servers/%d", id), nil, nil)
}
