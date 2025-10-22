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

func (repository *repository[E, C, M]) FindAllPaginated(ctx context.Context, pageable pagination.Pageable, preloads ...gormen.Preload) (*pagination.Page[M], error) {
	return repository.FindAllPaginatedBy(ctx, pageable, gormen.Where{}, preloads...)
}

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

	for k, v := range where.Conditions() {
		switch v {
		case types.Or:
			query = query.Or(k.Get())
		default:
			query = query.Where(k.Get())
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

func (repository *repository[E, C, M]) FindAll(ctx context.Context, preloads ...gormen.Preload) ([]M, error) {
	return repository.FindAllBy(ctx, gormen.Where{}, preloads...)
}

func (repository *repository[E, C, M]) FindAllBy(ctx context.Context, where gormen.Where, preloads ...gormen.Preload) ([]M, error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	for _, join := range where.Joins() {
		query = query.Joins(join)
	}

	for k, v := range where.Conditions() {
		switch v {
		case types.Or:
			query = query.Or(k.Get())
		default:
			query = query.Where(k.Get())
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

func (repository *repository[E, C, M]) FindAllOrdered(ctx context.Context, orders []sort.Order, preloads ...gormen.Preload) ([]M, error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	for _, o := range orders {
		query = query.Order(o.Prepared())
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

func (repository repository[E, _, _]) count(ctx context.Context, pageable pagination.Pageable, preloads ...gormen.Preload) (int64, error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	query = query.Model(*new(E))

	filter, err := pageable.Filter(query)
	if err != nil {
		return 0, err
	}

	var count int64
	results := filter.Count(&count)

	if err := results.Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *repository[E, C, M]) FindBy(ctx context.Context, where gormen.Where, preloads ...gormen.Preload) (nilo.Option[M], error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	for _, join := range where.Joins() {
		query = query.Joins(join)
	}

	for k, v := range where.Conditions() {
		switch v {
		case types.Or:
			query = query.Or(k.Get())
		default:
			query = query.Where(k.Get())
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

func (repository repository[E, _, _]) Count(ctx context.Context) (int64, error) {
	return repository.CountBy(ctx, gormen.Where{})
}

func (repository repository[E, _, _]) CountBy(ctx context.Context, where gormen.Where) (int64, error) {
	query := repository.db.WithContext(ctx)

	for _, join := range where.Joins() {
		query = query.Joins(join)
	}

	for k, v := range where.Conditions() {
		switch v {
		case types.Or:
			query = query.Or(k.Get())
		default:
			query = query.Where(k.Get())
		}
	}

	var count int64
	results := query.Model(*new(E)).Count(&count)

	if err := results.Error; err != nil {
		return 0, err
	}

	return count, nil
}
