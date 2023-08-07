package service

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/artnikel/ProfileService/internal/config"
	"github.com/artnikel/ProfileService/internal/model"
	"github.com/artnikel/ProfileService/internal/service/mocks"
	"github.com/caarlos0/env"
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
	cfg config.Variables
)

func TestMain(m *testing.M) {
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("could not parse config: ", err)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestSignUp(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)
	rep.On("SignUp", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
	err := srv.SignUp(context.Background(), &testUser)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestGetByLogin(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)
	rep.On("GetByLogin", mock.Anything, mock.AnythingOfType("string")).Return(testUser.Password, testUser.ID, nil).Once()
	password, id, err := srv.GetByLogin(context.Background(), testUser.Login)
	require.NoError(t, err)
	require.Equal(t, password, testUser.Password)
	require.Equal(t, id, testUser.ID)
	rep.AssertExpectations(t)
}

func TestGetRefreshTokenByID(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)
	rep.On("GetRefreshTokenByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testUser.RefreshToken, nil).Once()
	refreshToken, err := srv.GetRefreshTokenByID(context.Background(), testUser.ID)
	require.NoError(t, err)
	require.Equal(t, refreshToken, testUser.RefreshToken)
	rep.AssertExpectations(t)
}

func TestAddRefreshToken(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)
	rep.On("AddRefreshToken", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("string")).Return(nil).Once()
	err := srv.AddRefreshToken(context.Background(), testUser.ID, testUser.RefreshToken)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestDeleteAccount(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, &cfg)
	rep.On("DeleteAccount", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil)
	err := srv.DeleteAccount(context.Background(), testUser.ID)
	require.NoError(t, err)
}
