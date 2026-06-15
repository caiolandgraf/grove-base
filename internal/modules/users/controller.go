package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/caiolandgraf/grove-base/internal/app/types"
	"github.com/caiolandgraf/grove-base/internal/app/router"
	"github.com/go-fuego/fuego"
	"gorm.io/gorm"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func Wire(db *gorm.DB) *Controller {
	return NewController(WireService(db))
}

func (ctrl *Controller) Mount(api *fuego.Server, session *scs.SessionManager) {
	group := fuego.Group(api, "/users")

	router.Get(group, "/", ctrl.ListUsers, ListUsersDoc, session)
	router.Post(group, "/", ctrl.CreateUser, CreateUserDoc, session)
	router.Get(group, "/{user_id}", ctrl.GetUser, GetUserDoc, session)
	router.Put(group, "/{user_id}", ctrl.UpdateUser, UpdateUserDoc, session)
	router.Delete(group, "/{user_id}", ctrl.DeleteUser, DeleteUserDoc, session)
}

func (ctrl *Controller) GetUser(
	c fuego.ContextNoBody,
) (*UserResponse, error) {
	userID := c.PathParam("user_id")

	user, err := ctrl.service.GetUserByID(userID)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusNotFound,
			Err:    err,
		}
	}

	return user, nil
}

func (ctrl *Controller) ListUsers(
	c fuego.ContextNoBody,
) (*UsersListResponse, error) {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}

	users, err := ctrl.service.GetUsers(page, size)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}

	return users, nil
}

func (ctrl *Controller) CreateUser(
	c fuego.ContextWithBody[CreateUserRequest],
) (*UserResponse, error) {
	body, err := c.Body()
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	user, err := ctrl.service.CreateUser(&body)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			return nil, fuego.HTTPError{
				Status: http.StatusConflict,
				Err:    err,
				Title:  "user with this email already exists",
			}
		}

		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	return user, nil
}

func (ctrl *Controller) UpdateUser(
	c fuego.ContextWithBody[UpdateUserRequest],
) (*UserResponse, error) {
	userID := c.PathParam("user_id")
	body, err := c.Body()
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	user, err := ctrl.service.UpdateUser(userID, &body)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusNotFound,
			Err:    err,
		}
	}

	return user, nil
}

func (ctrl *Controller) DeleteUser(
	c fuego.ContextNoBody,
) (*types.MessageResponse, error) {
	userID := c.PathParam("user_id")

	err := ctrl.service.DeleteUser(userID)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusNotFound,
			Err:    err,
		}
	}

	return &types.MessageResponse{Message: "User deleted successfully"}, nil
}
