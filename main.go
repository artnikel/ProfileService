// Package main of a project
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/artnikel/ProfileService/internal/config"
	"github.com/artnikel/ProfileService/internal/handler"
	"github.com/artnikel/ProfileService/internal/repository"
	"github.com/artnikel/ProfileService/internal/service"
	"github.com/artnikel/ProfileService/proto"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func connectPostgres(connString string) (*pgxpool.Pool, error) {
	cfgPostgres, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfgPostgres)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}

// nolint gocritic
func main() {
	v := validator.New()
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("could not parse config: %v", err)
	}
	dbpool, errPool := connectPostgres(cfg.PostgresConnProfile)
	if errPool != nil {
		log.Fatalf("could not construct the pool: %v", errPool)
	}
	defer dbpool.Close()
	pgRep := repository.NewPgRepository(dbpool)
	pgServ := service.NewUserService(pgRep)
	pgHandl := handler.NewEntityUser(pgServ, v)
	lis, err := net.Listen("tcp", "localhost:8090")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	fmt.Println("Profile Server started")
	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, pgHandl)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve listener: %s", err)
	}
}
