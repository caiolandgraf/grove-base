package auth

import (
	"errors"

	"github.com/caiolandgraf/go-project-base/internal/modules/users"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Login(email, password string) (*users.UserResponse, error)
	Register(req *RegisterRequest) (*users.UserResponse, error)
	ValidateUser(userID string) (*users.UserResponse, error)
}

type service struct {
	userRepo *users.Repo
}

func NewService(userRepo *users.Repo) Service {
	return &service{userRepo: userRepo}
}

func WireService(db *gorm.DB) Service {
	return NewService(users.Users(db))
}

func (s *service) Login(email, password string) (*users.UserResponse, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &users.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *service) Register(req *RegisterRequest) (*users.UserResponse, error) {
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	return &users.UserResponse{
		Name:  req.Name,
		Email: req.Email,
	}, nil
}

func (s *service) ValidateUser(userID string) (*users.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &users.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
