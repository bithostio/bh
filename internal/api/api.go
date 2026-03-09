package api

import (
	"errors"
	"fmt"
)

// APIError represents an error returned by the API
type APIError struct {
	StatusCode int
	Errors     []string
}

func (e *APIError) Error() string {
	switch e.StatusCode {
	case 401:
		return "Authentication failed. Check your API key."
	case 403:
		return "Permission denied. Check your API key permissions."
	case 404:
		return "Resource not found."
	case 422:
		if len(e.Errors) > 0 {
			return fmt.Sprintf("Validation error: %s", e.Errors[0])
		}
		return "Validation error. Check your input."
	case 429:
		return "Rate limit exceeded. Please try again later."
	default:
		if len(e.Errors) > 0 {
			return e.Errors[0]
		}
		return fmt.Sprintf("API error: status %d", e.StatusCode)
	}
}

// IsAuthError reports whether the error is an API authentication (401) error.
func IsAuthError(err error) bool {
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == 401
}
