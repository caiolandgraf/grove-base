package controllers

import (
	"errors"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/caiolandgraf/grove-base/internal/dto"
	"github.com/caiolandgraf/grove-base/internal/models"
	"github.com/go-fuego/fuego"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// AuthController holds only the session manager since DB comes from app.DB.
type AuthController struct {
	session *scs.SessionManager
}

func NewAuthController(session *scs.SessionManager) *AuthController {
	return &AuthController{session: session}
}

// Login authenticates the user and creates a session.
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

	user, err := models.Users().FindUserByEmail(body.Email)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusUnauthorized,
			Err:    ErrInvalidCredentials,
		}
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(body.Password),
	); err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusUnauthorized,
			Err:    ErrInvalidCredentials,
		}
	}

	// Create session
	ctx := c.Context()
	ctrl.session.Put(ctx, "user_id", user.ID)
	ctrl.session.Put(ctx, "user_email", user.Email)
	ctrl.session.Put(ctx, "user_name", user.Name)

	return &dto.LoginResponse{
		User:    *toUserDTO(user),
		Message: "Login successful",
	}, nil
}

// Register creates a new user and logs in automatically.
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

	repo := models.Users()

	exists, _ := repo.Exists("email = ?", body.Email)
	if exists {
		return nil, fuego.HTTPError{
			Status: http.StatusConflict,
			Err:    ErrorEmailAlreadyExists,
			Title:  "email already exists",
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(body.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}

	user := &models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(hashedPassword),
	}

	if err := repo.Create(user); err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	// Auto-login after registration
	ctx := c.Context()
	ctrl.session.Put(ctx, "user_id", user.ID)
	ctrl.session.Put(ctx, "user_email", user.Email)
	ctrl.session.Put(ctx, "user_name", user.Name)

	return &dto.LoginResponse{
		User:    *toUserDTO(user),
		Message: "Registration successful",
	}, nil
}

// Logout destroys the session.
func (ctrl *AuthController) Logout(
	c fuego.ContextNoBody,
) (*dto.LogoutResponse, error) {
	if err := ctrl.session.Destroy(c.Context()); err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}

	return &dto.LogoutResponse{
		Message: "Logout successful",
	}, nil
}

// Me returns the logged-in user's data from the session.
func (ctrl *AuthController) Me(
	c fuego.ContextNoBody,
) (*dto.SessionResponse, error) {
	ctx := c.Context()
	userID := ctrl.session.GetString(ctx, "user_id")

	if userID == "" {
		return &dto.SessionResponse{
			Authenticated: false,
		}, nil
	}

	user, err := models.Users().Find(userID)
	if err != nil {
		return &dto.SessionResponse{
			Authenticated: false,
		}, nil
	}

	return &dto.SessionResponse{
		Authenticated: true,
		User:          toUserDTO(user),
	}, nil
}
