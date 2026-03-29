package postgres

import (
	"User/internal/repository"
	"github.com/jmoiron/sqlx"
)

type postgresDB struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) repository.UserRepo {
	return &postgresDB{
		db: db,
	}
}

func NewUrlPostgresRepo(db *sqlx.DB) repository.UrlRepo {
	return &postgresDB{
		db: db,
	}
}
