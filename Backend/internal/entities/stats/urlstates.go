package stats

import "time"

type UrlStats struct {
	ShortURLID uint64    `json:"short_url_id" db:"short_url_id"`
	Clicks     uint64    `json:"clicks" db:"clicks"`
	LastAccess time.Time `json:"last_access" db:"last_access"`
}
