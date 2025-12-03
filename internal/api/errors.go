package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError represents an error returned by the API
type APIError struct {
	StatusCode int
	Errors     []string
}

func (e *APIError) Error() string {
	if len(e.Errors) == 0 {
		return fmt.Sprintf("API error: status %d", e.StatusCode)
	}
	return fmt.Sprintf("API error: %s", e.Errors[0])
}

// HandleError converts an error to a user-friendly message
func HandleError(err error) string {
	if apiErr, ok := err.(*APIError); ok {
		switch apiErr.StatusCode {
		case 401:
			return "Authentication failed. Please run 'bh auth' to configure your API key."
		case 403:
			return "Permission denied. Check your API key permissions."
		case 404:
			return "Resource not found."
		case 422:
			if len(apiErr.Errors) > 0 {
				return fmt.Sprintf("Validation error: %s", apiErr.Errors[0])
			}
			return "Validation error. Check your input."
		case 429:
			return "Rate limit exceeded. Please try again later."
		default:
			if len(apiErr.Errors) > 0 {
				return apiErr.Errors[0]
			}
			return fmt.Sprintf("API error: %s", err)
		}
	}
	return fmt.Sprintf("Error: %s", err)
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
