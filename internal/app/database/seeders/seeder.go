package seeders

import (
	"fmt"

	"gorm.io/gorm"
)

// Seeder define o contrato de qualquer seeder.
type Seeder interface {
	Name() string
	Seed(db *gorm.DB) error
}

// Run executa todos os seeders em ordem.
// Se algum falhar, interrompe o processo.
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
