package auth

import "github.com/caiolandgraf/grove-base/internal/modules/users"

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email" example:"jane@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"secret123"`
}

type RegisterRequest struct {
	Name     string `json:"name"     validate:"required,min=3" example:"Jane Doe"`
	Email    string `json:"email"    validate:"required,email" example:"jane@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"secret123"`
}

type LoginResponse struct {
	User    users.UserResponse `json:"user"`
	Message string             `json:"message" example:"Login successful"`
}

type LogoutResponse struct {
	Message string `json:"message" example:"Logout successful"`
}

type SessionResponse struct {
	Authenticated bool                `json:"authenticated"  example:"true"`
	User          *users.UserResponse `json:"user,omitempty"`
}
