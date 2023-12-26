package handler

import (
	"context"
	"testing"

	"github.com/artnikel/ProfileService/internal/handler/mocks"
	"github.com/artnikel/ProfileService/internal/model"
	"github.com/artnikel/ProfileService/proto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testUser = model.User{
		ID:       uuid.New(),
		Login:    "testLogin",
		Password: []byte("testPassword"),
	}
	v = validator.New()
)

func TestSignUp(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewEntityUser(srv, v)
	protoUser := &proto.User{
		Login:    testUser.Login,
		Password: string(testUser.Password),
	}
	srv.On("SignUp", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
	_, err := hndl.SignUp(context.Background(), &proto.SignUpRequest{
		User: protoUser,
	})
	require.NoError(t, err)
	srv.AssertExpectations(t)
}

func TestGetByLogin(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewEntityUser(srv, v)
	protoLogin := proto.GetByLoginRequest{
		Login: testUser.Login,
	}
	srv.On("GetByLogin", mock.Anything, mock.AnythingOfType("string")).Return(testUser.Password, testUser.ID, nil).Once()
	resp, err := hndl.GetByLogin(context.Background(), &proto.GetByLoginRequest{
		Login: protoLogin.Login,
	})
	require.NoError(t, err)
	require.Equal(t, resp.Id, testUser.ID.String())
	require.Equal(t, resp.Password, string(testUser.Password))
	srv.AssertExpectations(t)
}

func TestDeleteAccount(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewEntityUser(srv, v)
	protoID := proto.DeleteAccountRequest{
		Id: testUser.ID.String(),
	}
	srv.On("DeleteAccount", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil)

	_, err := hndl.DeleteAccount(context.Background(), &proto.DeleteAccountRequest{
		Id: protoID.Id,
	})
	require.NoError(t, err)
	srv.AssertExpectations(t)
}
