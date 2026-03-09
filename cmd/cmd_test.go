package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/bithostio/bh/internal/api"
	"github.com/gsamokovarov/assert"
)

func executeCommand(args ...string) (string, error) {
	var execErr error
	output := captureOutput(func() {
		rootCmd.SetArgs(args)
		execErr = rootCmd.Execute()
	})

	return output, execErr
}

func TestVersionCommand(t *testing.T) {
	t.Run("prints version", func(t *testing.T) {
		Version = "0.1.0"
		defer func() { Version = "dev" }()

		output, err := executeCommand("version")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "0.1.0"))
	})

	t.Run("prints dev when unset", func(t *testing.T) {
		Version = "dev"

		output, err := executeCommand("version")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "dev"))
	})
}

func TestProvidersCommand(t *testing.T) {
	t.Setenv("BH_API_KEY", "test-api-key")

	t.Run("calls correct API endpoint", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(200, api.ProvidersResponse{
					Providers: []api.Provider{
						{ID: 1, Name: "Hetzner", Slug: "hetzner"},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("providers")
		assert.Nil(t, err)
		assert.NotNil(t, capturedReq)
		assert.Equal(t, "GET", capturedReq.Method)
		assert.True(t, strings.HasSuffix(capturedReq.URL.Path, "providers"))
		assert.True(t, strings.HasPrefix(capturedReq.Header.Get("Authorization"), "Bearer "))
	})

	t.Run("passes pagination parameter", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(200, api.ProvidersResponse{
					Providers: []api.Provider{},
					Meta:      &api.Meta{Pagination: api.Pagination{Page: 2}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("providers", "--page", "2")
		assert.Nil(t, err)
		assert.Equal(t, "page=2", capturedReq.URL.RawQuery)
	})

	t.Run("handles empty response", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.ProvidersResponse{
					Providers: []api.Provider{},
					Meta:      &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("providers")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "No providers found"))
	})

	t.Run("handles API error", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(401, api.ErrorResponse{
					Errors: []string{"Invalid API key"},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("providers")
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "Authentication failed"))
	})
}

func TestServersListCommand(t *testing.T) {
	t.Setenv("BH_API_KEY", "test-api-key")

	t.Run("calls correct API endpoint", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(200, api.ServersResponse{
					Servers: []api.Server{
						{ID: 1, Name: "web-server", Status: "active", Power: true},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("servers", "list")
		assert.Nil(t, err)
		assert.NotNil(t, capturedReq)
		assert.Equal(t, "GET", capturedReq.Method)
		assert.True(t, strings.HasSuffix(capturedReq.URL.Path, "servers"))
	})

	t.Run("passes pagination parameter", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(200, api.ServersResponse{
					Servers: []api.Server{},
					Meta:    &api.Meta{Pagination: api.Pagination{Page: 2}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("servers", "list", "--page", "2")
		assert.Nil(t, err)
		assert.Equal(t, "page=2", capturedReq.URL.RawQuery)
	})

	t.Run("handles empty response", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.ServersResponse{
					Servers: []api.Server{},
					Meta:    &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("servers", "list")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "servers found"))
	})

	t.Run("handles API error", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(500, api.ErrorResponse{
					Errors: []string{"Internal server error"},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("servers", "list")
		assert.NotNil(t, err)
	})
}

func TestServerDeleteCommand(t *testing.T) {
	t.Setenv("BH_API_KEY", "test-api-key")

	t.Run("calls correct API endpoint", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(204, nil), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("servers", "delete", "123", "--force")
		assert.Nil(t, err)
		assert.NotNil(t, capturedReq)
		assert.Equal(t, "DELETE", capturedReq.Method)
		assert.True(t, strings.HasSuffix(capturedReq.URL.Path, "servers/123"))
	})

	t.Run("handles not found error", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(404, api.ErrorResponse{
					Errors: []string{"Server not found"},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("servers", "delete", "999", "--force")
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "not found"))
	})
}

func TestUserCommand(t *testing.T) {
	t.Setenv("BH_API_KEY", "test-api-key")

	t.Run("displays user info", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				assert.True(t, strings.HasSuffix(req.URL.Path, "user"))
				return jsonResponse(200, api.UserResponse{
					FullName:    "John Doe",
					Email:       "john@example.com",
					Balance:     25.50,
					ServerLimit: 10,
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("user")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "John Doe"))
		assert.True(t, strings.Contains(output, "john@example.com"))
		assert.True(t, strings.Contains(output, "25.50"))
		assert.True(t, strings.Contains(output, "10"))
	})

	t.Run("hides server limit when zero", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.UserResponse{
					FullName:    "Admin User",
					Email:       "admin@example.com",
					Balance:     100.00,
					ServerLimit: 0,
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("user")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "Admin User"))
		assert.True(t, !strings.Contains(output, "Server Limit"))
	})

	t.Run("handles API error", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(401, api.ErrorResponse{
					Errors: []string{"Unauthorized"},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("user")
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "Authentication failed"))
	})
}

func TestRegionsCommand(t *testing.T) {
	t.Setenv("BH_API_KEY", "test-api-key")

	t.Run("requires provider flag", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.RegionsResponse{
					Regions: []api.Region{},
					Meta:    &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("regions")
		assert.NotNil(t, err)
	})

	t.Run("passes provider to API", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(200, api.RegionsResponse{
					Regions: []api.Region{
						{ID: 1, Name: "Frankfurt", Slug: "fsn1"},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("regions", "--provider", "hetzner")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(capturedReq.URL.RawQuery, "provider=hetzner"))
	})

	t.Run("displays multiple regions", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.RegionsResponse{
					Regions: []api.Region{
						{ID: 1, Name: "Frankfurt", Slug: "fsn1", Origin: "EU"},
						{ID: 2, Name: "Helsinki", Slug: "hel1", Origin: "EU"},
						{ID: 3, Name: "Ashburn", Slug: "ash", Origin: "US"},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("regions", "--provider", "hetzner")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "Frankfurt"))
		assert.True(t, strings.Contains(output, "Helsinki"))
		assert.True(t, strings.Contains(output, "Ashburn"))
	})

	t.Run("handles empty response", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.RegionsResponse{
					Regions: []api.Region{},
					Meta:    &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("regions", "--provider", "hetzner")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "No regions found"))
	})

	t.Run("passes pagination parameter", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(200, api.RegionsResponse{
					Regions: []api.Region{},
					Meta:    &api.Meta{Pagination: api.Pagination{Page: 3}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("regions", "--provider", "hetzner", "--page", "3")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(capturedReq.URL.RawQuery, "page=3"))
	})

	t.Run("handles API error", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(401, api.ErrorResponse{
					Errors: []string{"Invalid API key"},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("regions", "--provider", "hetzner")
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "Authentication failed"))
	})
}

func TestSizesCommand(t *testing.T) {
	t.Setenv("BH_API_KEY", "test-api-key")

	t.Run("requires provider flag", func(t *testing.T) {
		_, err := executeCommand("sizes", "--region", "1")
		assert.NotNil(t, err)
	})

	t.Run("requires region flag", func(t *testing.T) {
		_, err := executeCommand("sizes", "--provider", "hetzner")
		assert.NotNil(t, err)
	})

	t.Run("calls correct API endpoint", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(200, api.SizesResponse{
					Sizes: []api.Size{
						{ID: 1, Name: "CX11", Memory: "2048", Processor: "1 vCPU", Disk: "20480", Price: 4.15},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("sizes", "--provider", "hetzner", "--region", "1")
		assert.Nil(t, err)
		assert.NotNil(t, capturedReq)
		assert.Equal(t, "GET", capturedReq.Method)
		assert.True(t, strings.HasSuffix(capturedReq.URL.Path, "sizes"))
		assert.True(t, strings.Contains(capturedReq.URL.RawQuery, "provider=hetzner"))
		assert.True(t, strings.Contains(capturedReq.URL.RawQuery, "region_id=1"))
	})

	t.Run("displays multiple sizes", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.SizesResponse{
					Sizes: []api.Size{
						{ID: 1, Name: "CX11", Memory: "2048", Processor: "1 vCPU", Disk: "20480", Price: 4.15},
						{ID: 2, Name: "CX21", Memory: "4096", Processor: "2 vCPU", Disk: "40960", Price: 5.83},
						{ID: 3, Name: "CX31", Memory: "8192", Processor: "2 vCPU", Disk: "81920", Price: 10.49},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("sizes", "--provider", "hetzner", "--region", "1")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "CX11"))
		assert.True(t, strings.Contains(output, "CX21"))
		assert.True(t, strings.Contains(output, "CX31"))
	})

	t.Run("handles empty response", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.SizesResponse{
					Sizes: []api.Size{},
					Meta:  &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("sizes", "--provider", "hetzner", "--region", "1")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "No sizes found"))
	})

	t.Run("passes pagination parameter", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(200, api.SizesResponse{
					Sizes: []api.Size{},
					Meta:  &api.Meta{Pagination: api.Pagination{Page: 2}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("sizes", "--provider", "hetzner", "--region", "1", "--page", "2")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(capturedReq.URL.RawQuery, "page=2"))
	})

	t.Run("handles API error", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(500, api.ErrorResponse{
					Errors: []string{"Internal server error"},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("sizes", "--provider", "hetzner", "--region", "1")
		assert.NotNil(t, err)
	})
}

func TestImagesCommand(t *testing.T) {
	t.Setenv("BH_API_KEY", "test-api-key")

	t.Run("requires provider flag", func(t *testing.T) {
		_, err := executeCommand("images")
		assert.NotNil(t, err)
	})

	t.Run("calls correct API endpoint", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(200, api.ImagesResponse{
					Images: []api.Image{
						{ID: 1, Name: "Ubuntu 22.04", Distribution: "Ubuntu", Architecture: "x86"},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("images", "--provider", "hetzner")
		assert.Nil(t, err)
		assert.NotNil(t, capturedReq)
		assert.Equal(t, "GET", capturedReq.Method)
		assert.True(t, strings.HasSuffix(capturedReq.URL.Path, "images"))
		assert.True(t, strings.Contains(capturedReq.URL.RawQuery, "provider=hetzner"))
	})

	t.Run("passes architecture parameter", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(200, api.ImagesResponse{
					Images: []api.Image{},
					Meta:   &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("images", "--provider", "hetzner", "--arch", "arm")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(capturedReq.URL.RawQuery, "architecture=arm"))
	})

	t.Run("displays multiple images", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.ImagesResponse{
					Images: []api.Image{
						{ID: 1, Name: "Ubuntu 22.04", Distribution: "Ubuntu", Architecture: "x86"},
						{ID: 2, Name: "Debian 12", Distribution: "Debian", Architecture: "x86"},
						{ID: 3, Name: "CentOS 9", Distribution: "CentOS", Architecture: "x86"},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("images", "--provider", "hetzner")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "Ubuntu 22.04"))
		assert.True(t, strings.Contains(output, "Debian 12"))
		assert.True(t, strings.Contains(output, "CentOS 9"))
	})

	t.Run("handles empty response", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.ImagesResponse{
					Images: []api.Image{},
					Meta:   &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("images", "--provider", "hetzner")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "No images found"))
	})

	t.Run("passes pagination parameter", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(200, api.ImagesResponse{
					Images: []api.Image{},
					Meta:   &api.Meta{Pagination: api.Pagination{Page: 2}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("images", "--provider", "hetzner", "--page", "2")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(capturedReq.URL.RawQuery, "page=2"))
	})

	t.Run("handles API error", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(401, api.ErrorResponse{
					Errors: []string{"Unauthorized"},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("images", "--provider", "hetzner")
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "Authentication failed"))
	})
}

func TestSSHKeysListCommand(t *testing.T) {
	t.Setenv("BH_API_KEY", "test-api-key")

	t.Run("calls correct API endpoint", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(200, api.KeysResponse{
					Keys: []api.SSHKey{
						{ID: 1, Label: "My Key", Key: "ssh-rsa AAAA..."},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("ssh-keys", "list")
		assert.Nil(t, err)
		assert.NotNil(t, capturedReq)
		assert.Equal(t, "GET", capturedReq.Method)
		assert.True(t, strings.HasSuffix(capturedReq.URL.Path, "keys"))
	})

	t.Run("displays multiple keys", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.KeysResponse{
					Keys: []api.SSHKey{
						{ID: 1, Label: "Work Laptop", Key: "ssh-rsa AAAA1..."},
						{ID: 2, Label: "Home Desktop", Key: "ssh-ed25519 BBBB2..."},
						{ID: 3, Label: "CI Server", Key: "ssh-rsa CCCC3..."},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("ssh-keys", "list")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "Work Laptop"))
		assert.True(t, strings.Contains(output, "Home Desktop"))
		assert.True(t, strings.Contains(output, "CI Server"))
	})

	t.Run("passes pagination parameter", func(t *testing.T) {
		var capturedReq *http.Request
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				return jsonResponse(200, api.KeysResponse{
					Keys: []api.SSHKey{},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 2}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("ssh-keys", "list", "--page", "2")
		assert.Nil(t, err)
		assert.Equal(t, "page=2", capturedReq.URL.RawQuery)
	})

	t.Run("handles empty response", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.KeysResponse{
					Keys: []api.SSHKey{},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("ssh-keys", "list")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "No SSH keys found"))
	})

	t.Run("handles API error", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(500, api.ErrorResponse{
					Errors: []string{"Internal server error"},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("ssh-keys", "list")
		assert.NotNil(t, err)
	})
}

func TestServersListAllFlag(t *testing.T) {
	t.Setenv("BH_API_KEY", "test-api-key")

	t.Run("filters to active servers by default", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.ServersResponse{
					Servers: []api.Server{
						{ID: 1, Name: "active-server", Status: "active", Power: true},
						{ID: 2, Name: "failed-server", Status: "failed"},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("servers", "list")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "active-server"))
		assert.True(t, !strings.Contains(output, "failed-server"))
	})

	t.Run("shows all servers with --all flag", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.ServersResponse{
					Servers: []api.Server{
						{ID: 1, Name: "active-server", Status: "active", Power: true},
						{ID: 2, Name: "failed-server", Status: "failed"},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("servers", "list", "--all")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "active-server"))
		assert.True(t, strings.Contains(output, "failed-server"))
	})

	t.Run("includes pending servers by default", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.ServersResponse{
					Servers: []api.Server{
						{ID: 1, Name: "pending-server", Status: "active", Pending: true},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("servers", "list")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "pending-server"))
	})

	t.Run("displays multiple servers", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.ServersResponse{
					Servers: []api.Server{
						{ID: 1, Name: "web-1", Status: "active", Power: true, IPAddress: "1.2.3.4", Provider: "hetzner"},
						{ID: 2, Name: "web-2", Status: "active", Power: true, IPAddress: "5.6.7.8", Provider: "hetzner"},
						{ID: 3, Name: "db-1", Status: "active", Power: true, IPAddress: "9.10.11.12", Provider: "digital_ocean"},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("servers", "list")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "web-1"))
		assert.True(t, strings.Contains(output, "web-2"))
		assert.True(t, strings.Contains(output, "db-1"))
	})
}

func TestProvidersMultiple(t *testing.T) {
	t.Setenv("BH_API_KEY", "test-api-key")

	t.Run("displays multiple providers", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, api.ProvidersResponse{
					Providers: []api.Provider{
						{ID: 1, Name: "Hetzner", Slug: "hetzner"},
						{ID: 2, Name: "DigitalOcean", Slug: "digital_ocean"},
						{ID: 3, Name: "Vultr", Slug: "vultr"},
					},
					Meta: &api.Meta{Pagination: api.Pagination{Page: 1}},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		output, err := executeCommand("providers")
		assert.Nil(t, err)
		assert.True(t, strings.Contains(output, "Hetzner"))
		assert.True(t, strings.Contains(output, "DigitalOcean"))
		assert.True(t, strings.Contains(output, "Vultr"))
	})
}

func TestServersNewCommand(t *testing.T) {
	t.Setenv("BH_API_KEY", "test-api-key")

	t.Run("requires all flags", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(200, nil), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("servers", "new")
		assert.NotNil(t, err)

		_, err = executeCommand("servers", "new", "--name", "test")
		assert.NotNil(t, err)

		_, err = executeCommand("servers", "new", "--name", "test", "--provider", "hetzner")
		assert.NotNil(t, err)

		_, err = executeCommand("servers", "new", "--name", "test", "--provider", "hetzner", "--region", "1")
		assert.NotNil(t, err)

		_, err = executeCommand("servers", "new", "--name", "test", "--provider", "hetzner", "--region", "1", "--size", "1")
		assert.NotNil(t, err)

		_, err = executeCommand("servers", "new", "--name", "test", "--provider", "hetzner", "--region", "1", "--size", "1", "--image", "1")
		assert.NotNil(t, err)
	})

	t.Run("calls correct API endpoint with body", func(t *testing.T) {
		var capturedReq *http.Request
		var capturedBody api.CreateServerRequest
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				capturedReq = req
				json.NewDecoder(req.Body).Decode(&capturedBody)
				return jsonResponse(200, api.ServerResponse{
					Server: api.Server{ID: 42, Name: "my-server"},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("servers", "new",
			"--name", "my-server",
			"--provider", "hetzner",
			"--region", "1",
			"--size", "2",
			"--image", "3",
			"--keys", "4,5",
		)
		assert.Nil(t, err)
		assert.NotNil(t, capturedReq)
		assert.Equal(t, "POST", capturedReq.Method)
		assert.True(t, strings.HasSuffix(capturedReq.URL.Path, "servers"))
		assert.Equal(t, "my-server", capturedBody.Name)
		assert.Equal(t, "hetzner", capturedBody.Provider)
		assert.Equal(t, 1, capturedBody.RegionID)
		assert.Equal(t, 2, capturedBody.SizeID)
		assert.Equal(t, 3, capturedBody.ImageID)
		assert.Equal(t, true, capturedBody.Terms)
	})

	t.Run("passes backups flag", func(t *testing.T) {
		var capturedBody api.CreateServerRequest
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				json.NewDecoder(req.Body).Decode(&capturedBody)
				return jsonResponse(200, api.ServerResponse{
					Server: api.Server{ID: 42, Name: "my-server"},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("servers", "new",
			"--name", "my-server",
			"--provider", "hetzner",
			"--region", "1",
			"--size", "2",
			"--image", "3",
			"--keys", "4",
			"--backups",
		)
		assert.Nil(t, err)
		assert.Equal(t, true, capturedBody.BackupsEnabled)
	})

	t.Run("handles API error", func(t *testing.T) {
		api.HTTPTransport = &mockRoundTripper{
			handler: func(req *http.Request) (*http.Response, error) {
				return jsonResponse(422, api.ErrorResponse{
					Errors: []string{"Name has already been taken"},
				}), nil
			},
		}
		defer func() { api.HTTPTransport = nil }()

		_, err := executeCommand("servers", "new",
			"--name", "my-server",
			"--provider", "hetzner",
			"--region", "1",
			"--size", "2",
			"--image", "3",
			"--keys", "4",
		)
		assert.NotNil(t, err)
	})
}

type mockRoundTripper struct {
	handler func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.handler(req)
}

func jsonResponse(statusCode int, body any) *http.Response {
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(bodyBytes)),
		Header:     make(http.Header),
	}
}

func captureOutput(fn func()) string {
	r, w, err := os.Pipe()
	if err != nil {
		panic("failed to create pipe: " + err.Error())
	}

	old := os.Stdout
	defer func() { os.Stdout = old }()

	os.Stdout = w
	fn()

	if err := w.Close(); err != nil {
		panic("failed to close pipe writer: " + err.Error())
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		panic("failed to read captured output: " + err.Error())
	}

	if err := r.Close(); err != nil {
		panic("failed to close pipe reader: " + err.Error())
	}

	return buf.String()
}
