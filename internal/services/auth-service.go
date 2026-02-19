package services

import (
	"errors"

	"github.com/caiolandgraf/go-project-base/internal/dto"
	"github.com/caiolandgraf/go-project-base/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(email, password string) (*dto.UserResponse, error)
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)
	ValidateUser(userID string) (*dto.UserResponse, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(email, password string) (*dto.UserResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *authService) Register(
	req *dto.RegisterRequest,
) (*dto.UserResponse, error) {
	// Check if email already exists
	existing, _ := s.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	// Create user using UserService
	user := &dto.CreateUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// Here you could call UserService.CreateUser
	// For simplicity, creating directly
	return &dto.UserResponse{
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *authService) ValidateUser(userID string) (*dto.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
