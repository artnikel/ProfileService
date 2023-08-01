package handler

import (
	"context"

	"github.com/artnikel/ProfileService/internal/model"
	"github.com/artnikel/ProfileService/internal/service"
	"github.com/artnikel/ProfileService/proto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// UserService is an interface that contains methods of service for user
type UserService interface {
	SignUp(ctx context.Context, user *model.User) error
	Login(ctx context.Context, user *model.User) (*service.TokenPair, error)
	Refresh(ctx context.Context, tokenPair service.TokenPair) (*service.TokenPair, error)
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
func (handl *EntityUser) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	loginedUser := &model.User{
		Login:    req.User.Login,
		Password: []byte(req.User.Password),
	}

	err := handl.validate.StructCtx(ctx, loginedUser)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.LoginResponse{
			Error: "failed to validate",
		}, nil
	}

	tokenPair, err := handl.srvcUser.Login(ctx, loginedUser)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.LoginResponse{
			Error: "failed to login",
		}, nil
	}

	return &proto.LoginResponse{
		Tokenpair: &proto.TokenPair{
			AccessToken:  tokenPair.AccessToken,
			RefreshToken: tokenPair.RefreshToken,
		},
	}, nil
}

// Refresh refreshes pair of access and refresh tokens
func (handl *EntityUser) Refresh(ctx context.Context, req *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	tokenPair := service.TokenPair{
		AccessToken:  req.Tokenpair.AccessToken,
		RefreshToken: req.Tokenpair.RefreshToken,
	}
	refreshedTokenPair, err := handl.srvcUser.Refresh(ctx, tokenPair)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.RefreshResponse{
			Error: "failed to refresh",
		}, nil
	}

	return &proto.RefreshResponse{
		Tokenpair: &proto.TokenPair{
			AccessToken:  refreshedTokenPair.AccessToken,
			RefreshToken: refreshedTokenPair.RefreshToken,
		},
	}, nil
}

func (handl *EntityUser) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.DeleteAccountResponse, error){
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
			Error: "failed to read by there ID",
		}, nil
	}
	return &proto.DeleteAccountResponse{}, nil
}
