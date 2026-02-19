package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // not serialized in JSON
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // soft delete
}

// TableName specifies the table name
func (User) TableName() string {
	return "users"
}

// BeforeCreate GORM hook - executed before creating
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// Here you can add logic before creating
	// For example, password hashing
	return nil
}

// BeforeUpdate GORM hook - executed before updating
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// Logic before updating
	return nil
}
