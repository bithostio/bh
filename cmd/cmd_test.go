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
