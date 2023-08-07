package handler

import (
	"context"
	"testing"

	"github.com/artnikel/ProfileService/internal/handler/mocks"
	"github.com/artnikel/ProfileService/internal/model"
	"github.com/artnikel/ProfileService/uproto"
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
		RefreshToken: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
		eyJleHAiOjE2OTE1MzE2NzAsImlkIjoiMjE5NDkxNjctNTRhOC00NjAwLTk1NzMtM2EwYzAyZTE4NzFjIn0.
		RI9lxDrDlj0RS3FAtNSdwFGz14v9NX1tOxmLjSpZ2dU`,
	}
	v = validator.New()
)

func TestSignUp(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewEntityUser(srv, v)
	protoUser := &uproto.User{
		Login:    testUser.Login,
		Password: string(testUser.Password),
	}
	srv.On("SignUp", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
	_, err := hndl.SignUp(context.Background(), &uproto.SignUpRequest{
		User: protoUser,
	})
	require.NoError(t, err)
	srv.AssertExpectations(t)
}

func TestGetByLogin(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewEntityUser(srv, v)
	protoLogin := uproto.GetByLoginRequest{
		Login: testUser.Login,
	}
	srv.On("GetByLogin", mock.Anything, mock.AnythingOfType("string")).Return(testUser.Password, testUser.ID, nil).Once()
	resp, err := hndl.GetByLogin(context.Background(), &uproto.GetByLoginRequest{
		Login: protoLogin.Login,
	})
	require.NoError(t, err)
	require.Equal(t, resp.Id, testUser.ID.String())
	require.Equal(t, resp.Password, string(testUser.Password))
	srv.AssertExpectations(t)
}

func TestAddRefreshToken(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewEntityUser(srv, v)
	protoID := uproto.AddRefreshTokenRequest{
		Id: testUser.ID.String(),
	}
	srv.On("AddRefreshToken", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("string")).Return(nil).Once()
	_, err := hndl.AddRefreshToken(context.Background(), &uproto.AddRefreshTokenRequest{
		Id: protoID.Id,
	})
	require.NoError(t, err)
	srv.AssertExpectations(t)
}

func TestGetRefreshTokenByID(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewEntityUser(srv, v)
	protoID := uproto.GetRefreshTokenByIDRequest{
		Id: testUser.ID.String(),
	}
	srv.On("GetRefreshTokenByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testUser.RefreshToken, nil).Once()
	resp, err := hndl.GetRefreshTokenByID(context.Background(), &uproto.GetRefreshTokenByIDRequest{
		Id: protoID.Id,
	})
	require.Equal(t, resp.RefreshToken, testUser.RefreshToken)
	require.NoError(t, err)
	srv.AssertExpectations(t)
}

func TestDeleteAccount(t *testing.T) {
	srv := new(mocks.UserService)
	hndl := NewEntityUser(srv, v)
	protoID := uproto.DeleteAccountRequest{
		Id: testUser.ID.String(),
	}
	srv.On("DeleteAccount", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil)

	_, err := hndl.DeleteAccount(context.Background(), &uproto.DeleteAccountRequest{
		Id: protoID.Id,
	})
	require.NoError(t, err)
	srv.AssertExpectations(t)
}
