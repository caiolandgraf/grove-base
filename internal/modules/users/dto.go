package users

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3" example:"Jane Doe"`
	Email    string `json:"email" validate:"required,email" example:"jane@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"secret123"`
}

type UpdateUserRequest struct {
	Name  string `json:"name,omitempty" validate:"omitempty,min=3" example:"Jane Doe"`
	Email string `json:"email,omitempty" validate:"omitempty,email" example:"jane@example.com"`
}

type UserResponse struct {
	ID    string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name  string `json:"name" example:"Jane Doe"`
	Email string `json:"email" example:"jane@example.com"`
}

type UsersListResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total" example:"42"`
}
