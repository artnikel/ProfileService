// Package service contains business logic of a project
package service

import (
	"context"
	"errors"
	"fmt"

	berrors "github.com/artnikel/ProfileService/internal/errors"
	"github.com/artnikel/ProfileService/internal/model"
	"github.com/google/uuid"
)

// UserRepository is an interface that contains CRUD methods and GetAll
type UserRepository interface {
	SignUp(ctx context.Context, user *model.User) error
	GetByLogin(ctx context.Context, username string) ([]byte, uuid.UUID, error)
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
		var e *berrors.BusinessError
		if errors.As(err, &e) {
			return err
		}
		return fmt.Errorf("signUp %w", err)
	}
	return nil
}

// GetByLogin is a method of UserService that calls method of Repository
func (us *UserService) GetByLogin(ctx context.Context, login string) ([]byte, uuid.UUID, error) {
	hash, id, err := us.uRep.GetByLogin(ctx, login)
	if err != nil {
		return nil, uuid.Nil, fmt.Errorf("getByLogin %w", err)
	}
	return hash, id, nil
}

// DeleteAccount is a method from UserService that deleted account by id
func (us *UserService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	err := us.uRep.DeleteAccount(ctx, id)
	if err != nil {
		var e *berrors.BusinessError
		if errors.As(err, &e) {
			return err
		}
		return fmt.Errorf("deleteAccount %w", err)
	}
	return nil
}
