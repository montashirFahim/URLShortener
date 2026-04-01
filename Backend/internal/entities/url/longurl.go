package url

import "time"

type LongURL struct {
	ID        uint64    `json:"id" db:"id"`
	Url       string    `json:"url" db:"url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
