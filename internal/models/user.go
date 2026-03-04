package models

import (
	"time"

	"github.com/caiolandgraf/grove-base/internal/app"
	"github.com/caiolandgraf/grove-base/internal/database"
	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string { return "users" }

// Users returns a repository scoped to the User model.
func Users() *database.Repository[User] {
	return database.New[User](app.DB)
}

// Custom queries
func FindUserByEmail(email string) (*User, error) {
	return Users().FirstBy("email", email)
}
