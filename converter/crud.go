package converter

import (
	"context"
	"errors"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/nilo"
	"gorm.io/gorm"
)

const (
	Create = iota
	Delete
	Save
)

func (repository converterRepository[E, _, _]) Count(ctx context.Context) (int64, error) {
	query := repository.db.WithContext(ctx).Model(*new(E))

	var count int64
	results := query.Count(&count)

	if err := results.Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *converterRepository[E, C, M]) Create(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, Create)
}

func (repository *converterRepository[E, C, M]) Save(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, Save)
}

func (repository *converterRepository[E, C, M]) Delete(ctx context.Context, model *M) error {
	return repository.cud(ctx, model, Delete)
}

func (repository *converterRepository[E, C, M]) cud(ctx context.Context, model *M, method int) error {
	var entity C = new(E)
	entity.From(*model)

	result := repository.db.WithContext(ctx)

	switch method {
	case Save:
		result = result.Save(&entity)
	case Create:
		result = result.Create(&entity)
	case Delete:
		result = result.Delete(&entity)
	}

	if err := result.Error; err != nil {
		return err
	}

	*model = entity.Into()

	return nil
}

func (repository *converterRepository[E, C, M]) FindBy(ctx context.Context, sqlField gormen.SqlField, preloads ...gormen.Preload) (nilo.Option[M], error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	var entity C = new(E)
	result := query.First(&entity, sqlField.Prepared(), sqlField.Value())

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
