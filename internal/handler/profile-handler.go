// Package handler is the top level of the application and it contains request handlers
package handler

import (
	"context"

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

// SignUp calls SignUp method of Service by handler
func (handl *EntityUser) SignUp(ctx context.Context, req *proto.SignUpRequest) (*proto.SignUpResponse, error) {
	createdUser := &model.User{
		ID:       uuid.New(),
		Login:    req.User.Login,
		Password: []byte(req.User.Password),
	}

	err := handl.validate.StructCtx(ctx, createdUser)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.SignUpResponse{
			Error: "failed to validate",
		}, nil
	}

	err = handl.srvcUser.SignUp(ctx, createdUser)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.SignUpResponse{
			Error: "failed to signUp",
		}, nil
	}

	return &proto.SignUpResponse{
		Id: createdUser.ID.String(),
	}, nil
}

// Login authenticates user and returns access and refresh tokens
func (handl *EntityUser) GetByLogin(ctx context.Context, req *proto.GetByLoginRequest) (*proto.GetByLoginResponse, error) {
	err := handl.validate.VarCtx(ctx, req.Login, "required,min=3,max=15")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetByLoginResponse{
			Error: "failed to validate",
		}, nil
	}

	password, id, err := handl.srvcUser.GetByLogin(ctx, req.Login)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetByLoginResponse{
			Error: "failed to get password and id by login",
		}, nil
	}

	return &proto.GetByLoginResponse{
		Password: string(password),
		Id:       id.String(),
	}, nil
}

func (handl *EntityUser) AddRefreshToken(ctx context.Context, req *proto.AddRefreshTokenRequest) (*proto.AddRefreshTokenResponse, error) {
	err := handl.validate.VarCtx(ctx, req.Id, "required,uuid")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.AddRefreshTokenResponse{
			Error: "failed to validate",
		}, nil
	}
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.AddRefreshTokenResponse{
			Error: "failed to parse",
		}, nil
	}
	err = handl.srvcUser.AddRefreshToken(ctx, userID, req.RefreshToken)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.AddRefreshTokenResponse{
			Error: "failed to add refresh token by there ID",
		}, nil
	}
	return &proto.AddRefreshTokenResponse{}, nil
}

func (handl *EntityUser) GetRefreshTokenByID(ctx context.Context, req *proto.GetRefreshTokenByIDRequest) (*proto.GetRefreshTokenByIDResponse, error) {
	id := req.Id
	err := handl.validate.VarCtx(ctx, id, "required,uuid")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetRefreshTokenByIDResponse{
			Error: "failed to validate",
		}, nil
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetRefreshTokenByIDResponse{
			Error: "failed to parse id",
		}, nil
	}
	refreshToken, err := handl.srvcUser.GetRefreshTokenByID(ctx, idUUID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetRefreshTokenByIDResponse{
			Error: "failed to get refresh token by there ID",
		}, nil
	}
	return &proto.GetRefreshTokenByIDResponse{
		RefreshToken: refreshToken,
	}, nil
}

func (handl *EntityUser) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.DeleteAccountResponse, error) {
	id := req.Id
	err := handl.validate.VarCtx(ctx, id, "required,uuid")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.DeleteAccountResponse{
			Error: "failed to validate",
		}, nil
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.DeleteAccountResponse{
			Error: "failed to parse id",
		}, nil
	}
	err = handl.srvcUser.DeleteAccount(ctx, idUUID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.DeleteAccountResponse{
			Error: "failed to delete by there ID",
		}, nil
	}
	return &proto.DeleteAccountResponse{
		Id: "Account with ID " + req.Id + " successfully deleted.",
	}, nil
}
