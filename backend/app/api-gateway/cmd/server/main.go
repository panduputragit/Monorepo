package main

import (
	"fmt"
	"log"

	"github.com/panduputragit/gym/backend/app/api-gateway/internal/config"
	gatewayhttp "github.com/panduputragit/gym/backend/app/api-gateway/internal/http"
	"github.com/panduputragit/gym/backend/packages/httpserver"
)

func main() {
	cfg := config.Load()
	router := httpserver.NewRouter(cfg.Name, cfg.GinMode)

	if err := gatewayhttp.RegisterRoutes(router, cfg.ServiceURLs()); err != nil {
		log.Fatalf("register routes: %v", err)
	}

	addr := ":" + cfg.Port
	fmt.Printf("%s listening on %s\n", cfg.Name, addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
