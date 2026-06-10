package auth

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/caiolandgraf/go-project-base/internal/modules/users"
	"github.com/caiolandgraf/go-project-base/internal/app/router"
	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type Controller struct {
	authService    Service
	userService    users.Service
	sessionManager *scs.SessionManager
}

func NewController(
	authService Service,
	userService users.Service,
	session *scs.SessionManager,
) *Controller {
	return &Controller{
		authService:    authService,
		userService:    userService,
		sessionManager: session,
	}
}

func Wire(db *gorm.DB, session *scs.SessionManager) *Controller {
	repo := users.Users(db)
	return NewController(
		NewService(repo),
		users.NewService(repo),
		session,
	)
}

func (ctrl *Controller) Mount(api *fuego.Server, session *scs.SessionManager) {
	group := fuego.Group(api, "/auth")
	userHandlers := users.NewController(ctrl.userService)

	router.Post(group, "/login", ctrl.Login, LoginDoc, session)
	router.Post(group, "/register", userHandlers.CreateUser, RegisterDoc, session)
	router.Post(group, "/logout", ctrl.Logout, LogoutDoc, session)
	router.Get(group, "/me", ctrl.Me, MeDoc, session)
}

func (ctrl *Controller) Login(
	c fuego.ContextWithBody[LoginRequest],
) (*LoginResponse, error) {
	body, err := c.Body()
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	user, err := ctrl.authService.Login(body.Email, body.Password)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusUnauthorized,
			Err:    err,
		}
	}

	ctx := c.Context()
	ctrl.sessionManager.Put(ctx, "user_id", user.ID)
	ctrl.sessionManager.Put(ctx, "user_email", user.Email)
	ctrl.sessionManager.Put(ctx, "user_name", user.Name)

	return &LoginResponse{
		User:    *user,
		Message: "Login successful",
	}, nil
}

func (ctrl *Controller) Logout(
	c fuego.ContextNoBody,
) (*LogoutResponse, error) {
	ctx := c.Context()

	err := ctrl.sessionManager.Destroy(ctx)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}

	return &LogoutResponse{
		Message: "Logout successful",
	}, nil
}

func (ctrl *Controller) Me(
	c fuego.ContextNoBody,
) (*SessionResponse, error) {
	ctx := c.Context()
	userID := ctrl.sessionManager.GetString(ctx, "user_id")

	if userID == "" {
		return &SessionResponse{
			Authenticated: false,
		}, nil
	}

	user, err := ctrl.authService.ValidateUser(userID)
	if err != nil {
		return &SessionResponse{
			Authenticated: false,
		}, nil
	}

	return &SessionResponse{
		Authenticated: true,
		User:          user,
	}, nil
}
