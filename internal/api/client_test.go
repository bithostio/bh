package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gsamokovarov/assert"
)

func TestListServers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "GET", req.Method)
			assert.True(t, strings.HasSuffix(req.URL.Path, "servers"))
			assert.Equal(t, "Bearer test-api-key", req.Header.Get("Authorization"))

			return jsonResponse(200, ServersResponse{
				Servers: []Server{
					{ID: 1, Name: "server-1", Status: "running", IPAddress: "1.2.3.4"},
					{ID: 2, Name: "server-2", Status: "stopped", IPAddress: "5.6.7.8"},
				},
				Meta: &Meta{Pagination: Pagination{Page: 1, HasNextPage: false}},
			}), nil
		})

		resp, err := client.ListServers(0)
		assert.Nil(t, err)
		assert.Len(t, 2, resp.Servers)
		assert.Equal(t, "server-1", resp.Servers[0].Name)
	})

	t.Run("with pagination", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "page=2", req.URL.RawQuery)
			return jsonResponse(200, ServersResponse{Servers: []Server{}}), nil
		})

		_, err := client.ListServers(2)
		assert.Nil(t, err)
	})

	t.Run("empty response", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			return jsonResponse(200, ServersResponse{Servers: []Server{}}), nil
		})

		resp, err := client.ListServers(0)
		assert.Nil(t, err)
		assert.Len(t, 0, resp.Servers)
	})
}

func TestCreateServer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "POST", req.Method)

			body, _ := io.ReadAll(req.Body)
			var createReq CreateServerRequest
			err := json.Unmarshal(body, &createReq)
			assert.Nil(t, err)
			assert.Equal(t, "my-server", createReq.Name)

			return jsonResponse(200, ServerResponse{
				Server: Server{ID: 123, Name: "my-server", Status: "pending"},
			}), nil
		})

		server, err := client.CreateServer(&CreateServerRequest{
			Name:     "my-server",
			SizeID:   1,
			RegionID: 2,
			ImageID:  3,
			Provider: "hetzner",
			KeyIDs:   []int{1, 2},
			Terms:    true,
		})
		assert.Nil(t, err)
		assert.Equal(t, 123, server.ID)
	})

	t.Run("validation error", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			return jsonResponse(422, ErrorResponse{
				Errors: []ErrorDetail{{Message: "Name is required", Code: "blank"}},
			}), nil
		})

		_, err := client.CreateServer(&CreateServerRequest{})
		assert.NotNil(t, err)

		apiErr, ok := err.(*APIError)
		assert.True(t, ok)
		assert.Equal(t, 422, apiErr.StatusCode)
	})
}

func TestDeleteServer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, "DELETE", req.Method)
			assert.True(t, strings.HasSuffix(req.URL.Path, "servers/123"))
			return jsonResponse(200, ServerResponse{Server: Server{ID: 123, Status: "deleted"}}), nil
		})

		err := client.DeleteServer(123)
		assert.Nil(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			return jsonResponse(404, ErrorResponse{
				Errors: []ErrorDetail{{Message: "Server not found", Code: "not_found"}},
			}), nil
		})

		err := client.DeleteServer(999)
		assert.NotNil(t, err)

		apiErr, ok := err.(*APIError)
		assert.True(t, ok)
		assert.Equal(t, 404, apiErr.StatusCode)
	})
}

func TestListProviders(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			return jsonResponse(200, ProvidersResponse{
				Providers: []Provider{
					{ID: 1, Name: "Hetzner", Slug: "hetzner"},
					{ID: 2, Name: "DigitalOcean", Slug: "digitalocean"},
				},
			}), nil
		})

		resp, err := client.ListProviders(0)
		assert.Nil(t, err)
		assert.Len(t, 2, resp.Providers)
	})
}

func TestListRegions(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			assert.True(t, strings.Contains(req.URL.RawQuery, "provider=hetzner"))
			return jsonResponse(200, RegionsResponse{
				Regions: []Region{
					{ID: 1, Name: "Frankfurt", Slug: "fsn1"},
					{ID: 2, Name: "Helsinki", Slug: "hel1"},
				},
			}), nil
		})

		resp, err := client.ListRegions("hetzner", 0)
		assert.Nil(t, err)
		assert.Len(t, 2, resp.Regions)
	})
}

func TestListSizes(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			assert.True(t, strings.Contains(req.URL.RawQuery, "region_id=1"))
			assert.True(t, strings.Contains(req.URL.RawQuery, "provider=hetzner"))
			return jsonResponse(200, SizesResponse{
				Sizes: []Size{
					{ID: 1, Name: "CX11", Slug: "cx11", Price: 4.99},
					{ID: 2, Name: "CX21", Slug: "cx21", Price: 9.99},
				},
			}), nil
		})

		resp, err := client.ListSizes(1, "hetzner", 0)
		assert.Nil(t, err)
		assert.Len(t, 2, resp.Sizes)
	})
}

func TestListImages(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			assert.True(t, strings.Contains(req.URL.RawQuery, "provider=hetzner"))
			assert.True(t, strings.Contains(req.URL.RawQuery, "architecture=x86_64"))
			return jsonResponse(200, ImagesResponse{
				Images: []Image{
					{ID: 1, Name: "Ubuntu 22.04", Distribution: "ubuntu"},
					{ID: 2, Name: "Debian 12", Distribution: "debian"},
				},
			}), nil
		})

		resp, err := client.ListImages("hetzner", "x86_64", 0)
		assert.Nil(t, err)
		assert.Len(t, 2, resp.Images)
	})
}

func TestListSSHKeys(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			return jsonResponse(200, KeysResponse{
				Keys: []SSHKey{
					{ID: 1, Label: "my-key", Key: "ssh-rsa AAAA..."},
				},
			}), nil
		})

		resp, err := client.ListSSHKeys(0)
		assert.Nil(t, err)
		assert.Len(t, 1, resp.Keys)
	})
}

func TestCreateSSHKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			body, _ := io.ReadAll(req.Body)
			var wrapped map[string]any
			err := json.Unmarshal(body, &wrapped)
			assert.Nil(t, err)

			keyData, ok := wrapped["key"].(map[string]any)
			assert.True(t, ok)
			assert.Equal(t, "my-key", keyData["label"])

			return jsonResponse(200, KeyResponse{
				Key: SSHKey{ID: 1, Label: "my-key", Key: "ssh-rsa AAAA..."},
			}), nil
		})

		key, err := client.CreateSSHKey(&CreateSSHKeyRequest{
			Label:    "my-key",
			Key:      "ssh-rsa AAAA...",
			Provider: "hetzner",
		})
		assert.Nil(t, err)
		assert.Equal(t, 1, key.ID)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client := newMockClient(func(req *http.Request) (*http.Response, error) {
			assert.True(t, strings.HasSuffix(req.URL.Path, "user"))
			return jsonResponse(200, UserResponse{
				User: User{
					FullName:    "John Doe",
					Email:       "john@example.com",
					Balance:     50.00,
					ServerLimit: 10,
				},
			}), nil
		})

		user, err := client.GetUser()
		assert.Nil(t, err)
		assert.Equal(t, "John Doe", user.FullName)
		assert.Equal(t, 50.00, user.Balance)
	})
}

func TestAPIErrors(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		errors         []ErrorDetail
		expectedSubstr string
	}{
		{"unauthorized", 401, []ErrorDetail{{Message: "Invalid API key", Code: "unauthenticated"}}, "Authentication failed"},
		{"forbidden", 403, []ErrorDetail{{Message: "Access denied", Code: "forbidden"}}, "Permission denied"},
		{"not found", 404, []ErrorDetail{{Message: "Not found", Code: "not_found"}}, "Resource not found"},
		{"validation error", 422, []ErrorDetail{{Message: "Name is too short", Code: "blank"}}, "Name is too short"},
		{"rate limited", 429, []ErrorDetail{{Message: "Too many requests", Code: "rate_limited"}}, "Rate limit exceeded"},
		{"server error", 500, []ErrorDetail{{Message: "Internal server error", Code: "internal_error"}}, "Internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := newMockClient(func(req *http.Request) (*http.Response, error) {
				return jsonResponse(tt.statusCode, ErrorResponse{Errors: tt.errors}), nil
			})

			_, err := client.ListServers(0)
			assert.NotNil(t, err)

			apiErr, ok := err.(*APIError)
			assert.True(t, ok)
			assert.Equal(t, tt.statusCode, apiErr.StatusCode)
			assert.True(t, strings.Contains(apiErr.Error(), tt.expectedSubstr))
		})
	}
}

func TestRequestHeaders(t *testing.T) {
	client := newMockClient(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "Bearer test-api-key", req.Header.Get("Authorization"))
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		return jsonResponse(200, ServersResponse{Servers: []Server{}}), nil
	})

	_, _ = client.ListServers(0)
}

type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func newMockClient(fn func(req *http.Request) (*http.Response, error)) *Client {
	return &Client{
		BaseURL: "https://api.example.com/",
		APIKey:  "test-api-key",
		HTTPClient: &http.Client{
			Transport: &mockRoundTripper{roundTripFunc: fn},
		},
	}
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
