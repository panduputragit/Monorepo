package main

import (
	"fmt"
	"log"
	"net"

	"github.com/panduputragit/gym/backend/app/api-gateway/internal/config"
	gatewaygrpc "github.com/panduputragit/gym/backend/app/api-gateway/internal/grpc"
	gatewayhttp "github.com/panduputragit/gym/backend/app/api-gateway/internal/http"
	sharedconfig "github.com/panduputragit/gym/backend/packages/config"
	"github.com/panduputragit/gym/backend/packages/httpserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	sharedconfig.MustInit()
	cfg := config.Load()

	go startHTTP(&cfg)
	startGRPC(&cfg)
}

func startHTTP(cfg *config.Config) {
	router := httpserver.NewRouter(cfg.Name, cfg.GinMode)

	if err := gatewayhttp.RegisterRoutes(router, cfg.ServiceURLs()); err != nil {
		log.Fatalf("register routes: %v", err)
	}

	addr := ":" + cfg.Port
	fmt.Printf("%s HTTP listening on %s\n", cfg.Name, addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("HTTP server: %v", err)
	}
}

func startGRPC(cfg *config.Config) {
	lis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("gRPC listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	if err := gatewaygrpc.RegisterRoutes(grpcServer, cfg.AuthGRPCAddr); err != nil {
		log.Fatalf("register gRPC routes: %v", err)
	}

	fmt.Printf("%s gRPC listening on :%s\n", cfg.Name, cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("gRPC serve: %v", err)
	}
}
