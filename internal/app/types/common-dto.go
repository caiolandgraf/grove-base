// Package types defines shared API request and response DTOs.
package types

// ErrorResponse is the standard error payload returned by the API.
type ErrorResponse struct {
	Message string `json:"message" example:"validation failed"`
	Error   string `json:"error,omitempty" example:"email already exists"`
}

// MessageResponse is a simple success message payload.
type MessageResponse struct {
	Message string `json:"message" example:"User deleted successfully"`
}
