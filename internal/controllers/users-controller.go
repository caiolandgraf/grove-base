package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/caiolandgraf/go-project-base/internal/dto"
	"github.com/caiolandgraf/go-project-base/internal/services"
	"github.com/go-fuego/fuego"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

func (ctrl *UserController) GetUser(
	c fuego.ContextNoBody,
) (*dto.UserResponse, error) {
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

func (ctrl *UserController) ListUsers(
	c fuego.ContextNoBody,
) (*dto.UsersListResponse, error) {
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

func (ctrl *UserController) CreateUser(
	c fuego.ContextWithBody[dto.CreateUserRequest],
) (*dto.UserResponse, error) {
	body, err := c.Body()
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	user, err := ctrl.service.CreateUser(&body)
	if err != nil {
		if errors.Is(err, services.ErrorEmailAlreadyExists) {
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

func (ctrl *UserController) UpdateUser(
	c fuego.ContextWithBody[dto.UpdateUserRequest],
) (*dto.UserResponse, error) {
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

func (ctrl *UserController) DeleteUser(
	c fuego.ContextNoBody,
) (map[string]string, error) {
	userID := c.PathParam("user_id")

	err := ctrl.service.DeleteUser(userID)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusNotFound,
			Err:    err,
		}
	}

	return map[string]string{"message": "User deleted successfully"}, nil
}
