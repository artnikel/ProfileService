// Package handler is the top level of the application and it contains request handlers
package handler

import (
	"context"
	"fmt"

	"github.com/artnikel/ProfileService/internal/model"
	"github.com/artnikel/ProfileService/uproto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// UserService is an interface that contains methods of service for user
type UserService interface {
	SignUp(ctx context.Context, user *model.User) error
	GetByLogin(ctx context.Context, login string) ([]byte, uuid.UUID, error)
	AddRefreshToken(ctx context.Context, id uuid.UUID, refreshToken string) error
	GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (string, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

// EntityUser contains User Service interface
type EntityUser struct {
	srvcUser UserService
	validate *validator.Validate
	uproto.UnimplementedUserServiceServer
}

// NewEntityUser accepts User Service interface and returns an object of *EntityUser
func NewEntityUser(srvcUser UserService, validate *validator.Validate) *EntityUser {
	return &EntityUser{srvcUser: srvcUser, validate: validate}
}

// SignUp calls method of Service by handler
func (handl *EntityUser) SignUp(ctx context.Context, req *uproto.SignUpRequest) (*uproto.SignUpResponse, error) {
	createdUser := &model.User{
		ID:       uuid.New(),
		Login:    req.User.Login,
		Password: []byte(req.User.Password),
	}

	err := handl.validate.StructCtx(ctx, createdUser)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &uproto.SignUpResponse{}, fmt.Errorf("failed to validate")
	}

	err = handl.srvcUser.SignUp(ctx, createdUser)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &uproto.SignUpResponse{}, fmt.Errorf("failed to signUp")
	}

	return &uproto.SignUpResponse{
		Id: createdUser.ID.String(),
	}, nil
}

// GetByLogin calls method of Service by handler
func (handl *EntityUser) GetByLogin(ctx context.Context, req *uproto.GetByLoginRequest) (*uproto.GetByLoginResponse, error) {
	err := handl.validate.VarCtx(ctx, req.Login, "required,min=5,max=20")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &uproto.GetByLoginResponse{}, fmt.Errorf("failed to validate")
	}

	password, id, err := handl.srvcUser.GetByLogin(ctx, req.Login)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &uproto.GetByLoginResponse{}, fmt.Errorf("failed to get password and id by login")
	}

	return &uproto.GetByLoginResponse{
		Password: string(password),
		Id:       id.String(),
	}, nil
}

// AddRefreshToken calls method of Service by handler
func (handl *EntityUser) AddRefreshToken(ctx context.Context, req *uproto.AddRefreshTokenRequest) (*uproto.AddRefreshTokenResponse, error) {
	err := handl.validate.VarCtx(ctx, req.Id, "required,uuid")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &uproto.AddRefreshTokenResponse{}, fmt.Errorf("failed to validate")
	}
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &uproto.AddRefreshTokenResponse{}, fmt.Errorf("failed to parse id")
	}
	err = handl.srvcUser.AddRefreshToken(ctx, userID, req.RefreshToken)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &uproto.AddRefreshTokenResponse{}, fmt.Errorf("failed to add refresh token by there ID")
	}
	return &uproto.AddRefreshTokenResponse{}, nil
}

// GetRefreshTokenByID calls method of Service by handler
func (handl *EntityUser) GetRefreshTokenByID(ctx context.Context, req *uproto.GetRefreshTokenByIDRequest) (*uproto.GetRefreshTokenByIDResponse, error) {
	id := req.Id
	err := handl.validate.VarCtx(ctx, id, "required,uuid")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &uproto.GetRefreshTokenByIDResponse{}, fmt.Errorf("failed to validate")
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &uproto.GetRefreshTokenByIDResponse{}, fmt.Errorf("failed to parse id")
	}
	refreshToken, err := handl.srvcUser.GetRefreshTokenByID(ctx, idUUID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &uproto.GetRefreshTokenByIDResponse{}, fmt.Errorf("failed to get refresh token by there ID")
	}
	return &uproto.GetRefreshTokenByIDResponse{
		RefreshToken: refreshToken,
	}, nil
}

// DeleteAccount calls method of Service by handler
func (handl *EntityUser) DeleteAccount(ctx context.Context, req *uproto.DeleteAccountRequest) (*uproto.DeleteAccountResponse, error) {
	id := req.Id
	err := handl.validate.VarCtx(ctx, id, "required,uuid")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &uproto.DeleteAccountResponse{}, fmt.Errorf("failed to validate")
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &uproto.DeleteAccountResponse{}, fmt.Errorf("failed to parse id")
	}
	err = handl.srvcUser.DeleteAccount(ctx, idUUID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &uproto.DeleteAccountResponse{}, fmt.Errorf("failed to delete by there ID")
	}
	return &uproto.DeleteAccountResponse{
		Id: "Account with ID " + req.Id + " successfully deleted.",
	}, nil
}
