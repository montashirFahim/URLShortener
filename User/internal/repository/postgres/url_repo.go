package postgres

import (
	"User/internal/domain/url"
	"context"
	"database/sql"
	"errors"
)

func (r *postgresDB) CreateUrl(ctx context.Context, u *url.URL) (int64, error) {
	query := `
		INSERT INTO urls (id, uid, original_url, short_url, created_at, updated_at)
		VALUES (:id, :uid, :original_url, :short_url, :created_at, :updated_at)`
	
	result, err := r.db.NamedExecContext(ctx, query, u)
	if err != nil {
		return 0, err
	}
	
	return result.RowsAffected()
}

func (r *postgresDB) GetUrlByID(ctx context.Context, id string) (*url.URL, error) {
	query := `
		SELECT id, uid, original_url, short_url, created_at, updated_at
		FROM urls
		WHERE id = $1`
	
	u := &url.URL{}
	err := r.db.GetContext(ctx, u, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	return u, nil
}

func (r *postgresDB) GetUrlsByUserID(ctx context.Context, userID int64, limit, offset int) ([]*url.URL, error) {
	query := `
		SELECT id, uid, original_url, short_url, created_at, updated_at
		FROM urls
		WHERE uid = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`
	
	var urls []*url.URL
	err := r.db.SelectContext(ctx, &urls, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	
	return urls, nil
}

func (r *postgresDB) UpdateUrl(ctx context.Context, u *url.URL) (int64, error) {
	query := `
		UPDATE urls
		SET uid = :uid, original_url = :original_url, short_url = :short_url, updated_at = :updated_at
		WHERE id = :id`
	
	result, err := r.db.NamedExecContext(ctx, query, u)
	if err != nil {
		return 0, err
	}
	
	return result.RowsAffected()
}

func (r *postgresDB) DeleteUrl(ctx context.Context, id string) (int64, error) {
	query := `DELETE FROM urls WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, err
	}
	
	return result.RowsAffected()
}
