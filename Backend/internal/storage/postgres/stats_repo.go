package postgres

import (
	"Server/internal/entities/stats"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type statsRepo struct {
	db *sqlx.DB
}

func NewStatsRepo(db *sqlx.DB) *statsRepo {
	return &statsRepo{db: db}
}

func (r *statsRepo) RecordClick(ctx context.Context, click *stats.Click) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO clicks (short_url_id, ip_address, user_agent) VALUES ($1, $2, $3)`,
		click.ShortURLID, click.IpAddress, click.UserAgent)
	return err
}

func (r *statsRepo) IncrementStats(ctx context.Context, shortURLID uint64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO url_stats (short_url_id, clicks, last_access)
		 VALUES ($1, 1, $2)
		 ON CONFLICT (short_url_id)
		 DO UPDATE SET clicks = url_stats.clicks + 1, last_access = $2`,
		shortURLID, time.Now())
	return err
}

func (r *statsRepo) GetStats(ctx context.Context, shortURLID uint64) (*stats.UrlStats, error) {
	s := &stats.UrlStats{}
	err := r.db.GetContext(ctx, s,
		"SELECT short_url_id, clicks, last_access FROM url_stats WHERE short_url_id = $1", shortURLID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &stats.UrlStats{ShortURLID: shortURLID, Clicks: 0}, nil
		}
		return nil, err
	}
	return s, nil
}
