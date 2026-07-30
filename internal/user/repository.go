package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/you/gobase/internal/httputil"
)

type UserRepo struct {
	db *sqlx.DB
}

func (r *UserRepo) Create(ctx context.Context, u *User) error {
	query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query, u.Name, u.Email).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *UserRepo) FindByID(ctx context.Context, id string) (*User, error) {
	var u User
	query := `SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.GetContext(ctx, &u, query, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user %s: %w", id, httputil.ErrNotFound)
	}
	return &u, err
}

func (r *UserRepo) List(ctx context.Context) ([]*User, error) {
	var users []*User
	query := `SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC`
	err := r.db.SelectContext(ctx, &users, query)
	return users, err
}

func (r *UserRepo) Update(ctx context.Context, u *User) error {
	query := `UPDATE users SET name = $1, email = $2, updated_at = NOW() WHERE id = $3`
	result, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.ID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user %s: %w", u.ID, httputil.ErrNotFound)
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user %s: %w", id, httputil.ErrNotFound)
	}
	return nil
}
