package port

import (
	"context"
	"hex-arch-fiber/application/model"

	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/nilo"
)

type UserRepository interface {
	FindAll(ctx context.Context, pageable pagination.Pageable) (*pagination.Page[model.User], error)
	FindByUsername(ctx context.Context, username string) (nilo.Option[model.User], error)
}

type UserService interface {
	FindAll(ctx context.Context, pageable pagination.Pageable) (*pagination.Page[model.User], error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}
