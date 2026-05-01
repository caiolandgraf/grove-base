package models

import (
	"time"

	"github.com/caiolandgraf/grove-base/internal/app"
	"github.com/caiolandgraf/grove-base/internal/app/database"
	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null"                     json:"name"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null"         json:"email"`
	Password  string         `gorm:"type:varchar(255);not null"                     json:"-"`
	CreatedAt time.Time      `gorm:"autoCreateTime"                                 json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"                                 json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                          json:"-"`
}

func (User) TableName() string { return "users" }

// recurrentBillRepository is a wrapper around the generic Repository for the recurrentBill model,
// allowing for custom queries and methods specific to Users.
type usersRepository struct {
	*database.Repository[User]
}

// Users returns a repository scoped to the recurrentBill model.
func Users() *usersRepository {
	return &usersRepository{
		Repository: database.New[User](app.DB),
	}
}

// Custom queries
func (r *usersRepository) FindUserByEmail(email string) (*User, error) {
	return Users().FirstBy("email", email)
}
