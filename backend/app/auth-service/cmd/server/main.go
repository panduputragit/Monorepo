package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/panduputragit/gym/backend/app/auth-service/internal/config"
	"github.com/panduputragit/gym/backend/app/auth-service/internal/handler"
	authhttp "github.com/panduputragit/gym/backend/app/auth-service/internal/http"
	authgrpc "github.com/panduputragit/gym/backend/app/auth-service/internal/routes/gRPC"
	"github.com/panduputragit/gym/backend/app/auth-service/internal/token"
	sharedconfig "github.com/panduputragit/gym/backend/packages/config"
	"github.com/panduputragit/gym/backend/packages/database"
	"github.com/panduputragit/gym/backend/packages/httpserver"
	"google.golang.org/grpc"
)

func main() {
	sharedconfig.MustInit()
	cfg := config.Load()

	db, err := database.Connect(context.Background(), database.Config{URL: cfg.DatabaseURL})
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()
	fmt.Printf("%s connected to database\n", cfg.Name)

	tokenMaker, err := token.NewMakerWithRandomKey()
	if err != nil {
		log.Fatalf("create token maker: %v", err)
	}

	go startREST(&cfg, db, tokenMaker)
	startGRPC(&cfg, db, tokenMaker)
}

func startREST(cfg *config.Config, db *sql.DB, tokenMaker *token.Maker) {
	router := httpserver.NewRouter(cfg.Name, cfg.GinMode)
	authhttp.RegisterRoutes(router, handler.New(db, tokenMaker))

	addr := ":" + cfg.Port
	fmt.Printf("%s REST listening on %s\n", cfg.Name, addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("REST server: %v", err)
	}
}

func startGRPC(cfg *config.Config, db *sql.DB, tokenMaker *token.Maker) {
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("gRPC listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	if err := authgrpc.RegisterRoutes(grpcServer, db, tokenMaker); err != nil {
		log.Fatalf("register gRPC routes: %v", err)
	}

	fmt.Printf("%s gRPC listening on :%s\n", cfg.Name, cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC serve: %v", err)
	}
}
