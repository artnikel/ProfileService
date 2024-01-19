// Package repository is a lower level of project
package repository

import (
	"context"
	"fmt"

	"github.com/artnikel/ProfileService/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	berrors "github.com/artnikel/ProfileService/internal/errors"
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
		return fmt.Errorf("queryRow %w", err)
	}
	if count != 0 {
		return berrors.New(berrors.LoginAlreadyExist, "Login is occupied by another user")
	}
	_, err = p.pool.Exec(ctx, "INSERT INTO users (id, login, password) VALUES ($1, $2, $3)", user.ID, user.Login, user.Password)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}
	return nil
}

// GetByLogin get password and id of user.
func (p *PgRepository) GetByLogin(ctx context.Context, login string) ([]byte, uuid.UUID, error) {
	var id uuid.UUID
	var password []byte
	err := p.pool.QueryRow(ctx, "SELECT password, id FROM users WHERE login = $1", login).Scan(&password, &id)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("queryRow %w", err)
	}
	return password, id, nil
}

// DeleteAccount deleted account by id.
func (p *PgRepository) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	var count int
	err := p.pool.QueryRow(ctx, "SELECT COUNT(id) FROM users WHERE id = $1", id).Scan(&count)
	if err != nil {
		return fmt.Errorf("queryRow %w", err)
	}
	if count == 0 {
		return berrors.New(berrors.UserDoesntExists, "User doesnt exist")
	}
	_, err = p.pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}
	return nil
}
