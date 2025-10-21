package std

import (
	"context"
	"errors"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/internal/types"
	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/gormen/pagination/sort"
	"github.com/javiorfo/nilo"
	"gorm.io/gorm"
)

func (repository *repository[M]) FindAllPaginated(ctx context.Context, pageable pagination.Pageable, preloads ...gormen.Preload) (*pagination.Page[M], error) {
	return repository.FindAllPaginatedBy(ctx, pageable, gormen.Where{}, preloads...)
}

func (repository *repository[M]) FindAllPaginatedBy(ctx context.Context, pageable pagination.Pageable, where gormen.Where, preloads ...gormen.Preload) (*pagination.Page[M], error) {
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

	var entities []M
	page, err := pageable.Paginate(query)
	if err != nil {
		return nil, err
	}

	for k, v := range where.Conditions() {
		switch v {
		case types.Or:
			query = query.Or(k.Prepared(), k.Value())
		default:
			query = query.Where(k.Prepared(), k.Value())
		}
	}

	results := page.Find(&entities)

	if err := results.Error; err != nil {
		return nil, err
	}

	return &pagination.Page[M]{Total: total, Elements: entities}, nil
}

func (repository *repository[M]) FindAll(ctx context.Context, preloads ...gormen.Preload) ([]M, error) {
	return repository.FindAllBy(ctx, gormen.Where{}, preloads...)
}

func (repository *repository[M]) FindAllBy(ctx context.Context, where gormen.Where, preloads ...gormen.Preload) ([]M, error) {
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
			query = query.Or(k.Prepared(), k.Value())
		default:
			query = query.Where(k.Prepared(), k.Value())
		}
	}

	var entities []M
	results := query.Find(&entities)

	if err := results.Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (repository *repository[M]) FindAllOrdered(ctx context.Context, orders []sort.Order, preloads ...gormen.Preload) ([]M, error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	for _, o := range orders {
		query = query.Order(o.Prepared())
	}

	var entities []M
	results := query.Find(&entities)

	if err := results.Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (repository repository[M]) count(ctx context.Context, pageable pagination.Pageable, preloads ...gormen.Preload) (int64, error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	query = query.Model(*new(M))

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

func (repository *repository[M]) FindBy(ctx context.Context, where gormen.Where, preloads ...gormen.Preload) (nilo.Option[M], error) {
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
			query = query.Or(k.Prepared(), k.Value())
		default:
			query = query.Where(k.Prepared(), k.Value())
		}
	}

	var entity M
	result := query.First(&entity)

	if err := result.Error; err != nil {
		none := nilo.None[M]()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return none, nil
		}
		return none, err
	}

	return nilo.Some(entity), nil
}

func (repository repository[_]) Count(ctx context.Context) (int64, error) {
	return repository.CountBy(ctx, gormen.Where{})
}

func (repository repository[M]) CountBy(ctx context.Context, where gormen.Where) (int64, error) {
	query := repository.db.WithContext(ctx)

	for _, join := range where.Joins() {
		query = query.Joins(join)
	}

	for k, v := range where.Conditions() {
		switch v {
		case types.Or:
			query = query.Or(k.Prepared(), k.Value())
		default:
			query = query.Where(k.Prepared(), k.Value())
		}
	}

	var count int64
	results := query.Model(*new(M)).Count(&count)

	if err := results.Error; err != nil {
		return 0, err
	}

	return count, nil
}
