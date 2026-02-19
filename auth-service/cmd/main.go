package main

import (
	"auth-service/internal/config"
	"auth-service/internal/handler"

	"auth-service/internal/repository/mysql"
	"auth-service/internal/repository/redis"
	"auth-service/internal/service"
	"auth-service/pkg/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize JWT
	if err := utils.InitJWT(cfg.JWT.PrivateKey, cfg.JWT.PublicKey); err != nil {
		log.Fatalf("Failed to initialize JWT: %v", err)
	}

	// Initialize PostgreSQL
	pgRepo, err := mysql.NewUserRepository(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Initialize Redis
	redisRepo, err := redis.NewTokenRepository(cfg.Redis.Addr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Create service
	authService := service.NewAuthService(pgRepo, redisRepo)

	// Create Gin router
	router := gin.Default()

	// Register handlers
	authHandler := handler.NewAuthHandler(authService)
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	// Start server
	log.Printf("Starting server on port %d", cfg.Server.Port)
	if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
