package paging

import (
	"context"

	"github.com/javiorfo/gormix/crud"
	"github.com/javiorfo/gormix/pagination"
)

type Repository[M any] interface {
	FindAll(ctx context.Context, pageable pagination.Pageable, preloads ...crud.Preload) (*pagination.Page[M], error)
}
