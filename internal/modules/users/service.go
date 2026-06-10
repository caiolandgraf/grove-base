package users

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	CreateUser(req *CreateUserRequest) (*UserResponse, error)
	GetUserByID(id string) (*UserResponse, error)
	GetUsers(page, size int) (*UsersListResponse, error)
	UpdateUser(id string, req *UpdateUserRequest) (*UserResponse, error)
	DeleteUser(id string) error
}

type service struct {
	repo *Repo
}

func NewService(repo *Repo) Service {
	return &service{repo: repo}
}

func WireService(db *gorm.DB) Service {
	return NewService(Users(db))
}

var ErrEmailAlreadyExists = errors.New("email already exists")

func (s *service) CreateUser(req *CreateUserRequest) (*UserResponse, error) {
	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return modelToDTO(user), nil
}

func (s *service) GetUserByID(id string) (*UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return modelToDTO(user), nil
}

func (s *service) GetUsers(page, size int) (*UsersListResponse, error) {
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

	userDTOs := make([]UserResponse, len(users))
	for i, user := range users {
		userDTOs[i] = *modelToDTO(&user)
	}

	return &UsersListResponse{
		Users: userDTOs,
		Total: int(total),
	}, nil
}

func (s *service) UpdateUser(
	id string,
	req *UpdateUserRequest,
) (*UserResponse, error) {
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

	return modelToDTO(user), nil
}

func (s *service) DeleteUser(id string) error {
	return s.repo.Delete(id)
}

func modelToDTO(user *User) *UserResponse {
	return &UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
