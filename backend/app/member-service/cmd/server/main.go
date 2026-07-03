package main

import (
	"context"
	"fmt"
	"log"

	"github.com/panduputragit/gym/backend/app/member-service/internal/config"
	servicehttp "github.com/panduputragit/gym/backend/app/member-service/internal/http"
	"github.com/panduputragit/gym/backend/packages/database"
	"github.com/panduputragit/gym/backend/packages/httpserver"
)

func main() {
	cfg := config.Load()
	db, err := database.ConnectOptional(context.Background(), database.Config{URL: cfg.DatabaseURL})
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	if db != nil {
		defer db.Close()
		fmt.Printf("%s connected to database\n", cfg.Name)
	} else {
		fmt.Printf("%s database disabled; set MEMBER_DATABASE_URL to enable it\n", cfg.Name)
	}

	router := httpserver.NewRouter(cfg.Name, cfg.GinMode)
	servicehttp.RegisterRoutes(router)

	addr := ":" + cfg.Port
	fmt.Printf("%s listening on %s\n", cfg.Name, addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("start server: %v", err)
	}
}
