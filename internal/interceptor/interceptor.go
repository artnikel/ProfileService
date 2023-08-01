// Package interceptor is a package for intercept the execution of RPC methods.
package interceptor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/artnikel/ProfileService/internal/config"
	"github.com/caarlos0/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// JWTInterceptor is a gRPC unary server interceptor that performs JWT token validation.
// It checks the authorization header in the metadata of the incoming context and validates the token.
// The interceptor also verifies the token's expiration and performs additional authorization checks for specific methods.
// If the token is valid and the authorization checks pass, the interceptor calls the handler method.
// nolint:gocyclo //JWTInterceptor its a function with a too many checks of token
func JWTInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var cfg config.Variables
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("could not parse config: ", err)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}
	if strings.Contains(info.FullMethod, "/UserService") {
		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "Missing authorization header")
		}
		tokenString := extractTokenFromHeader(authHeader[0])
		if tokenString == "" {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid authorization header format")
		}
		token, err := ValidateToken(tokenString, cfg.TokenSignature)
		if err != nil || !token.Valid {
			logrus.Errorf("%v", err)
			return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			exp := claims["exp"].(float64)
			if exp < float64(time.Now().Unix()) {
				return nil, status.Errorf(codes.Unauthenticated, "Token is expired")
			}
		}
		resp, err := handler(ctx, req)
		if err != nil {
			logrus.Errorf("%v", err)
			return nil, status.Errorf(codes.Unauthenticated, "Failed to run handler method")
		}
		return resp, err
	}
	resp, err := handler(ctx, req)
	if err != nil {
		logrus.Errorf("%v", err)
		return nil, status.Errorf(codes.Unauthenticated, "Failed to run handler method")
	}
	return resp, err
}

// ValidateToken validates the JWT token using the secret key.
// It checks the signing method and returns the parsed token if it is valid.
func ValidateToken(tokenString, secretKey string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func extractTokenFromHeader(authHeader string) string {
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || !strings.EqualFold(strings.ToLower(parts[0]), "bearer") {
		return ""
	}
	return parts[1]
}
