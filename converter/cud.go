package converter

import (
	"context"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/internal/types"
	"github.com/javiorfo/steams"
)

// Create converts the model M into the entity E using the converter C,
// then creates the entity in the database using GORM.
func (repository *repository[E, C, M]) Create(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, types.Create)
}

// Save converts the model M into the entity E, then saves (creates or updates)
// the entity in the database using GORM.
func (repository *repository[E, C, M]) Save(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, types.Save)
}

// Delete converts the model M into the entity E, then deletes
// the entity from the database using GORM.
func (repository *repository[E, C, M]) Delete(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, types.Delete)
}

// cud is a helper that converts the model to an entity,
// performs the create, save or delete operation using GORM,
// handles any error, and converts back the entity into the model.
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

// CreateAll converts multiple models into entities,
// then creates them in batches in the database.
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

// DeleteAll converts multiple models into entities,
// then deletes them from the database.
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

// DeleteAllBy deletes all entities of type E matching the given Where clause conditions.
func (repository repository[E, _, _]) DeleteAllBy(ctx context.Context, where gormen.Where) error {
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

	query = query.Delete(*new(E))
	if err := query.Error; err != nil {
		return err
	}

	return nil
}

// SaveAll converts multiple models into entities,
// then saves (creates or updates) them in the database.
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
