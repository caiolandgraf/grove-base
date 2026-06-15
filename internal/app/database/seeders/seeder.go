package seeders

import (
	"fmt"

	"gorm.io/gorm"
)

type Seeder interface {
	Name() string
	Seed(db *gorm.DB) error
}

func Run(db *gorm.DB) error {
	seeders := []Seeder{
		UsersSeeder{},
	}

	for _, seeder := range seeders {
		if err := seeder.Seed(db); err != nil {
			return fmt.Errorf("seeder %s failed: %w", seeder.Name(), err)
		}
	}

	return nil
}
