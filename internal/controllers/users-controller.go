package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/caiolandgraf/go-project-base/internal/dto"
	"github.com/caiolandgraf/go-project-base/internal/models"
	"github.com/go-fuego/fuego"
	"golang.org/x/crypto/bcrypt"
)

var ErrorEmailAlreadyExists = errors.New("email already exists")

func GetUser(c fuego.ContextNoBody) (*dto.UserResponse, error) {
	userID := c.PathParam("user_id")

	user, err := models.Users().Find(userID)

	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusNotFound,
			Err:    err,
		}
	}

	return toUserDTO(user), nil
}

func ListUsers(c fuego.ContextNoBody) (*dto.UsersListResponse, error) {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("size"))

	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 15
	}

	users, total, err := models.Users().Paginate(page, size)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}

	items := make([]dto.UserResponse, len(users))
	for i, u := range users {
		items[i] = *toUserDTO(&u)
	}

	return &dto.UsersListResponse{
		Users: items,
		Total: int(total),
	}, nil
}

func CreateUser(
	c fuego.ContextWithBody[dto.CreateUserRequest],
) (*dto.UserResponse, error) {
	if err := fuego.ValidateParams(c); err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
			Title:  "invalid request body",
		}
	}

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
			Title:  "user with this email already exists",
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

	return toUserDTO(user), nil
}

func UpdateUser(
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

	repo := models.Users()

	user, err := repo.Find(userID)
	if err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusNotFound,
			Err:    err,
		}
	}

	if body.Name != "" {
		user.Name = body.Name
	}
	if body.Email != "" {
		user.Email = body.Email
	}

	if err := repo.Update(user); err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusBadRequest,
			Err:    err,
		}
	}

	return toUserDTO(user), nil
}

func DeleteUser(c fuego.ContextNoBody) (map[string]string, error) {
	userID := c.PathParam("user_id")

	if err := models.Users().Delete(userID); err != nil {
		return nil, fuego.HTTPError{
			Status: http.StatusNotFound,
			Err:    err,
		}
	}

	return map[string]string{"message": "User deleted successfully"}, nil
}

func toUserDTO(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
