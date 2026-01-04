package std

import (
	"context"
	"errors"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/internal/types"
	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/gormen/pagination/sort"
	"gorm.io/gorm"
)

// FindAllPaginated retrieves a paginated list of records of type M without filters,
// supporting preloading related associations.
func (repository *repository[M]) FindAllPaginated(ctx context.Context, pageable pagination.Pageable, preloads ...gormen.Preload) (*pagination.Page[M], error) {
	return repository.FindAllPaginatedBy(ctx, pageable, gormen.Where{}, preloads...)
}

// FindAllPaginatedBy retrieves a paginated list of records of type M filtered by the given Where clause,
// supports preloading related associations, and returns total count and current page entities.
func (repository *repository[M]) FindAllPaginatedBy(ctx context.Context, pageable pagination.Pageable, where gormen.Where, preloads ...gormen.Preload) (*pagination.Page[M], error) {
	total, err := repository.count(ctx, pageable, where, preloads...)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &pagination.Page[M]{Total: total}, nil
	}

	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	for _, join := range where.Joins() {
		query = query.Joins(join)
	}

	page, err := pageable.Paginate(query)
	if err != nil {
		return nil, err
	}

	for cond, op := range where.Conditions() {
		switch op {
		case types.Or:
			query = query.Or(cond.Get())
		default:
			query = query.Where(cond.Get())
		}
	}

	var entities []M
	results := page.Find(&entities)
	if err := results.Error; err != nil {
		return nil, err
	}

	return &pagination.Page[M]{Total: total, Elements: entities}, nil
}

// FindAll retrieves all records of type M with optional preloading, no filtering.
func (repository *repository[M]) FindAll(ctx context.Context, preloads ...gormen.Preload) ([]M, error) {
	return repository.FindAllBy(ctx, gormen.Where{}, preloads...)
}

// FindAllBy retrieves all records of type M filtered by the given Where clause,
// supports preloading related associations.
func (repository *repository[M]) FindAllBy(ctx context.Context, where gormen.Where, preloads ...gormen.Preload) ([]M, error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	for _, join := range where.Joins() {
		query = query.Joins(join)
	}

	for cond, op := range where.Conditions() {
		switch op {
		case types.Or:
			query = query.Or(cond.Get())
		default:
			query = query.Where(cond.Get())
		}
	}

	var entities []M
	results := query.Find(&entities)
	if err := results.Error; err != nil {
		return nil, err
	}

	return entities, nil
}

// FindAllOrdered retrieves all records of type M ordered by the given orders,
// supports preloading related associations.
func (repository *repository[M]) FindAllOrdered(ctx context.Context, orders []sort.Order, preloads ...gormen.Preload) ([]M, error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	for _, o := range orders {
		query = query.Order(o.Get())
	}

	var entities []M
	results := query.Find(&entities)
	if err := results.Error; err != nil {
		return nil, err
	}

	return entities, nil
}

// count calculates the total number of records available based on the Pageable filtering and preloads.
func (repository repository[M]) count(ctx context.Context, pageable pagination.Pageable, where gormen.Where, preloads ...gormen.Preload) (int64, error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	for _, join := range where.Joins() {
		query = query.Joins(join)
	}

	for cond, op := range where.Conditions() {
		if op == types.Or {
			query = query.Or(cond.Get())
		} else {
			query = query.Where(cond.Get())
		}
	}

	query = query.Model(*new(M))

	filteredQuery, err := pageable.Filter(query)
	if err != nil {
		return 0, err
	}

	var count int64
	results := filteredQuery.Count(&count)
	if err := results.Error; err != nil {
		return 0, err
	}

	return count, nil
}

// FindBy fetches the first record of type M matching the Where clause with preloads,
// returns an pointer value with the record if found, otherwise Nil.
func (repository *repository[M]) FindBy(ctx context.Context, where gormen.Where, preloads ...gormen.Preload) (*M, error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	for _, join := range where.Joins() {
		query = query.Joins(join)
	}

	for cond, op := range where.Conditions() {
		switch op {
		case types.Or:
			query = query.Or(cond.Get())
		default:
			query = query.Where(cond.Get())
		}
	}

	var entity M
	result := query.First(&entity)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &entity, nil
}

// Count returns the total number of records without filters.
func (repository repository[_]) Count(ctx context.Context) (int64, error) {
	return repository.CountBy(ctx, gormen.Where{})
}

// CountBy returns the number of records matching the given Where clause.
func (repository repository[M]) CountBy(ctx context.Context, where gormen.Where) (int64, error) {
	query := repository.db.WithContext(ctx)

	for _, join := range where.Joins() {
		query = query.Joins(join)
	}

	for cond, op := range where.Conditions() {
		switch op {
		case types.Or:
			query = query.Or(cond.Get())
		default:
			query = query.Where(cond.Get())
		}
	}

	var count int64
	results := query.Model(*new(M)).Count(&count)
	if err := results.Error; err != nil {
		return 0, err
	}

	return count, nil
}
