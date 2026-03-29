package postgres

import (
	"User/internal/domain/user"
	"context"
	"database/sql"
	"errors"
)

func (r *postgresDB) CreateUser(ctx context.Context, u *user.User) (int64, error) {
	query := `
		INSERT INTO users (name, username, password, email, phone, address, city, state, zip, country, created_at, updated_at)
		VALUES (:name, :username, :password, :email, :phone, :address, :city, :state, :zip, :country, :created_at, :updated_at)
		RETURNING id`
	
	rows, err := r.db.NamedQueryContext(ctx, query, u)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var id int64
	if rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return 0, err
		}
	}
	
	u.UID = id
	return id, nil
}

func (r *postgresDB) GetUserByMobileNo(ctx context.Context, mobileNo string) (*user.User, error) {
	query := `
		SELECT id as uid, name, username, password, email, phone, address, city, state, zip, country, created_at, updated_at
		FROM users
		WHERE phone = $1`
	
	u := &user.User{}
	err := r.db.GetContext(ctx, u, query, mobileNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil 
		}
		return nil, err
	}
	
	return u, nil
}

func (r *postgresDB) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	query := `
		SELECT id as uid, name, username, password, email, phone, address, city, state, zip, country, created_at, updated_at
		FROM users
		WHERE email = $1`
	
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

func (r *postgresDB) GetUserByID(ctx context.Context, id int64) (*user.User, error) {
	query := `
		SELECT id as uid, name, username, password, email, phone, address, city, state, zip, country, created_at, updated_at
		FROM users
		WHERE id = $1`
	
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

func (r *postgresDB) UpdateUser(ctx context.Context, u *user.User) (int64, error) {
	query := `
		UPDATE users
		SET name = :name, username = :username, password = :password, email = :email, phone = :phone, address = :address, city = :city, state = :state, zip = :zip, country = :country, updated_at = :updated_at
		WHERE id = :uid`
	
	result, err := r.db.NamedExecContext(ctx, query, u)
	if err != nil {
		return 0, err
	}
	
	return result.RowsAffected()
}

func (r *postgresDB) DeleteUser(ctx context.Context, id int64) (int64, error) {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, err
	}
	
	return result.RowsAffected()
}
