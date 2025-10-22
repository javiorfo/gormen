package converter

import (
	"context"
	"errors"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/internal/types"
	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/gormen/pagination/sort"
	"github.com/javiorfo/nilo"
	"github.com/javiorfo/steams"
	"gorm.io/gorm"
)

// FindAllPaginated returns a paginated list of models M based on the given Pageable and optional preloads.
// It delegates to FindAllPaginatedBy with an empty Where condition.
func (repository *repository[E, C, M]) FindAllPaginated(ctx context.Context, pageable pagination.Pageable, preloads ...gormen.Preload) (*pagination.Page[M], error) {
	return repository.FindAllPaginatedBy(ctx, pageable, gormen.Where{}, preloads...)
}

// FindAllPaginatedBy fetches a paginated list of models M filtered by the specified Where conditions,
// supports eager loading of associations via preloads, and returns total count and models.
func (repository *repository[E, C, M]) FindAllPaginatedBy(ctx context.Context, pageable pagination.Pageable, where gormen.Where, preloads ...gormen.Preload) (*pagination.Page[M], error) {
	total, err := repository.count(ctx, pageable, preloads...)
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

	var entities []E
	page, err := pageable.Paginate(query)
	if err != nil {
		return nil, err
	}

	for cond, op := range where.Conditions() {
		if op == types.Or {
			query = query.Or(cond.Get())
		} else {
			query = query.Where(cond.Get())
		}
	}

	results := page.Find(&entities)
	if err := results.Error; err != nil {
		return nil, err
	}

	models := steams.Mapper(steams.OfSlice(entities), func(entity E) M {
		var c C = &entity
		return c.Into()
	}).Collect()

	return &pagination.Page[M]{Total: total, Elements: models}, nil
}

// FindAll gets all records of model M, applying optional preloads.
func (repository *repository[E, C, M]) FindAll(ctx context.Context, preloads ...gormen.Preload) ([]M, error) {
	return repository.FindAllBy(ctx, gormen.Where{}, preloads...)
}

// FindAllBy gets all records of model M filtered by Where conditions and applies preloads.
func (repository *repository[E, C, M]) FindAllBy(ctx context.Context, where gormen.Where, preloads ...gormen.Preload) ([]M, error) {
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

	var entities []E
	results := query.Find(&entities)
	if err := results.Error; err != nil {
		return nil, err
	}

	models := steams.Mapper(steams.OfSlice(entities), func(entity E) M {
		var c C = &entity
		return c.Into()
	}).Collect()

	return models, nil
}

// FindAllOrdered retrieves all records of model M ordered by specified sort orders and applies preloads.
func (repository *repository[E, C, M]) FindAllOrdered(ctx context.Context, orders []sort.Order, preloads ...gormen.Preload) ([]M, error) {
	query := repository.db.WithContext(ctx)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	for _, o := range orders {
		query = query.Order(o.Get())
	}

	var entities []E
	results := query.Find(&entities)
	if err := results.Error; err != nil {
		return nil, err
	}

	models := steams.Mapper(steams.OfSlice(entities), func(entity E) M {
		var c C = &entity
		return c.Into()
	}).Collect()

	return models, nil
}

// count returns the total number of records matching the Pageable's filter criteria and preloads.
func (repository repository[E, _, _]) count(ctx context.Context, pageable pagination.Pageable, preloads ...gormen.Preload) (int64, error) {
	query := repository.db.WithContext(ctx)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	query = query.Model(*new(E))

	filteredQuery, err := pageable.Filter(query)
	if err != nil {
		return 0, err
	}

	var count int64
	result := filteredQuery.Count(&count)
	if err := result.Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindBy retrieves the first record matching the Where conditions with preloads,
// returns an Option of model M — None if not found.
func (repository *repository[E, C, M]) FindBy(ctx context.Context, where gormen.Where, preloads ...gormen.Preload) (nilo.Option[M], error) {
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

	var entity C = new(E)
	result := query.First(&entity)
	if err := result.Error; err != nil {
		none := nilo.None[M]()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return none, nil
		}
		return none, err
	}

	model := entity.Into()
	return nilo.Some(model), nil
}

// Count returns the total number of records without filters.
func (repository repository[E, _, _]) Count(ctx context.Context) (int64, error) {
	return repository.CountBy(ctx, gormen.Where{})
}

// CountBy returns count of records matching the Where conditions.
func (repository repository[E, _, _]) CountBy(ctx context.Context, where gormen.Where) (int64, error) {
	query := repository.db.WithContext(ctx)
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

	var count int64
	results := query.Model(*new(E)).Count(&count)
	if err := results.Error; err != nil {
		return 0, err
	}

	return count, nil
}
