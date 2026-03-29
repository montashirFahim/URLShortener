package url

import "time"

type URL struct {
	ID          string    `json:"id"`
	UID         int64     `json:"uid"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	Count       int64     `json:"count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u *URL) New() *URL {
	return &URL{}
}
