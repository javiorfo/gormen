package std

import (
	"context"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/internal/types"
)

// Create inserts the given model into the database using the GORM Create method.
func (repository *repository[M]) Create(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, types.Create)
}

// Save creates or updates the given model in the database using the GORM Save method.
func (repository *repository[M]) Save(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, types.Save)
}

// Delete removes the given model from the database using the GORM Delete method.
func (repository *repository[M]) Delete(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, types.Delete)
}

// cud is a private helper that executes create, update/save, or delete database operations
// based on the specified method constant.
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

// CreateAll inserts multiple models in batches specified by batchSize.
func (repository *repository[M]) CreateAll(ctx context.Context, models *[]M, batchSize int) error {
	result := repository.db.WithContext(ctx).CreateInBatches(models, batchSize)

	if err := result.Error; err != nil {
		return err
	}

	return nil
}

// DeleteAll removes all given models from the database.
func (repository *repository[M]) DeleteAll(ctx context.Context, models []M) error {
	result := repository.db.WithContext(ctx).Delete(models)

	if err := result.Error; err != nil {
		return err
	}

	return nil
}

// DeleteAllBy deletes all records matching the conditions and joins defined in the Where clause.
func (repository repository[M]) DeleteAllBy(ctx context.Context, where gormen.Where) error {
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

	// Delete all matching records of model type M
	query = query.Delete(*new(M))

	if err := query.Error; err != nil {
		return err
	}

	return nil
}

// SaveAll creates or updates multiple models in the database.
func (repository *repository[M]) SaveAll(ctx context.Context, models []M) error {
	result := repository.db.WithContext(ctx).Save(models)

	if err := result.Error; err != nil {
		return err
	}

	return nil
}
