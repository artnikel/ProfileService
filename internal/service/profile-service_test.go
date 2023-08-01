package service

import (
	"context"
	"crypto/sha256"
	"testing"

	"github.com/artnikel/ProfileService/internal/config"
	"github.com/artnikel/ProfileService/internal/model"
	"github.com/artnikel/ProfileService/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

var (
	// testTokens = &TokenPair{
	// 	AccessToken:  "",
	// 	RefreshToken: "",
	// }

	testUser = model.User{
		ID:           uuid.New(),
		Login:        "testLogin",
		Password:     []byte("testPassword"),
		RefreshToken: "",
	}
	cfg *config.Variables
)

func TestSignUp(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, cfg)
	rep.On("SignUp", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
	err := srv.SignUp(context.Background(), &testUser)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	rep := new(mocks.UserRepository)

	hashedbytes, err := bcrypt.GenerateFromPassword(testUser.Password, bcryptCost)
	require.NoError(t, err)

	rep.On("GetByLogin", mock.Anything, mock.AnythingOfType("string")).
		Return(hashedbytes, testUser.ID, nil)
	rep.On("AddRefreshToken", mock.Anything, mock.AnythingOfType("*model.User")).
		Return(nil)

	srv := NewUserService(rep, cfg)

	_, err = srv.Login(context.Background(), &testUser)
	require.NoError(t, err)
}

func TestRefresh(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, cfg)

	tokenPair, err := srv.GenerateTokenPair(testUser.ID)
	require.NoError(t, err)
	sum := sha256.Sum256([]byte(tokenPair.RefreshToken))

	hashedbytes, err := bcrypt.GenerateFromPassword(sum[:], bcryptCost)
	require.NoError(t, err)

	rep.On("GetRefreshTokenByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(string(hashedbytes), nil)
	rep.On("AddRefreshToken", mock.Anything, mock.AnythingOfType("*model.User")).
		Return(nil)

	_, err = srv.Refresh(context.Background(), tokenPair)
	require.NoError(t, err)
}

func TestDeleteAccount(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, cfg)

	rep.On("DeleteAccount", mock.Anything, mock.AnythingOfType("uuid.UUID")).
		Return(nil)

	err := srv.DeleteAccount(context.Background(), testUser.ID)
	require.NoError(t, err)
}

func TestTokensIDCompare(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, cfg)
	tokenPair, err := srv.GenerateTokenPair(testUser.ID)
	require.NoError(t, err)
	id, err := srv.TokensIDCompare(tokenPair)
	require.NoError(t, err)
	require.Equal(t, testUser.ID, id)
}

func TestHashPassword(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, cfg)
	testBytes := []byte("test")
	_, err := srv.HashPassword(testBytes)
	require.NoError(t, err)
}

func TestCheckPasswordHash(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, cfg)
	testBytes := []byte("test")
	hashedBytes, err := srv.HashPassword(testBytes)
	require.NoError(t, err)
	isEqual, err := srv.CheckPasswordHash(hashedBytes, testBytes)
	require.NoError(t, err)
	require.True(t, isEqual)
}

func TestGenerateTokenPair(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, cfg)
	_, err := srv.GenerateTokenPair(testUser.ID)
	require.NoError(t, err)
}

func TestGenerateJWTToken(t *testing.T) {
	rep := new(mocks.UserRepository)
	srv := NewUserService(rep, cfg)
	_, err := srv.GenerateJWTToken(accessTokenExpiration, testUser.ID)
	require.NoError(t, err)
}
