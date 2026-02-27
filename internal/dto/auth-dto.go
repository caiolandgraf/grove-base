package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type SessionResponse struct {
	Authenticated bool          `json:"authenticated"`
	User          *UserResponse `json:"user,omitempty"`
}
