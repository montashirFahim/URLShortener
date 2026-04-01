package url

import "time"

type ShortURL struct {
	ID        uint64    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	LongURLID uint64    `json:"long_url_id" db:"long_url_id"`
	Code      string    `json:"code" db:"code"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
