package service

import (
	"User/internal/domain/url"
	"User/internal/domain/user"
	"context"
)

type UserService interface {
	Register(ctx context.Context, u *user.User) (int64, error)
	Login(ctx context.Context, username, password string) (string, error)
	GetUserByID(ctx context.Context, id int64) (*user.User, error)
}

type UrlService interface {
	GetUserUrls(ctx context.Context, userID int64, limit, page int) ([]*url.URL, error)
}

type Service struct {
	User UserService
	Url  UrlService
}

func NewService(userSvc UserService, urlSvc UrlService) *Service {
	return &Service{
		User: userSvc,
		Url:  urlSvc,
	}
}
