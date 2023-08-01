package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/artnikel/ProfileService/internal/config"
	"github.com/artnikel/ProfileService/internal/interceptor"
	"github.com/artnikel/ProfileService/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository is an interface that contains CRUD methods and GetAll
type UserRepository interface {
	SignUp(ctx context.Context, user *model.User) error
	GetByLogin(ctx context.Context, username string) ([]byte, uuid.UUID, error)
	AddRefreshToken(ctx context.Context, user *model.User) error
	GetRefreshTokenByID(ctx context.Context, id uuid.UUID) (string, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

// UserService contains UserRepository interface
type UserService struct {
	uRep UserRepository
	cfg     *config.Variables
}

// NewUserService accepts UserRepository object and returnes an object of type *UserService
func NewUserService(uRep UserRepository, cfg *config.Variables) *UserService {
	return &UserService{uRep: uRep, cfg: cfg}
}

// Expiration time for an access and a refresh tokens
const (
	accessTokenExpiration  = 15 * time.Minute
	refreshTokenExpiration = 72 * time.Hour
	bcryptCost             = 14
)

// TokenPair contains an Access and a refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// SignUp is a method of UserService that calls  method of Repository
func (us *UserService) SignUp(ctx context.Context, user *model.User) error {
	var err error
	user.Password, err = us.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("UserService -> HashPassword -> error: %w", err)
	}
	err = us.uRep.SignUp(ctx, user)
	if err != nil {
		return fmt.Errorf("UserService -> UserRepository -> SignUp -> error: %w", err)
	}
	return nil
}

// Login is a method of UserService that calls method of Repository
func (us *UserService) Login(ctx context.Context, user *model.User) (*TokenPair, error) {
	hash, id, err := us.uRep.GetByLogin(ctx, user.Login)
	user.ID = id
	if err != nil {
		return &TokenPair{}, fmt.Errorf("UserService ->  Login -> RepositoryUser -> GetPasswordByUsernsame -> error: %w", err)
	}
	verified, err := us.CheckPasswordHash(hash, user.Password)
	if err != nil || !verified {
		return &TokenPair{}, fmt.Errorf("UserService ->  Login -> CheckPasswordHash -> error: %w", err)
	}
	tokenPair, err := us.GenerateTokenPair(user.ID)
	if err != nil {
		return &TokenPair{}, fmt.Errorf("UserService ->  Login -> GenerateTokenPair -> error: %w", err)
	}
	sum := sha256.Sum256([]byte(tokenPair.RefreshToken))
	hashedRefreshToken, err := us.HashPassword(sum[:])
	if err != nil {
		return &TokenPair{}, fmt.Errorf("UserService -> Login -> HashPassword -> error: %w", err)
	}
	user.RefreshToken = string(hashedRefreshToken)
	err = us.uRep.AddRefreshToken(context.Background(), user)
	if err != nil {
		return &TokenPair{}, fmt.Errorf("UserService ->  Login -> RepositoryUser -> AddRefreshToken -> error: %w", err)
	}
	return &tokenPair, nil
}

// Refresh is a method of UserService that refeshes access token and refresh token
func (us *UserService) Refresh(ctx context.Context, tokenPair TokenPair) (*TokenPair, error) {
	id, err := us.TokensIDCompare(tokenPair)
	if err != nil {
		return &TokenPair{}, fmt.Errorf("UserService -> Refresh -> TokensIDCompare -> error: %w", err)
	}
	hash, err := us.uRep.GetRefreshTokenByID(ctx, id)
	if err != nil {
		return &TokenPair{}, fmt.Errorf("UserService ->  Refresh -> RepositoryUser -> GetPasswordByUsernsame -> error: %w", err)
	}
	sum := sha256.Sum256([]byte(tokenPair.RefreshToken))
	verified, err := us.CheckPasswordHash([]byte(hash), sum[:])
	if err != nil || !verified {
		return &TokenPair{}, fmt.Errorf("UserService ->  Refresh -> CheckPasswordHash -> error: refreshToken invalid")
	}
	tokenPair, err = us.GenerateTokenPair(id)
	if err != nil {
		return &TokenPair{}, fmt.Errorf("UserService ->  Refresh -> GenerateTokenPair -> error: %w", err)
	}
	sum = sha256.Sum256([]byte(tokenPair.RefreshToken))
	hashedRefreshToken, err := us.HashPassword(sum[:])
	if err != nil {
		return &TokenPair{}, fmt.Errorf("UserService -> Refresh -> HashPassword -> error: %w", err)
	}
	var user model.User
	user.RefreshToken = string(hashedRefreshToken)
	user.ID = id
	err = us.uRep.AddRefreshToken(context.Background(), &user)
	if err != nil {
		return &TokenPair{}, fmt.Errorf("UserService ->  Refresh -> RepositoryUser -> AddRefreshToken -> error: %w", err)
	}
	return &tokenPair, nil
}

// DeleteAccount is a method from UserService that deleted account by id
func (us *UserService) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	err := us.uRep.DeleteAccount(ctx,id)
	if err != nil {
		return fmt.Errorf("UserService-DeleteAccount: error%w", err)
	}
	return nil
}

// TokensIDCompare compares IDs from refresh and access token for being equal
func (us *UserService) TokensIDCompare(tokenPair TokenPair) (uuid.UUID, error) {
	accessToken, err := interceptor.ValidateToken(tokenPair.AccessToken, us.cfg.TokenSignature)
	if err != nil {
		return uuid.Nil, fmt.Errorf("UserService -> TokensIDCompare -> accessToken -> middleware -> ValidateToken -> error: %w", err)
	}
	var accessID uuid.UUID
	var uuidID uuid.UUID
	if claims, ok := accessToken.Claims.(jwt.MapClaims); ok && accessToken.Valid {
		uuidID, err = uuid.Parse(claims["id"].(string))
		if err != nil {
			return uuid.Nil, fmt.Errorf("UserService -> TokensIDCompare -> accessToken -> uuid.Parse -> error: %w", err)
		}
		accessID = uuidID
	}
	refreshToken, err := interceptor.ValidateToken(tokenPair.RefreshToken, us.cfg.TokenSignature)
	if err != nil {
		return uuid.Nil, fmt.Errorf("UserService -> TokensIDCompare -> refreshToken -> middleware -> ValidateToken -> error: %w", err)
	}
	var refreshID uuid.UUID
	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		exp := claims["exp"].(float64)
		uuidID, err = uuid.Parse(claims["id"].(string))
		if err != nil {
			return uuid.Nil, fmt.Errorf("UserService -> TokensIDCompare -> accessToken  -> uuid.Parse -> error: %w", err)
		}
		refreshID = uuidID
		if exp < float64(time.Now().Unix()) {
			return uuid.Nil, fmt.Errorf("UserService ->  TokensIDCompare -> middleware -> ValidateToken -> error: %w", err)
		}
	}
	if accessID != refreshID {
		return uuid.Nil, fmt.Errorf("user ID in acess token doesn't equal user ID in refresh token")
	}
	return accessID, nil
}

// HashPassword is a method that makes from bytes hashed value
func (us *UserService) HashPassword(password []byte) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword(password, bcryptCost)
	if err != nil {
		return bytes, fmt.Errorf("UserService -> HashPassword -> GenerateFromPassword -> error: %w", err)
	}
	return bytes, nil
}

// CheckPasswordHash is a method  that checks if hash is equal hash from given password
func (us *UserService) CheckPasswordHash(hash, password []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, password)
	if err != nil {
		return false, fmt.Errorf("UserService -> CheckPasswordHash -> CompareHashAndPassword -> error: %w", err)
	}
	return true, nil
}

// GenerateTokenPair generates pair of access and refresh tokens
func (us *UserService) GenerateTokenPair(id uuid.UUID) (TokenPair, error) {
	accessToken, err := us.GenerateJWTToken(accessTokenExpiration, id)
	if err != nil {
		return TokenPair{}, fmt.Errorf("UserService -> GenerateTokenPair -> accessToken -> GenerateJWTToken -> error: %w", err)
	}
	refreshToken, err := us.GenerateJWTToken(refreshTokenExpiration, id)
	if err != nil {
		return TokenPair{}, fmt.Errorf("UserService -> GenerateTokenPair -> refreshToken -> GenerateJWTToken -> error: %w", err)
	}
	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// GenerateJWTToken is a method that generate JWT token with given expiration with user id
func (us *UserService) GenerateJWTToken(expiration time.Duration, id uuid.UUID) (string, error) {
	claims := &jwt.MapClaims{
		"exp": time.Now().Add(expiration).Unix(),
		"id":  id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(us.cfg.TokenSignature))
	if err != nil {
		return "", fmt.Errorf("UserService -> GenerateJWTToken -> token.SignedString -> error: %w", err)
	}
	return tokenString, nil
}
