package service

import (
	"Server/internal/entities/auth"
	"Server/internal/entities/stats"
	"Server/internal/entities/url"
	"Server/internal/entities/user"
	"Server/internal/storage"
	"context"
)

type UserService interface {
	Register(ctx context.Context, u *user.User) (*user.User, error)
	Login(ctx context.Context, req *auth.TokenRequest) (*auth.TokenResponse, error)
	RefreshToken(ctx context.Context, req *auth.TokenRequest) (*auth.TokenResponse, error)
	RevokeToken(ctx context.Context, token string) error
	GetUserByID(ctx context.Context, id string) (*user.User, error)
}

type URLService interface {
	ShortenURL(ctx context.Context, userID string, originalURL string, customCode string) (*url.ShortURL, *url.LongURL, error)
	GetLongURL(ctx context.Context, code string) (string, error)
	GetUserURLs(ctx context.Context, userID string, limit, offset int) ([]URLDetail, int64, error)
	GetURLDetail(ctx context.Context, userID string, shortURLID uint64) (*URLDetail, error)
	DeleteURL(ctx context.Context, userID string, shortURLID uint64) error
	GetAnalytics(ctx context.Context, userID string, shortURLID uint64) (*stats.UrlStats, error)
	RecordClick(ctx context.Context, code string, ipAddress string, userAgent string)
}

// URLDetail is a joined view for API responses
type URLDetail struct {
	ShortURLID uint64 `json:"id"`
	Code       string `json:"code"`
	LongURL    string `json:"long_url"`
	CreatedAt  string `json:"created_at"`
	ExpiresAt  string `json:"expires_at,omitempty"`
}

type Service struct {
	User UserService
	URL  URLService
}

func NewService(store *storage.Storage, jwtSecret string) *Service {
	return &Service{
		User: NewUserService(store.User, jwtSecret),
		URL:  NewURLService(store.LongURL, store.ShortURL, store.Stats, store.Cache),
	}
}
