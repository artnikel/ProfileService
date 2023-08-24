// Package service contains business logic of a project
package service

import (
	"context"
	"fmt"

	"github.com/artnikel/ProfileService/internal/model"
	"github.com/google/uuid"
)

// UserRepository is an interface that contains CRUD methods and GetAll
type UserRepository interface {
	SignUp(ctx context.Context, user *model.User) error
	GetByLogin(ctx context.Context, username string) ([]byte, uuid.UUID, error)
	AddRefreshToken(ctx context.Context, id uuid.UUID, refreshToken string) error
	GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (string, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

// UserService contains UserRepository interface
type UserService struct {
	uRep UserRepository
}

// NewUserService accepts UserRepository object and returnes an object of type *UserService
func NewUserService(uRep UserRepository) *UserService {
	return &UserService{uRep: uRep}
}

// SignUp is a method of UserService that calls  method of Repository
func (us *UserService) SignUp(ctx context.Context, user *model.User) error {
	err := us.uRep.SignUp(ctx, user)
	if err != nil {
		return fmt.Errorf("UserService-SignUp: error: %w", err)
	}
	return nil
}

// GetByLogin is a method of UserService that calls method of Repository
func (us *UserService) GetByLogin(ctx context.Context, login string) ([]byte, uuid.UUID, error) {
	hash, id, err := us.uRep.GetByLogin(ctx, login)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("UserService-GetByLogin: error: %w", err)
	}
	return hash, id, nil
}

// GetRefreshTokenByID is a method of UserService that calls method of Repository
func (us *UserService) GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (string, error) {
	refreshToken, err := us.uRep.GetRefreshTokenByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("UserService-GetRefreshTokenByID: error: %w", err)
	}
	return refreshToken, nil
}

// AddRefreshToken is a method of UserService that calls method of Repository
func (us *UserService) AddRefreshToken(ctx context.Context, id uuid.UUID, refreshToken string) error {
	err := us.uRep.AddRefreshToken(ctx, id, refreshToken)
	if err != nil {
		return fmt.Errorf("UserService-GetRefreshTokenByID: error: %w", err)
	}
	return nil
}

// DeleteAccount is a method from UserService that deleted account by id
func (us *UserService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	err := us.uRep.DeleteAccount(ctx, id)
	if err != nil {
		return fmt.Errorf("UserService-DeleteAccount: error%w", err)
	}
	return nil
}
