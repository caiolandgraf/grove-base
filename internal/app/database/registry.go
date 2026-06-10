package database

var registry []any

// Register adds a GORM model to the Atlas migration registry.
// Call from each module's init() in model.go.
func Register(model any) {
	registry = append(registry, model)
}

// All returns every registered GORM model for Atlas migrations.
func All() []any {
	return registry
}
