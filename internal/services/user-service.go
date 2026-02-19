package services

import (
	"errors"

	"github.com/caiolandgraf/go-project-base/internal/dto"
	"github.com/caiolandgraf/go-project-base/internal/models"
	"github.com/caiolandgraf/go-project-base/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
	GetUserByID(id string) (*dto.UserResponse, error)
	GetUsers(page, size int) (*dto.UsersListResponse, error)
	UpdateUser(id string, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

var ErrorEmailAlreadyExists = errors.New("email already exists")

func (s *userService) CreateUser(
	req *dto.CreateUserRequest,
) (*dto.UserResponse, error) {
	// Check if email already exists
	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, ErrorEmailAlreadyExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return s.modelToDTO(user), nil
}

func (s *userService) GetUserByID(id string) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.modelToDTO(user), nil
}

func (s *userService) GetUsers(page, size int) (*dto.UsersListResponse, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	offset := (page - 1) * size
	users, total, err := s.repo.FindAll(size, offset)
	if err != nil {
		return nil, err
	}

	userDTOs := make([]dto.UserResponse, len(users))
	for i, user := range users {
		userDTOs[i] = *s.modelToDTO(&user)
	}

	return &dto.UsersListResponse{
		Users: userDTOs,
		Total: int(total),
	}, nil
}

func (s *userService) UpdateUser(
	id string,
	req *dto.UpdateUserRequest,
) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.modelToDTO(user), nil
}

func (s *userService) DeleteUser(id string) error {
	return s.repo.Delete(id)
}

// Helper to convert Model -> DTO
func (s *userService) modelToDTO(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
