package std

/* import (
	"context"
	"errors"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/nilo"
	"gorm.io/gorm"
)

func (repository *stdRepository[M]) FindBy(ctx context.Context, sqlField gormen.SqlField, preloads ...gormen.Preload) (nilo.Option[M], error) {
	query := repository.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	var entity M
	result := query.First(&entity, sqlField.Prepared(), sqlField.Value())

	if err := result.Error; err != nil {
		none := nilo.None[M]()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
		}
		return none, err
	}

	return nilo.Some(entity), nil
} */
