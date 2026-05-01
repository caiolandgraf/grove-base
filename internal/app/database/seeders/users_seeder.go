package seeders

import (
	"fmt"

	"github.com/caiolandgraf/grove-base/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UsersSeeder popula usuários iniciais do sistema.
type UsersSeeder struct{}

func (UsersSeeder) Name() string { return "UsersSeeder" }

func (UsersSeeder) Seed(db *gorm.DB) error {
	defaultUsers := []struct {
		Name     string
		Email    string
		Password string
	}{
		{
			Name:     "Admin",
			Email:    "admin@grove.local",
			Password: "admin123",
		},
		{
			Name:     "Test User",
			Email:    "user@grove.local",
			Password: "user1234",
		},
	}

	for _, item := range defaultUsers {
		var existing models.User
		err := db.
			Where("email = ?", item.Email).
			First(&existing).
			Error

		switch err {
		case nil:
			continue

		case gorm.ErrRecordNotFound:
			hash, hashErr := bcrypt.GenerateFromPassword(
				[]byte(item.Password),
				bcrypt.DefaultCost,
			)
			if hashErr != nil {
				return fmt.Errorf(
					"hash password for %s: %w",
					item.Email,
					hashErr,
				)
			}

			newUser := models.User{
				Name:     item.Name,
				Email:    item.Email,
				Password: string(hash),
			}

			if createErr := db.Create(&newUser).Error; createErr != nil {
				return fmt.Errorf("create user %s: %w", item.Email, createErr)
			}

		default:
			return fmt.Errorf("query user %s: %w", item.Email, err)
		}
	}

	return nil
}
