package handler

// import (
// 	"context"
// 	"testing"

// 	"github.com/artnikel/ProfileService/internal/handler/mocks"
// 	"github.com/artnikel/ProfileService/internal/model"
// 	"github.com/artnikel/ProfileService/internal/service"
// 	"github.com/artnikel/ProfileService/proto"
// 	"github.com/go-playground/validator/v10"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// var (
// 	testTokens = &service.TokenPair{
// 		AccessToken:  "",
// 		RefreshToken: "",
// 	}
// 	testUser = model.User{
// 		ID:           uuid.New(),
// 		Login:        "testLogin",
// 		Password:     []byte("testPassword"),
// 		RefreshToken: "",
// 	}
// 	v = validator.New()
// )

// func TestSignUp(t *testing.T) {
// 	srv := new(mocks.UserService)
// 	hndl := NewEntityUser(srv, v)
// 	protoUser := &proto.User{
// 		Login:    testUser.Login,
// 		Password: string(testUser.Password),
// 	}
// 	srv.On("SignUp", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()
// 	resp, err := hndl.SignUp(context.Background(), &proto.SignUpRequest{
// 		User: protoUser,
// 	})
// 	if resp.Error != "" && resp.Id == uuid.Nil.String() {
// 		t.Errorf("error %v:", resp.Error)
// 	}
// 	assert.NoError(t, err)
// 	srv.AssertExpectations(t)
// }

// func TestLogin(t *testing.T) {
// 	srv := new(mocks.UserService)
// 	hndl := NewEntityUser(srv, v)
// 	protoUser := &proto.User{
// 		Login:    testUser.Login,
// 		Password: string(testUser.Password),
// 	}
// 	srv.On("Login", mock.Anything, mock.AnythingOfType("*model.User")).Return(testTokens, nil).Once()
// 	resp, err := hndl.Login(context.Background(), &proto.LoginRequest{
// 		User: protoUser,
// 	})
// 	if resp.Error != "" {
// 		t.Errorf("error: %v", resp.Error)
// 	}
// 	assert.NoError(t, err)
// 	srv.AssertExpectations(t)
// }

// func TestRefresh(t *testing.T){
// 	srv := new(mocks.UserService)
// 	hndl := NewEntityUser(srv, v)
// 	protoTokens := proto.TokenPair{
// 		AccessToken:  testTokens.AccessToken,
// 		RefreshToken: testTokens.RefreshToken,
// 	}
// 	srv.On("Refresh", mock.Anything, mock.AnythingOfType("service.TokenPair")).Return(testTokens, nil).Once()
// 	resp, err := hndl.Refresh(context.Background(), &proto.RefreshRequest{
// 		Tokenpair: &protoTokens,
// 	})
// 	if resp.Error != "" {
// 		t.Errorf("error: %v", resp.Error)
// 	}
// 	assert.NoError(t, err)
// 	srv.AssertExpectations(t)
// }

// func TestDeleteAccount(t *testing.T){
// 	srv := new(mocks.UserService)
// 	hndl := NewEntityUser(srv, v)
// 	protoID := proto.DeleteAccountRequest{
// 		Id: testUser.ID.String(),
// 	}
// 	srv.On("DeleteAccount", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil)

// 	resp, err := hndl.DeleteAccount(context.Background(), &proto.DeleteAccountRequest{
// 		Id: protoID.Id,
// 	})
// 	if resp.Error != "" {
// 		t.Errorf("error %v:", resp.Error)
// 	}
// 	assert.NoError(t, err)
// 	srv.AssertExpectations(t)
// }
