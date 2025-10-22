package std

import (
	"context"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/internal/types"
)

func (repository *repository[M]) Create(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, types.Create)
}

func (repository *repository[M]) Save(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, types.Save)
}

func (repository *repository[M]) Delete(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, types.Delete)
}

func (repository *repository[M]) cud(ctx context.Context, model *M, method int) error {
	result := repository.db.WithContext(ctx)

	switch method {
	case types.Save:
		result = result.Save(model)
	case types.Create:
		result = result.Create(model)
	case types.Delete:
		result = result.Delete(model)
	}

	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func (repository *repository[M]) CreateAll(ctx context.Context, models *[]M, batchSize int) error {
	result := repository.db.WithContext(ctx).CreateInBatches(&models, batchSize)

	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func (repository *repository[M]) DeleteAll(ctx context.Context, models []M) error {
	result := repository.db.WithContext(ctx).Delete(&models)

	if err := result.Error; err != nil {
		return err
	}

	return nil
}

func (repository repository[M]) DeleteAllBy(ctx context.Context, where gormen.Where) error {
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

	query = query.Delete(*new(M))

	if err := query.Error; err != nil {
		return err
	}

	return nil
}

func (repository *repository[M]) SaveAll(ctx context.Context, models []M) error {
	result := repository.db.WithContext(ctx).Save(&models)

	if err := result.Error; err != nil {
		return err
	}

	return nil
}
