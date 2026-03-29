package repository

import (
	"User/internal/domain/url"
	"User/internal/domain/user"
	"context"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *user.User) (int64, error)
	GetUserByMobileNo(ctx context.Context, mobileNo string) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	GetUserByID(ctx context.Context, id int64) (*user.User, error)
	UpdateUser(ctx context.Context, user *user.User) (int64, error)
	DeleteUser(ctx context.Context, id int64) (int64, error)
}

type UrlRepo interface {
	CreateUrl(ctx context.Context, url *url.URL) (int64, error)
	GetUrlByID(ctx context.Context, id string) (*url.URL, error)
	GetUrlsByUserID(ctx context.Context, userID int64, limit, offset int) ([]*url.URL, error)
	UpdateUrl(ctx context.Context, url *url.URL) (int64, error)
	DeleteUrl(ctx context.Context, id string) (int64, error)
}

type Repository struct {
	userRepository UserRepo
	urlRepository  UrlRepo
}

func NewRepository(userRepo UserRepo, urlRepo UrlRepo) *Repository {
	return &Repository{
		userRepository: userRepo,
		urlRepository:  urlRepo,
	}
}
