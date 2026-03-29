package service

import (
	"User/internal/domain/url"
	"User/internal/repository"
	"context"
)

type urlService struct {
	repo repository.UrlRepo
}

func NewUrlService(repo repository.UrlRepo) UrlService {
	return &urlService{repo: repo}
}

func (s *urlService) GetUserUrls(ctx context.Context, userID int64, limit, page int) ([]*url.URL, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	return s.repo.GetUrlsByUserID(ctx, userID, limit, offset)
}
