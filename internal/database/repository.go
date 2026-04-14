package database

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Repository is a generic CRUD repository for any GORM model.
// It provides Laravel Eloquent-like methods so you never need
// to write a per-entity repository again.
//
// Usage:
//
//	repo := database.New[models.User](app.DB)
//	user, err := repo.Find("some-uuid")
//	users, total, err := repo.Paginate(1, 15)
type Repository[T any] struct {
	db *gorm.DB
}

// New creates a Repository for the given model type.
func New[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

// ──────────────────────────────────────────────
// Single record queries
// ──────────────────────────────────────────────

// Find retrieves a record by its primary key.
func (r *Repository[T]) Find(id any) (*T, error) {
	var entity T
	if err := r.db.First(&entity, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%T not found", entity)
		}
		return nil, err
	}
	return &entity, nil
}

// FindOrFail is an alias for Find (same behavior – returns error when not found).
func (r *Repository[T]) FindOrFail(id any) (*T, error) {
	return r.Find(id)
}

// First returns the first record that matches the conditions.
func (r *Repository[T]) First(query string, args ...any) (*T, error) {
	var entity T
	if err := r.db.Where(query, args...).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%T not found", entity)
		}
		return nil, err
	}
	return &entity, nil
}

// FirstBy is a shortcut: repo.FirstBy("email", "a@b.com")
func (r *Repository[T]) FirstBy(field string, value any) (*T, error) {
	return r.First(field+" = ?", value)
}

// ──────────────────────────────────────────────
// Multiple record queries
// ──────────────────────────────────────────────

// All returns every record (use with care).
func (r *Repository[T]) All() ([]T, error) {
	var entities []T
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Where returns records matching the condition.
func (r *Repository[T]) Where(query string, args ...any) ([]T, error) {
	var entities []T
	if err := r.db.Where(query, args...).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Paginate returns a page of results plus the total count.
// page is 1-based; perPage defaults to 15 when <= 0.
func (r *Repository[T]) Paginate(page, perPage int) ([]T, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 15
	}

	var (
		entities []T
		total    int64
	)

	if err := r.db.Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	if err := r.db.Limit(perPage).
		Offset(offset).
		Find(&entities).
		Error; err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// ──────────────────────────────────────────────
// Write operations
// ──────────────────────────────────────────────

// Create persists a new record.
func (r *Repository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *Repository[T]) CreateWithTx(
	tx *gorm.DB,
	entity *T,
) error {
	if entity == nil {
		var zero T
		return fmt.Errorf("cannot create: %T entity is nil", zero)
	}

	db := r.DB()
	if tx != nil {
		db = tx
	}

	return db.Create(entity).Error
}

// Update saves all fields of the record.
func (r *Repository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// UpdateFields updates only the specified fields by primary key.
//
//	repo.UpdateFields(id, map[string]any{"name": "New Name"})
func (r *Repository[T]) UpdateFields(id any, fields map[string]any) error {
	return r.db.Model(new(T)).Where("id = ?", id).Updates(fields).Error
}

// Delete soft-deletes (or hard-deletes) a record by primary key.
func (r *Repository[T]) Delete(id any) error {
	return r.db.Delete(new(T), "id = ?", id).Error
}

// ──────────────────────────────────────────────
// Existence / counting helpers
// ──────────────────────────────────────────────

// Exists returns true if at least one record matches.
func (r *Repository[T]) Exists(query string, args ...any) (bool, error) {
	var count int64
	if err := r.db.Model(new(T)).
		Where(query, args...).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count returns the number of records matching the condition.
func (r *Repository[T]) Count(query string, args ...any) (int64, error) {
	var count int64
	if err := r.db.Model(new(T)).
		Where(query, args...).
		Count(&count).
		Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ──────────────────────────────────────────────
// Escape hatch – access the underlying *gorm.DB
// ──────────────────────────────────────────────

// Query returns a *gorm.DB scoped to the model so you can
// chain arbitrary GORM calls when the helpers aren't enough.
//
//	repo.Query().Joins("Profile").Where("age > ?", 18).Find(&users)
func (r *Repository[T]) Query() *gorm.DB {
	return r.db.Model(new(T))
}

// DB returns the raw *gorm.DB without model scoping.
func (r *Repository[T]) DB() *gorm.DB {
	return r.db
}
