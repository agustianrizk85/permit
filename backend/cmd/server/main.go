package main

import (
	"log"

	"legalpermit/internal/config"
	"legalpermit/internal/database"
	"legalpermit/internal/handler"
	"legalpermit/internal/repository"
	"legalpermit/internal/seed"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	// Seed the two default accounts on first run.
	if err := seed.Accounts(repository.NewUserRepository(db), cfg); err != nil {
		log.Fatalf("seeding failed: %v", err)
	}

	router := handler.NewRouter(db, cfg)

	addr := ":" + cfg.AppPort
	log.Printf("Legal Permit API listening on %s (env=%s)", addr, cfg.AppEnv)
	if err := router.Run(addr); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
