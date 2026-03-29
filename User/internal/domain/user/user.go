package user

import "time"

type User struct {
	UID       int64     `json:"uid"`
	Name      string    `json:"name"`
	UserName  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Zip       string    `json:"zip"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) New() *User {
	return &User{}
}
