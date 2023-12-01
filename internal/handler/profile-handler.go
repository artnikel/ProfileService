// Package handler is the top level of the application and it contains request handlers
package handler

import (
	"context"
	"fmt"

	"github.com/artnikel/ProfileService/internal/model"
	"github.com/artnikel/ProfileService/proto"
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
	proto.UnimplementedUserServiceServer
}

// NewEntityUser accepts User Service interface and returns an object of *EntityUser
func NewEntityUser(srvcUser UserService, validate *validator.Validate) *EntityUser {
	return &EntityUser{srvcUser: srvcUser, validate: validate}
}

// SignUp calls method of Service by handler
func (handl *EntityUser) SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.SignUpResponse, error) {
	createdUser := &model.User{
		ID:       uuid.New(),
		Login:    req.User.Login,
		Password: []byte(req.User.Password),
	}

	err := handl.validate.StructCtx(ctx, createdUser)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.SignUpResponse{}, fmt.Errorf("structCtx %w", err)
	}

	err = handl.srvcUser.SignUp(ctx, createdUser)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.SignUpResponse{}, fmt.Errorf("signUp %w", err)
	}

	return &proto.SignUpResponse{
		Id: createdUser.ID.String(),
	}, nil
}

// GetByLogin calls method of Service by handler
func (handl *EntityUser) GetByLogin(ctx context.Context, req *proto.GetByLoginRequest) (*proto.GetByLoginResponse, error) {
	err := handl.validate.VarCtx(ctx, req.Login, "required,min=5,max=20")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetByLoginResponse{}, fmt.Errorf("varCtx %w", err)
	}

	password, id, err := handl.srvcUser.GetByLogin(ctx, req.Login)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetByLoginResponse{}, fmt.Errorf("getByLogin %w", err)
	}

	return &proto.GetByLoginResponse{
		Password: string(password),
		Id:       id.String(),
	}, nil
}

// AddRefreshToken calls method of Service by handler
func (handl *EntityUser) AddRefreshToken(ctx context.Context, req *proto.AddRefreshTokenRequest) (*proto.AddRefreshTokenResponse, error) {
	err := handl.validate.VarCtx(ctx, req.Id, "required,uuid")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.AddRefreshTokenResponse{}, fmt.Errorf("varCtx %w", err)
	}
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.AddRefreshTokenResponse{}, fmt.Errorf("parse %w", err)
	}
	err = handl.srvcUser.AddRefreshToken(ctx, userID, req.RefreshToken)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.AddRefreshTokenResponse{}, fmt.Errorf("addRefreshToken %w", err)
	}
	return &proto.AddRefreshTokenResponse{}, nil
}

// GetRefreshTokenByID calls method of Service by handler
func (handl *EntityUser) GetRefreshTokenByID(ctx context.Context, req *proto.GetRefreshTokenByIDRequest) (*proto.GetRefreshTokenByIDResponse, error) {
	id := req.Id
	err := handl.validate.VarCtx(ctx, id, "required,uuid")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetRefreshTokenByIDResponse{}, fmt.Errorf("varCtx %w", err)
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetRefreshTokenByIDResponse{}, fmt.Errorf("parse %w", err)
	}
	refreshToken, err := handl.srvcUser.GetRefreshTokenByID(ctx, idUUID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetRefreshTokenByIDResponse{}, fmt.Errorf("getRefreshTokenByID %w", err)
	}
	return &proto.GetRefreshTokenByIDResponse{
		RefreshToken: refreshToken,
	}, nil
}

// DeleteAccount calls method of Service by handler
func (handl *EntityUser) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.DeleteAccountResponse, error) {
	id := req.Id
	err := handl.validate.VarCtx(ctx, id, "required,uuid")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.DeleteAccountResponse{}, fmt.Errorf("varCtx %w", err)
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.DeleteAccountResponse{}, fmt.Errorf("parse %w", err)
	}
	err = handl.srvcUser.DeleteAccount(ctx, idUUID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.DeleteAccountResponse{}, fmt.Errorf("deleteAccount %w", err)
	}
	return &proto.DeleteAccountResponse{
		Id: "Account with ID " + req.Id + " successfully deleted.",
	}, nil
}
