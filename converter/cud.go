package converter

import (
	"context"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/internal"
	"github.com/javiorfo/steams"
)

const (
	create = iota
	delete
	save
)

func (repository *repository[E, C, M]) Create(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, create)
}

func (repository *repository[E, C, M]) Save(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, save)
}

func (repository *repository[E, C, M]) Delete(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, delete)
}

func (repository *repository[E, C, M]) cud(ctx context.Context, model *M, method int) error {
	var entity C = new(E)
	entity.From(*model)

	result := repository.db.WithContext(ctx)

	switch method {
	case save:
		result = result.Save(&entity)
	case create:
		result = result.Create(&entity)
	case delete:
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

	for k, v := range where {
		switch v {
		case internal.Or:
			query = query.Or(k.Prepared(), k.Value())
		default:
			query = query.Where(k.Prepared(), k.Value())
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
