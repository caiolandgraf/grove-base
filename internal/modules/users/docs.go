package users

import (
	"github.com/caiolandgraf/grove-base/internal/app/router"
	"github.com/caiolandgraf/grove-base/internal/app/types"
)

var ListUsersDoc = router.Doc{
	Detail: router.Detail{
		Tags:        []string{"Users"},
		Summary:     "List users",
		Description: "Returns a paginated list of users",
		OperationID: "listUsers",
	},
	Query: []router.QueryParam{
		{Name: "page", Description: "Page number", Type: router.QueryInt, Default: 1},
		{Name: "size", Description: "Items per page", Type: router.QueryInt, Default: 10},
	},
	Responses: map[int]router.Response{
		200: {Description: "Paginated user list", Type: UsersListResponse{}},
		429: {Description: "Rate limit exceeded", Type: types.ErrorResponse{}},
		500: {Description: "Internal server error", Type: types.ErrorResponse{}},
	},
}

var CreateUserDoc = router.Doc{
	Detail: router.Detail{
		Tags:        []string{"Users"},
		Summary:     "Create a new user",
		Description: "Creates a new user account",
		OperationID: "createUser",
	},
	Responses: map[int]router.Response{
		200: {Description: "User created", Type: UserResponse{}},
		400: {Description: "Validation error", Type: types.ErrorResponse{}},
		409: {Description: "Email already exists", Type: types.ErrorResponse{}},
		429: {Description: "Rate limit exceeded", Type: types.ErrorResponse{}},
	},
}

var GetUserDoc = router.Doc{
	Detail: router.Detail{
		Tags:        []string{"Users"},
		Summary:     "Get user by ID",
		Description: "Returns a single user by UUID",
		OperationID: "getUser",
	},
	Responses: map[int]router.Response{
		200: {Description: "User found", Type: UserResponse{}},
		404: {Description: "User not found", Type: types.ErrorResponse{}},
		429: {Description: "Rate limit exceeded", Type: types.ErrorResponse{}},
	},
}

var UpdateUserDoc = router.Doc{
	Detail: router.Detail{
		Tags:        []string{"Users"},
		Summary:     "Update user",
		Description: "Updates an existing user by UUID",
		OperationID: "updateUser",
	},
	Responses: map[int]router.Response{
		200: {Description: "User updated", Type: UserResponse{}},
		400: {Description: "Validation error", Type: types.ErrorResponse{}},
		404: {Description: "User not found", Type: types.ErrorResponse{}},
		429: {Description: "Rate limit exceeded", Type: types.ErrorResponse{}},
	},
}

var DeleteUserDoc = router.Doc{
	Detail: router.Detail{
		Tags:        []string{"Users"},
		Summary:     "Delete user",
		Description: "Soft-deletes a user by UUID",
		OperationID: "deleteUser",
	},
	Responses: map[int]router.Response{
		200: {Description: "User deleted", Type: types.MessageResponse{}},
		404: {Description: "User not found", Type: types.ErrorResponse{}},
		429: {Description: "Rate limit exceeded", Type: types.ErrorResponse{}},
	},
}
