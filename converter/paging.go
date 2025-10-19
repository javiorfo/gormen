package converter

import (
	"context"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/steams"
)

// TODO 
// WithPreload() WithWhere

func (repository *converterRepository[E, C, M]) FindAllPaginated(ctx context.Context, pageable pagination.Pageable, preloads ...gormen.Preload) (*pagination.Page[M], error) {
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

	var entities []E
	page, err := pageable.Paginate(query)
	if err != nil {
		return nil, err
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

func (repository converterRepository[E, _, _]) count(ctx context.Context, pageable pagination.Pageable, preloads ...gormen.Preload) (int64, error) {
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
