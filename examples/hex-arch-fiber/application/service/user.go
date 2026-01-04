package service

import (
	"context"
	"errors"
	"hex-arch-fiber/application/model"
	"hex-arch-fiber/port"

	"github.com/javiorfo/gormen/pagination"
)

type userService struct {
	userRepo port.UserRepository
}

func NewUserService(r port.UserRepository) port.UserService {
	return &userService{userRepo: r}
}

func (u *userService) FindAll(ctx context.Context, pageable pagination.Pageable) (*pagination.Page[model.User], error) {
	return u.userRepo.FindAll(ctx, pageable)
}

func (u *userService) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := u.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
