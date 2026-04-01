package postgres

import (
	"Server/internal/entities/user"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(ctx context.Context, u *user.User) error {
	query := `
		INSERT INTO users (id, name, username, password, email, phone, created_at, updated_at)
		VALUES (:id, :name, :username, :password, :email, :phone, :created_at, :updated_at)`
	
	_, err := r.db.NamedExecContext(ctx, query, u)
	return err
}

func (r *userRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `SELECT id, name, username, password, email, phone, created_at, updated_at FROM users WHERE email = $1`
	
	u := &user.User{}
	err := r.db.GetContext(ctx, u, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (r *userRepo) GetUserByUsername(ctx context.Context, username string) (*user.User, error) {
	query := `SELECT id, name, username, password, email, phone, created_at, updated_at FROM users WHERE username = $1`
	
	u := &user.User{}
	err := r.db.GetContext(ctx, u, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	query := `SELECT id, name, username, password, email, phone, created_at, updated_at FROM users WHERE id = $1`
	
	u := &user.User{}
	err := r.db.GetContext(ctx, u, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}
