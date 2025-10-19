package gormen

import (
	"context"

	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/gormen/pagination/sort"
	"github.com/javiorfo/nilo"
)

type Repository[M any] interface {
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, model *M) error
	CreateAll(ctx context.Context, model []M) error
	Delete(ctx context.Context, model *M) error
	DeleteAll(ctx context.Context, model []M) error
	DeleteAllBy(ctx context.Context, sqlField []SqlField) error
	ExistsBy(ctx context.Context, sqlField SqlField) (bool, error)
	FindBy(ctx context.Context, sqlField SqlField, preloads ...Preload) (nilo.Option[M], error)
	FindAllBy(ctx context.Context, sqlField []SqlField, preloads ...Preload) ([]M, error)
	FindAllPaginated(ctx context.Context, pageable pagination.Pageable, preloads ...Preload) (*pagination.Page[M], error)
	FindAllOrdered(ctx context.Context, orders []sort.Order, preloads ...Preload) ([]M, error)
	Save(ctx context.Context, model *M) error
	SaveAll(ctx context.Context, model []M) error
}
