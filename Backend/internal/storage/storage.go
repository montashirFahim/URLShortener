package storage

import (
	"Server/internal/entities/stats"
	"Server/internal/entities/url"
	"Server/internal/entities/user"
	"context"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *user.User) error
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	GetUserByID(ctx context.Context, id string) (*user.User, error)
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
}

type LongURLRepo interface {
	FindOrCreate(ctx context.Context, rawURL string) (*url.LongURL, error)
	GetByID(ctx context.Context, id uint64) (*url.LongURL, error)
}

type ShortURLRepo interface {
	Create(ctx context.Context, s *url.ShortURL) error
	GetByCode(ctx context.Context, code string) (*url.ShortURL, error)
	GetByID(ctx context.Context, id uint64) (*url.ShortURL, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*url.ShortURL, int64, error)
	Delete(ctx context.Context, userID string, id uint64) error
}

type StatsRepo interface {
	RecordClick(ctx context.Context, click *stats.Click) error
	IncrementStats(ctx context.Context, shortURLID uint64) error
	GetStats(ctx context.Context, shortURLID uint64) (*stats.UrlStats, error)
}

// URLCache caches the code -> long_url mapping for fast redirects
type URLCache interface {
	GetLongURL(ctx context.Context, code string) (string, error)
	SetLongURL(ctx context.Context, code string, longURL string, expiresAt *string) error
	DeleteURL(ctx context.Context, code string) error
}

type Storage struct {
	User     UserRepo
	LongURL  LongURLRepo
	ShortURL ShortURLRepo
	Stats    StatsRepo
	Cache    URLCache
}

func NewStorage(userRepo UserRepo, longURLRepo LongURLRepo, shortURLRepo ShortURLRepo, statsRepo StatsRepo, urlCache URLCache) *Storage {
	return &Storage{
		User:     userRepo,
		LongURL:  longURLRepo,
		ShortURL: shortURLRepo,
		Stats:    statsRepo,
		Cache:    urlCache,
	}
}
