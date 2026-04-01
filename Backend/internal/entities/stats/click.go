package stats

import "time"

type Click struct {
	ID         uint64    `json:"id" db:"id"`
	ShortURLID uint64    `json:"short_url_id" db:"short_url_id"`
	IpAddress  string    `json:"ip_address" db:"ip_address"`
	UserAgent  string    `json:"user_agent" db:"user_agent"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
