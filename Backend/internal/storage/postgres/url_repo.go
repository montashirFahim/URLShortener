package postgres

import (
	"Server/internal/entities/url"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type longURLRepo struct {
	db *sqlx.DB
}

func NewLongURLRepo(db *sqlx.DB) *longURLRepo {
	return &longURLRepo{db: db}
}

func (r *longURLRepo) FindOrCreate(ctx context.Context, rawURL string) (*url.LongURL, error) {
	// Try to find existing
	existing := &url.LongURL{}
	err := r.db.GetContext(ctx, existing, "SELECT id, url, created_at FROM long_urls WHERE url = $1", rawURL)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Insert new
	lu := &url.LongURL{}
	err = r.db.QueryRowxContext(ctx,
		"INSERT INTO long_urls (url) VALUES ($1) RETURNING id, url, created_at", rawURL).StructScan(lu)
	if err != nil {
		return nil, err
	}
	return lu, nil
}

func (r *longURLRepo) GetByID(ctx context.Context, id uint64) (*url.LongURL, error) {
	lu := &url.LongURL{}
	err := r.db.GetContext(ctx, lu, "SELECT id, url, created_at FROM long_urls WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return lu, nil
}

// --- ShortURL Repo ---

type shortURLRepo struct {
	db *sqlx.DB
}

func NewShortURLRepo(db *sqlx.DB) *shortURLRepo {
	return &shortURLRepo{db: db}
}

func (r *shortURLRepo) Create(ctx context.Context, s *url.ShortURL) error {
	err := r.db.QueryRowxContext(ctx,
		`INSERT INTO short_urls (code, long_url_id, user_id, expires_at)
		 VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		s.Code, s.LongURLID, s.UserID, s.ExpiresAt).Scan(&s.ID, &s.CreatedAt)
	return err
}

func (r *shortURLRepo) GetByCode(ctx context.Context, code string) (*url.ShortURL, error) {
	s := &url.ShortURL{}
	err := r.db.GetContext(ctx, s,
		"SELECT id, code, long_url_id, user_id, expires_at, created_at FROM short_urls WHERE code = $1", code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func (r *shortURLRepo) GetByID(ctx context.Context, id uint64) (*url.ShortURL, error) {
	s := &url.ShortURL{}
	err := r.db.GetContext(ctx, s,
		"SELECT id, code, long_url_id, user_id, expires_at, created_at FROM short_urls WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return s, nil
}

func (r *shortURLRepo) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*url.ShortURL, int64, error) {
	var total int64
	err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) FROM short_urls WHERE user_id = $1", userID)
	if err != nil {
		return nil, 0, err
	}

	var urls []*url.ShortURL
	err = r.db.SelectContext(ctx, &urls,
		`SELECT id, code, long_url_id, user_id, expires_at, created_at
		 FROM short_urls WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return urls, total, nil
}

func (r *shortURLRepo) Delete(ctx context.Context, userID string, id uint64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM short_urls WHERE id = $1 AND user_id = $2", id, userID)
	return err
}
