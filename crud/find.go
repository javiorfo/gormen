package crud

import (
	"context"
	"errors"

	"github.com/javiorfo/gormix/pagination"
	"github.com/javiorfo/nilo"
	"github.com/javiorfo/steams"
	"gorm.io/gorm"
)

func (repository *hexaRepository[E, C, M]) FindBy(ctx context.Context, sqlField SqlField, preloads ...Preload) (nilo.Option[M], error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	var entity C = new(E)
	result := query.First(&entity, sqlField.Prepared(), sqlField.value)

	if err := result.Error; err != nil {
		none := nilo.None[M]()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return none, err
	}

	model := entity.Into()

	return nilo.Some(model), nil
}

func (repository *stdRepository[M]) FindBy(ctx context.Context, sqlField SqlField, preloads ...Preload) (nilo.Option[M], error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	var entity M
	result := query.First(&entity, sqlField.Prepared(), sqlField.value)

	if err := result.Error; err != nil {
		none := nilo.None[M]()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return none, err
	}

	return nilo.Some(entity), nil
}

func (repository *hexaRepository[E, C, M]) FindAll(ctx context.Context, pageable pagination.Pageable, preloads ...Preload) (*pagination.Page[M], error) {
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

func (repository hexaRepository[E, _, _]) count(ctx context.Context, pageable pagination.Pageable, preloads ...Preload) (int64, error) {
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
