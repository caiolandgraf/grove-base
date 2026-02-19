package controllers

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/caiolandgraf/go-project-base/internal/dto"
	"github.com/caiolandgraf/go-project-base/internal/services"
	"github.com/go-fuego/fuego"
)

type AuthController struct {
	authService    services.AuthService
	sessionManager *scs.SessionManager
}

func NewAuthController(
	authService services.AuthService,
	sessionManager *scs.SessionManager,
) *AuthController {
	return &AuthController{
		authService:    authService,
		sessionManager: sessionManager,
	}
}

// Login authenticates the user and creates a session
func (ctrl *AuthController) Login(
	c fuego.ContextWithBody[dto.LoginRequest],
) (*dto.LoginResponse, error) {
	body, err := c.Body()
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	// Authenticate user
	user, err := ctrl.authService.Login(body.Email, body.Password)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusUnauthorized,
			Err:    err,
		}
	}

	// Create session
	ctx := c.Context()
	ctrl.sessionManager.Put(ctx, "user_id", user.ID)
	ctrl.sessionManager.Put(ctx, "user_email", user.Email)
	ctrl.sessionManager.Put(ctx, "user_name", user.Name)

	return &dto.LoginResponse{
		User:    *user,
		Message: "Login successful",
	}, nil
}

// Register creates a new user and logs in automatically
func (ctrl *AuthController) Register(
	c fuego.ContextWithBody[dto.RegisterRequest],
) (*dto.LoginResponse, error) {
	body, err := c.Body()
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	// Register user
	user, err := ctrl.authService.Register(&body)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	// Create session automatically
	ctx := c.Context()
	ctrl.sessionManager.Put(ctx, "user_id", user.ID)
	ctrl.sessionManager.Put(ctx, "user_email", user.Email)
	ctrl.sessionManager.Put(ctx, "user_name", user.Name)

	return &dto.LoginResponse{
		User:    *user,
		Message: "Registration successful",
	}, nil
}

// Logout destroys the session
func (ctrl *AuthController) Logout(
	c fuego.ContextNoBody,
) (*dto.LogoutResponse, error) {
	ctx := c.Context()

	// Destroy session
	err := ctrl.sessionManager.Destroy(ctx)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}

	return &dto.LogoutResponse{
		Message: "Logout successful",
	}, nil
}

// Me returns the logged-in user's data
func (ctrl *AuthController) Me(
	c fuego.ContextNoBody,
) (*dto.SessionResponse, error) {
	ctx := c.Context()
	userID := ctrl.sessionManager.GetString(ctx, "user_id")

	if userID == "" {
		return &dto.SessionResponse{
			Authenticated: false,
		}, nil
	}

	// Validate if user still exists
	user, err := ctrl.authService.ValidateUser(userID)
	if err != nil {
		return &dto.SessionResponse{
			Authenticated: false,
		}, nil
	}

	return &dto.SessionResponse{
		Authenticated: true,
		User:          user,
	}, nil
}
