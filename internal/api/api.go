package api

import "fmt"

// APIError represents an error returned by the API
type APIError struct {
	StatusCode int
	Errors     []string
}

func (e *APIError) Error() string {
	switch e.StatusCode {
	case 401:
		return "Authentication failed. Please run 'bh auth' to configure your API key."
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
