package database

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	db *gorm.DB
}

func New[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

// ──────────────────────────────────────────────
// Chainable query builder methods
// ──────────────────────────────────────────────

// Where adds a condition to the query chain.
//
// It's smart about the operator:
//   - If arg contains "%" it uses LIKE/ILIKE automatically
//   - If arg is nil it uses IS NULL
//   - Otherwise it uses = ?
//
// You can also pass a raw SQL expression as pointer:
//
//	r.Where("user_id", userID)                    → WHERE user_id = ?
//	r.Where("name", "%foo%")                      → WHERE name ILIKE ?
//	r.Where("deleted_at", nil)                    → WHERE deleted_at IS NULL
//	r.Where("age > ?", 18)                        → WHERE age > 18  (raw, arg ignored as field)
func (r *Repository[T]) Where(pointer string, arg any) *Repository[T] {
	var query string

	switch {
	// raw expression: pointer already contains "?" — treat arg as the bind value
	case strings.Contains(pointer, "?"):
		return &Repository[T]{
			db: r.db.Where(pointer, arg),
		}

	// nil arg → IS NULL
	case arg == nil:
		query = fmt.Sprintf("%s IS NULL", pointer)
		return &Repository[T]{
			db: r.db.Where(query),
		}

	// string with "%" → ILIKE (postgres) or LIKE (mysql/sqlite)
	case isLikePattern(arg):
		query = fmt.Sprintf("%s ILIKE ?", pointer)

	// default → equality
	default:
		query = fmt.Sprintf("%s = ?", pointer)
	}

	return &Repository[T]{
		db: r.db.Where(query, arg),
	}
}

// WhereRaw adds a raw WHERE clause with explicit args.
//
//	r.WhereRaw("age > ? AND active = ?", 18, true)
func (r *Repository[T]) WhereRaw(query string, args ...any) *Repository[T] {
	return &Repository[T]{db: r.db.Where(query, args...)}
}

// WhereIn adds a WHERE field IN (...) condition.
//
//	r.WhereIn("status", []string{"active", "pending"})
func (r *Repository[T]) WhereIn(field string, values any) *Repository[T] {
	return &Repository[T]{
		db: r.db.Where(fmt.Sprintf("%s IN ?", field), values),
	}
}

// WhereNotIn adds a WHERE field NOT IN (...) condition.
func (r *Repository[T]) WhereNotIn(field string, values any) *Repository[T] {
	return &Repository[T]{
		db: r.db.Where(fmt.Sprintf("%s NOT IN ?", field), values),
	}
}

// WhereNull adds WHERE field IS NULL.
func (r *Repository[T]) WhereNull(field string) *Repository[T] {
	return &Repository[T]{
		db: r.db.Where(fmt.Sprintf("%s IS NULL", field)),
	}
}

// WhereNotNull adds WHERE field IS NOT NULL.
func (r *Repository[T]) WhereNotNull(field string) *Repository[T] {
	return &Repository[T]{
		db: r.db.Where(fmt.Sprintf("%s IS NOT NULL", field)),
	}
}

// WhereBetween adds WHERE field BETWEEN from AND to.
func (r *Repository[T]) WhereBetween(
	field string,
	from, to any,
) *Repository[T] {
	return &Repository[T]{
		db: r.db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", field), from, to),
	}
}

// OrWhere adds an OR WHERE condition.
func (r *Repository[T]) OrWhere(query string, args ...any) *Repository[T] {
	return &Repository[T]{db: r.db.Or(query, args...)}
}

// Order adds an ORDER BY clause. Accepts any raw SQL expression.
//
//	r.Order("created_at DESC")
//	r.Order("char_length(name) ASC, name ASC")
//
// Multiple calls are additive:
//
//	r.Order("status ASC").Order("created_at DESC")
func (r *Repository[T]) Order(value string) *Repository[T] {
	return &Repository[T]{db: r.db.Order(value)}
}

// OrderBy is a typed helper for simple column ordering.
//
//	r.OrderBy("created_at", "DESC")
func (r *Repository[T]) OrderBy(field, direction string) *Repository[T] {
	direction = strings.ToUpper(direction)
	if direction != "ASC" && direction != "DESC" {
		direction = "ASC"
	}
	return r.Order(fmt.Sprintf("%s %s", field, direction))
}

// Latest orders by created_at DESC.
func (r *Repository[T]) Latest() *Repository[T] {
	return r.Order("created_at DESC")
}

// Oldest orders by created_at ASC.
func (r *Repository[T]) Oldest() *Repository[T] {
	return r.Order("created_at ASC")
}

// Limit sets the maximum number of records to return.
func (r *Repository[T]) Limit(n int) *Repository[T] {
	return &Repository[T]{db: r.db.Limit(n)}
}

// Offset skips n records.
func (r *Repository[T]) Offset(n int) *Repository[T] {
	return &Repository[T]{db: r.db.Offset(n)}
}

// With eager-loads relations (Preload).
//
//	r.With("Profile", "Orders")
func (r *Repository[T]) With(relations ...string) *Repository[T] {
	db := r.db
	for _, rel := range relations {
		db = db.Preload(rel)
	}
	return &Repository[T]{db: db}
}

// WithArgs eager-loads a relation with a condition callback.
//
//	r.WithArgs("Orders", func(q *gorm.DB) *gorm.DB {
//	    return q.Where("status = ?", "paid")
//	})
func (r *Repository[T]) WithArgs(
	relation string,
	fn func(*gorm.DB) *gorm.DB,
) *Repository[T] {
	return &Repository[T]{db: r.db.Preload(relation, fn)}
}

// Select specifies which columns to retrieve.
//
//	r.Select("id", "email", "created_at")
func (r *Repository[T]) Select(fields ...string) *Repository[T] {
	return &Repository[T]{db: r.db.Select(fields)}
}

// Distinct selects distinct records.
func (r *Repository[T]) Distinct(fields ...string) *Repository[T] {
	if len(fields) == 0 {
		return &Repository[T]{db: r.db.Distinct()}
	}
	args := make([]any, len(fields))
	for i, f := range fields {
		args[i] = f
	}
	return &Repository[T]{db: r.db.Distinct(args...)}
}

// ──────────────────────────────────────────────
// Execution — single record
// ──────────────────────────────────────────────

// Get executes the query and returns the first matching record.
func (r *Repository[T]) Get() (*T, error) {
	var entity T
	if err := r.db.First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%T not found", entity)
		}
		return nil, err
	}
	return &entity, nil
}

// Find retrieves a record by primary key.
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

// FindOrFail is an alias for Find.
func (r *Repository[T]) FindOrFail(id any) (*T, error) {
	return r.Find(id)
}

// FindMany retrieves multiple records by a slice of primary keys.
func (r *Repository[T]) FindMany(ids any) ([]T, error) {
	var entities []T
	if err := r.db.Where("id IN ?", ids).Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// First returns the first record matching a raw condition.
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
// Execution — multiple records
// ──────────────────────────────────────────────

// GetAll executes the query and returns all matching records.
func (r *Repository[T]) GetAll() ([]T, error) {
	var entities []T
	if err := r.db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// All returns every record (use with care on large tables).
func (r *Repository[T]) All() ([]T, error) {
	return r.GetAll()
}

// Pluck retrieves a single column into a slice.
//
//	var emails []string
//	err := repo.WhereRaw("active = ?", true).Pluck("email", &emails)
func (r *Repository[T]) Pluck(field string, dest any) error {
	return r.db.Pluck(field, dest).Error
}

// Chunk processes records in batches to avoid loading everything into memory.
//
//	repo.Where("active", true).Chunk(100, func(batch []User) error {
//	    // process
//	    return nil
//	})
func (r *Repository[T]) Chunk(size int, fn func([]T) error) error {
	var offset int
	for {
		var batch []T
		if err := r.db.Limit(size).
			Offset(offset).
			Find(&batch).
			Error; err != nil {
			return err
		}
		if len(batch) == 0 {
			break
		}
		if err := fn(batch); err != nil {
			return err
		}
		if len(batch) < size {
			break
		}
		offset += size
	}
	return nil
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

	var total int64
	// count on a clean session so LIMIT/OFFSET don't interfere
	if err := r.db.Session(&gorm.Session{}).Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var entities []T
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
// Fill
// ──────────────────────────────────────────────

// Fill populates the fields of dst with values from src.
//
// Rules:
//   - src can be a value or a pointer to a struct
//   - src fields of type T (non-pointer): only filled if not zero-value
//   - src fields of type *T (pointer): nil = ignored, non-nil = fills (even if zero, e.g. "")
//   - matching by field name (case-insensitive) or tag `fill:"DstFieldName"`
//
// Example with tag:
//
//	type UpdateBody struct {
//	    CategoryID *string `fill:"BillCategoryID"`
//	}
//
// Usage:
//
//	item, _ := Bills().Find(id)
//	Bills().Fill(item, body)
func (r *Repository[T]) Fill(dst *T, src any) *Repository[T] {
	if dst == nil {
		return r
	}

	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Ptr {
		if srcVal.IsNil() {
			return r
		}
		srcVal = srcVal.Elem()
	}
	if srcVal.Kind() != reflect.Struct {
		return r
	}

	srcType := srcVal.Type()
	srcIndex := make(map[string]reflect.Value, srcType.NumField())
	for i := range srcType.NumField() {
		f := srcType.Field(i)
		srcIndex[strings.ToLower(f.Name)] = srcVal.Field(i)
	}

	dstVal := reflect.ValueOf(dst).Elem()
	dstType := dstVal.Type()

	for i := range dstType.NumField() {
		dstField := dstType.Field(i)
		dstFieldVal := dstVal.Field(i)

		if !dstFieldVal.CanSet() {
			continue
		}

		lookupName := dstField.Name
		if tag, ok := dstField.Tag.Lookup(
			"fill",
		); ok && tag != "" &&
			tag != "-" {
			lookupName = tag
		}

		srcFieldVal, found := srcIndex[strings.ToLower(lookupName)]
		if !found {
			continue
		}

		switch {
		case srcFieldVal.Type() == dstFieldVal.Type():
			if !srcFieldVal.IsZero() {
				dstFieldVal.Set(srcFieldVal)
			}

		case srcFieldVal.Kind() == reflect.Ptr && dstFieldVal.Kind() != reflect.Ptr:
			if !srcFieldVal.IsNil() &&
				srcFieldVal.Elem().Type() == dstFieldVal.Type() {
				dstFieldVal.Set(srcFieldVal.Elem())
			}

		case srcFieldVal.Kind() == reflect.Ptr && dstFieldVal.Kind() == reflect.Ptr:
			if !srcFieldVal.IsNil() &&
				srcFieldVal.Type() == dstFieldVal.Type() {
				dstFieldVal.Set(srcFieldVal)
			}
		}
	}

	return r
}

// ──────────────────────────────────────────────
// Write operations
// ──────────────────────────────────────────────

// Create persists a new record.
func (r *Repository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// CreateMany inserts multiple records in a single statement.
func (r *Repository[T]) CreateMany(entities []T) error {
	if len(entities) == 0 {
		return nil
	}
	return r.db.Create(&entities).Error
}

// CreateWithTx persists a new record using an existing transaction.
func (r *Repository[T]) CreateWithTx(tx *gorm.DB, entity *T) error {
	if entity == nil {
		var zero T
		return fmt.Errorf("cannot create: %T entity is nil", zero)
	}
	db := r.db
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

// DeleteWhere deletes all records matching a condition.
func (r *Repository[T]) DeleteWhere(query string, args ...any) error {
	return r.db.Where(query, args...).Delete(new(T)).Error
}

// ForceDelete permanently removes a record even with soft-delete enabled.
func (r *Repository[T]) ForceDelete(id any) error {
	return r.db.Unscoped().Delete(new(T), "id = ?", id).Error
}

// ──────────────────────────────────────────────
// Existence / counting helpers
// ──────────────────────────────────────────────

// Exists returns true if at least one record matches the current query chain.
//
//	repo.Where("email", "test@example.com").Exists()
func (r *Repository[T]) Exists() (bool, error) {
	var count int64
	// Use a clean session to avoid side effects from LIMIT/OFFSET on the count
	// and apply the existing query chain.
	if err := r.db.Session(&gorm.Session{}).
		Model(new(T)).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// Count returns the number of records matching the current query chain.
func (r *Repository[T]) Count() (int64, error) {
	var count int64
	if err := r.db.Model(new(T)).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ──────────────────────────────────────────────
// Aggregate functions
// ──────────────────────────────────────────────

// Sum calculates the sum of a column and scans the result into dest.
//
//	var total float64
//	err := repo.Where("user_id", userID).Sum("amount", &total)
func (r *Repository[T]) Sum(column string, dest any) error {
	return r.db.Model(new(T)).
		Select("COALESCE(SUM(" + column + "), 0)").
		Scan(dest).Error
}

// Avg calculates the average of a column and scans the result into dest.
//
//	var avg float64
//	err := repo.Where("status", "active").Avg("score", &avg)
func (r *Repository[T]) Avg(column string, dest any) error {
	return r.db.Model(new(T)).
		Select("COALESCE(AVG(" + column + "), 0)").
		Scan(dest).Error
}

// Min finds the minimum value of a column and scans the result into dest.
//
//	var earliest time.Time
//	err := repo.Where("user_id", userID).Min("created_at", &earliest)
func (r *Repository[T]) Min(column string, dest any) error {
	return r.db.Model(new(T)).
		Select("MIN(" + column + ")").
		Scan(dest).Error
}

// Max finds the maximum value of a column and scans the result into dest.
//
//	var highest float64
//	err := repo.Where("user_id", userID).Max("amount", &highest)
func (r *Repository[T]) Max(column string, dest any) error {
	return r.db.Model(new(T)).
		Select("MAX(" + column + ")").
		Scan(dest).Error
}

// GroupBy adds a GROUP BY clause to the query chain.
// Typically combined with Select + Scan for aggregated results.
//
//	type Result struct {
//	    UserID uuid.UUID
//	    Total  float64
//	}
//	var results []Result
//	err := repo.
//	    Select("user_id", "SUM(amount) as total").
//	    GroupBy("user_id").
//	    Scan(&results)
func (r *Repository[T]) GroupBy(columns ...string) *Repository[T] {
	return &Repository[T]{db: r.db.Group(strings.Join(columns, ", "))}
}

// Having adds a HAVING clause (use after GroupBy).
//
//	repo.GroupBy("user_id").Having("SUM(amount) > ?", 1000)
func (r *Repository[T]) Having(query string, args ...any) *Repository[T] {
	return &Repository[T]{db: r.db.Having(query, args...)}
}

// Scan executes the current query and scans the results into dest.
// Useful for custom SELECT projections or aggregate queries.
//
//	var results []MyDTO
//	err := repo.Select("user_id", "COUNT(*) as total").GroupBy("user_id").Scan(&results)
func (r *Repository[T]) Scan(dest any) error {
	return r.db.Scan(dest).Error
}

// SumMap returns a map of group → sum for a given group column and sum column.
// Convenience wrapper around GroupBy + Sum for simple breakdowns.
//
//	// total spend per status
//	totals, err := repo.SumMap("status", "amount")
//	// totals["paid"] = 9800.00, totals["pending"] = 250.00
func (r *Repository[T]) SumMap(
	groupColumn, sumColumn string,
) (map[string]float64, error) {
	type row struct {
		Group string  `gorm:"column:grp"`
		Total float64 `gorm:"column:total"`
	}
	var rows []row
	err := r.db.Model(new(T)).
		Select(groupColumn + " as grp, COALESCE(SUM(" + sumColumn + "), 0) as total").
		Group(groupColumn).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[string]float64, len(rows))
	for _, row := range rows {
		result[row.Group] = row.Total
	}
	return result, nil
}

// ──────────────────────────────────────────────
// Transactions
// ──────────────────────────────────────────────

// Transaction executes fn inside a DB transaction.
// Rolls back automatically on error, commits on success.
//
//	repo.Transaction(func(tx *gorm.DB) error {
//	    return database.New[models.Order](tx).Create(&order)
//	})
func (r *Repository[T]) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

// ──────────────────────────────────────────────
// Escape hatches
// ──────────────────────────────────────────────

// Query returns a *gorm.DB scoped to the model for arbitrary GORM calls.
func (r *Repository[T]) Query() *gorm.DB {
	return r.db.Model(new(T))
}

// DB returns the raw *gorm.DB without model scoping.
func (r *Repository[T]) DB() *gorm.DB {
	return r.db
}

// ──────────────────────────────────────────────
// Internal helpers
// ──────────────────────────────────────────────

// isLikePattern returns true if the value is a string containing "%".
func isLikePattern(arg any) bool {
	s, ok := arg.(string)
	return ok && strings.Contains(s, "%")
}
