package main

import (
	"context"
	"log"
	"net"

	"github.com/artnikel/ProfileService/internal/config"
	"github.com/artnikel/ProfileService/internal/handler"
	"github.com/artnikel/ProfileService/internal/repository"
	"github.com/artnikel/ProfileService/internal/service"
	"github.com/artnikel/ProfileService/proto"
	"github.com/caarlos0/env"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func connectPostgres(cfg *config.Variables) (*pgxpool.Pool, error) {
	cfgPgx, err := pgxpool.ParseConfig(cfg.PostgresConnProfile)
	if err != nil {
		return nil, err
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfgPgx)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}

func main() {
	var cfg config.Variables
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("could not parse config: ", err)
	}
	validate := validator.New()
	dbpool, errPool := connectPostgres(&cfg)
	if errPool != nil {
		log.Fatal("could not construct the pool: ", errPool)
	}
	defer dbpool.Close()
	pgRep := repository.NewPgRepository(dbpool)
	pgServ := service.NewUserService(pgRep, &cfg)
	pgHandl := handler.NewEntityUser(pgServ, validate)
	lis, err := net.Listen("tcp", "localhost:8090")
	if err != nil {
		log.Fatalf("Cannot create listener: %s", err)
	}
	//interceptorOption := grpc.UnaryInterceptor(interceptor.JWTInterceptor)
	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, pgHandl)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve listener: %s", err)
	}
}
