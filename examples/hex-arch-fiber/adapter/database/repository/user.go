package repository

import (
	"context"
	"hex-arch-fiber/adapter/database/entities"
	"hex-arch-fiber/application/model"
	"hex-arch-fiber/port"

	"github.com/javiorfo/gormen"
	"github.com/javiorfo/gormen/converter"
	"github.com/javiorfo/gormen/pagination"
	"github.com/javiorfo/gormen/where"
	"gorm.io/gorm"
)

type userRepo struct {
	*gorm.DB
	gr gormen.Repository[model.User]
}

func NewUserRepo(db *gorm.DB) port.UserRepository {
	return &userRepo{
		DB: db,
		gr: converter.NewRepository[entities.UserDB, *entities.UserDB](db),
	}
}

func (u *userRepo) FindAll(ctx context.Context, pageable pagination.Pageable) (*pagination.Page[model.User], error) {
	return u.gr.FindAllPaginatedBy(ctx, pageable, gormen.NewWhere(where.Equal("enable", true)).Build(), "Person")
}

func (u *userRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	return u.gr.FindBy(ctx, gormen.NewWhere(where.Equal("username", username)).Build(), "Person")
}
