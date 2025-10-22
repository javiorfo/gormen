package converter

import (
	"context"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/internal/types"
	"github.com/javiorfo/steams"
)

func (repository *repository[E, C, M]) Create(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, types.Create)
}

func (repository *repository[E, C, M]) Save(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, types.Save)
}

func (repository *repository[E, C, M]) Delete(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, types.Delete)
}

func (repository *repository[E, C, M]) cud(ctx context.Context, model *M, method int) error {
	var entity C = new(E)
	entity.From(*model)

	result := repository.db.WithContext(ctx)

	switch method {
	case types.Save:
		result = result.Save(&entity)
	case types.Create:
		result = result.Create(&entity)
	case types.Delete:
		result = result.Delete(&entity)
	}

	if err := result.Error; err != nil {
		return err
	}

	*model = entity.Into()

	return nil
}

func (repository *repository[E, C, M]) CreateAll(ctx context.Context, models *[]M, batchSize int) error {
	entities := make([]C, len(*models))

	for i, model := range *models {
		var entity C = new(E)
		entity.From(model)
		entities[i] = entity
	}

	result := repository.db.WithContext(ctx).CreateInBatches(&entities, batchSize)

	if err := result.Error; err != nil {
		return err
	}

	*models = steams.Mapper(steams.OfSlice(entities), func(entity C) M {
		return entity.Into()
	}).Collect()

	return nil
}

func (repository *repository[E, C, M]) DeleteAll(ctx context.Context, models []M) error {
	entities := make([]C, len(models))

	for i, model := range models {
		var entity C = new(E)
		entity.From(model)
		entities[i] = entity
	}

	result := repository.db.WithContext(ctx).Delete(&entities)

	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func (repository repository[E, _, _]) DeleteAllBy(ctx context.Context, where gormen.Where) error {
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

	query = query.Delete(*new(E))

	if err := query.Error; err != nil {
		return err
	}

	return nil
}

func (repository *repository[E, C, M]) SaveAll(ctx context.Context, models []M) error {
	entities := make([]C, len(models))

	for i, model := range models {
		var entity C = new(E)
		entity.From(model)
		entities[i] = entity
	}

	result := repository.db.WithContext(ctx).Save(&entities)

	if err := result.Error; err != nil {
		return err
	}

	return nil
}
