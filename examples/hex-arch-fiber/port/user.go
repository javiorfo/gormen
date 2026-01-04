package port

import (
	"context"
	"hex-arch-fiber/application/model"

	"github.com/javiorfo/gormen/pagination"
)

type UserRepository interface {
	FindAll(ctx context.Context, pageable pagination.Pageable) (*pagination.Page[model.User], error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}

type UserService interface {
	FindAll(ctx context.Context, pageable pagination.Pageable) (*pagination.Page[model.User], error)
	FindByUsername(ctx context.Context, username string) (*model.User, error)
}
