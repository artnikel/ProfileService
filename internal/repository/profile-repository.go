// Package repository is a lower level of project
package repository

import (
	"context"
	"fmt"

	"github.com/artnikel/ProfileService/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgRepository represents the PostgreSQL repository implementation.
type PgRepository struct {
	pool *pgxpool.Pool
}

// NewPgRepository creates and returns a new instance of PgRepository, using the provided pgxpool.Pool.
func NewPgRepository(pool *pgxpool.Pool) *PgRepository {
	return &PgRepository{
		pool: pool,
	}
}

// SignUp creates a new user record in the database.
func (p *PgRepository) SignUp(ctx context.Context, user *model.User) error {
	var count int
	err := p.pool.QueryRow(ctx, "SELECT COUNT(id) FROM users WHERE login = $1", user.Login).Scan(&count)
	if err != nil {
		return fmt.Errorf("PgRepository-SignUpUser: error in method r.pool.QuerryRow(): %w", err)
	}
	if count != 0 {
		return fmt.Errorf("PgRepository-SignUpUser: the login is occupied by another user")
	}
	_, err = p.pool.Exec(ctx, "INSERT INTO users (id, login, password) VALUES ($1, $2, $3)", user.ID, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("PgRepository-SignUpUser: error in method r.pool.Exec(): %w", err)
	}
	return nil
}

// GetByLogin get password and id of user.
func (p *PgRepository) GetByLogin(ctx context.Context, login string) ([]byte, uuid.UUID, error) {
	var id uuid.UUID
	var password []byte
	err := p.pool.QueryRow(ctx, "SELECT password, id FROM users WHERE login = $1", login).Scan(&password, &id)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("PgRepository-GetByLOgin: error in method r.pool.QuerryRow(): %w", err)
	}
	return password, id, nil
}

// AddRefreshToken adds a token to the user's record in the database.
func (p *PgRepository) AddRefreshToken(ctx context.Context, id uuid.UUID, refreshToken string) error {
	var count int
	err := p.pool.QueryRow(ctx, "SELECT COUNT(id) FROM users WHERE id = $1", id).Scan(&count)
	if err != nil {
		return fmt.Errorf("PgRepository-AddRefreshToken: error in method r.pool.QuerryRow(): %w", err)
	}
	if count == 0 {
		return fmt.Errorf("PgRepository-DeleteAccount: cannot add refresh token to non-existent user")
	}
	_, err = p.pool.Exec(ctx, "UPDATE users SET refreshtoken = $1 WHERE id = $2", refreshToken, id)
	if err != nil {
		return fmt.Errorf("PgRepository-AddRefreshToken : r.pool.Exec(): %w", err)
	}
	return nil
}

// GetRefreshTokenByID returns refresh token by id.
func (p *PgRepository) GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (string, error) {
	var refreshToken string
	err := p.pool.QueryRow(ctx, "SELECT refreshtoken FROM users WHERE id = $1", id).Scan(&refreshToken)
	if err != nil {
		return "", fmt.Errorf("PgRepository-GetRefreshTokenByID: error in method r.pool.QuerryRow(): %w", err)
	}
	return refreshToken, nil
}

// DeleteAccount deleted account by id.
func (p *PgRepository) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	var count int
	err := p.pool.QueryRow(ctx, "SELECT COUNT(id) FROM users WHERE id = $1", id).Scan(&count)
	if err != nil {
		return fmt.Errorf("PgRepository-DeleteAccount: error in method r.pool.QuerryRow(): %w", err)
	}
	if count == 0 {
		return fmt.Errorf("PgRepository-DeleteAccount: cannot delete non-existent user")
	}
	_, err = p.pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("PgRepository-DeleteAccount: error in method r.pool.Exec(): %w", err)
	}
	return nil
}
