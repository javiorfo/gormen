package gormen

import (
	"context"

	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/gormen/pagination/sort"
	"github.com/javiorfo/nilo"
)

type Repository[M any] interface {
	CudRepository[M]
	ReadRepository[M]
}

type CudRepository[M any] interface {
	Create(ctx context.Context, model *M) error
	CreateAll(ctx context.Context, model *[]M, batchSize int) error
	Delete(ctx context.Context, model *M) error
	DeleteAll(ctx context.Context, model []M) error
	DeleteAllBy(ctx context.Context, where Where) error
	Save(ctx context.Context, model *M) error
	SaveAll(ctx context.Context, model []M) error
}

type ReadRepository[M any] interface {
	Count(ctx context.Context) (int64, error)
	CountBy(ctx context.Context, where Where) (int64, error)
	FindBy(ctx context.Context, where Where, preloads ...Preload) (nilo.Option[M], error)
	FindAll(ctx context.Context, preloads ...Preload) ([]M, error)
	FindAllBy(ctx context.Context, where Where, preloads ...Preload) ([]M, error)
	FindAllPaginated(ctx context.Context, pageable pagination.Pageable, preloads ...Preload) (*pagination.Page[M], error)
	FindAllPaginatedBy(ctx context.Context, pageable pagination.Pageable, where Where, preloads ...Preload) (*pagination.Page[M], error)
	FindAllOrdered(ctx context.Context, orders []sort.Order, preloads ...Preload) ([]M, error)
}
