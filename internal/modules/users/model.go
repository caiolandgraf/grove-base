package users

import (
	"errors"
	"time"

	"github.com/caiolandgraf/grove-base/internal/app/database"
	"gorm.io/gorm"
)

func init() {
	database.Register(&User{})
}

var ErrNotFound = errors.New("user not found")

type User struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null"                     json:"name"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null"         json:"email"`
	Password  string         `gorm:"type:varchar(255);not null"                     json:"-"`
	CreatedAt time.Time      `gorm:"autoCreateTime"                                 json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"                                 json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                          json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

// Repo wraps the generic repository for User-specific queries.
type Repo struct {
	*database.Repository[User]
}

func Users(db *gorm.DB) *Repo {
	return &Repo{
		Repository: database.New[User](db),
	}
}

func (r *Repo) FindByEmail(email string) (*User, error) {
	user, err := r.FirstBy("email", email)
	if err != nil {
		return nil, ErrNotFound
	}
	return user, nil
}

func (r *Repo) FindByID(id string) (*User, error) {
	user, err := r.Find(id)
	if err != nil {
		return nil, ErrNotFound
	}
	return user, nil
}

func (r *Repo) FindAll(page, perPage int) ([]User, int64, error) {
	return r.Paginate(page, perPage)
}
