package handler

import "errors"

// containsError checks if an error contains a specific error
func containsError(err, target error) bool {
	return errors.Is(err, target)
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

