package user

import (
	"time"
)

type User struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	UserName  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Phone     string    `json:"phone" db:"phone"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
