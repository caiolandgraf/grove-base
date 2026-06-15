package auth

import (
	"github.com/caiolandgraf/grove-base/internal/app/router"
	"github.com/caiolandgraf/grove-base/internal/app/types"
	"github.com/caiolandgraf/grove-base/internal/modules/users"
)

var LoginDoc = router.Doc{
	Detail: router.Detail{
		Tags:        []string{"Auth"},
		Summary:     "Login",
		Description: "Authenticates a user and creates a session",
		OperationID: "login",
	},
	Responses: map[int]router.Response{
		200: {Description: "Login successful", Type: LoginResponse{}},
		400: {Description: "Validation error", Type: types.ErrorResponse{}},
		401: {Description: "Invalid credentials", Type: types.ErrorResponse{}},
	},
}

var RegisterDoc = router.Doc{
	Detail: router.Detail{
		Tags:        []string{"Auth"},
		Summary:     "Register a new user",
		Description: "Creates a new user account via the auth endpoint",
		OperationID: "registerUser",
	},
	Responses: map[int]router.Response{
		200: {Description: "User created", Type: users.UserResponse{}},
		400: {Description: "Validation error", Type: types.ErrorResponse{}},
		409: {Description: "Email already exists", Type: types.ErrorResponse{}},
	},
}

var LogoutDoc = router.Doc{
	Detail: router.Detail{
		Tags:        []string{"Auth"},
		Summary:     "Logout",
		Description: "Destroys the current session",
		OperationID: "logout",
	},
	Responses: map[int]router.Response{
		200: {Description: "Logout successful", Type: LogoutResponse{}},
		500: {Description: "Internal server error", Type: types.ErrorResponse{}},
	},
}

var MeDoc = router.Doc{
	Auth: true,
	Detail: router.Detail{
		Tags:        []string{"Auth"},
		Summary:     "Current session",
		Description: "Returns the authenticated user session data",
		OperationID: "me",
	},
	Responses: map[int]router.Response{
		200: {Description: "Session data", Type: SessionResponse{}},
		401: {Description: "Unauthorized", Type: types.ErrorResponse{}},
	},
}
