package main

import (
	"digital-community/internal/config"
	"digital-community/internal/handlers"
	"digital-community/internal/router"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	cfg := config.Load()

	if err := config.InitDB(cfg.DBPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	handlers.StartThumbnailWarmup()

	r := router.Setup()

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
