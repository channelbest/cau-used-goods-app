package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"cau-used-goods-app/backend/internal/ai"
	"cau-used-goods-app/backend/internal/auth"
	"cau-used-goods-app/backend/internal/config"
	"cau-used-goods-app/backend/internal/db"
	"cau-used-goods-app/backend/internal/middleware"
	"cau-used-goods-app/backend/internal/product"
	"cau-used-goods-app/backend/internal/sensitive"
	"cau-used-goods-app/backend/internal/stats"
	"cau-used-goods-app/backend/internal/upload"
	"cau-used-goods-app/backend/internal/user"
)

func main() {
	configPath := os.Getenv("APP_CONFIG")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	if err := db.Init(cfg.Database); err != nil {
		log.Fatalf("init database failed: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("close database failed: %v", err)
		}
	}()

	log.Printf("database connected: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	authRepo := auth.NewRepository(db.DB())
	authService := auth.NewService(authRepo, cfg.JWT, cfg.Wechat)
	authHandler := auth.NewHandler(authService)

	userRepo := user.NewRepository(db.DB())
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	sensitiveRepo := sensitive.NewRepository(db.DB())
	sensitiveService := sensitive.NewService(sensitiveRepo)

	productRepo := product.NewRepository(db.DB())
	productService := product.NewService(productRepo, sensitiveService)
	productHandler := product.NewHandler(productService)
	uploadService := upload.NewService()
	uploadHandler := upload.NewHandler(uploadService)
	statsRepo := stats.NewRepository(db.DB())
	statsService := stats.NewService(statsRepo)
	statsHandler := stats.NewHandler(statsService)
	aiService := ai.NewService(cfg.AI.APIKey)
	aiHandler := ai.NewHandler(aiService)
	r := gin.Default()
	r.Static("/uploads", "./uploads")

	authMiddleware := middleware.Auth(cfg.JWT.Secret)

	auth.RegisterRoutes(r, authHandler, authMiddleware, cfg.Server.Env == "dev")
	user.RegisterRoutes(r, userHandler, authMiddleware)
	user.RegisterAdminRoutes(r, userHandler, authMiddleware, middleware.Admin())
	product.RegisterRoutes(r, productHandler, authMiddleware)
	upload.RegisterRoutes(r, uploadHandler, authMiddleware)
	ai.RegisterRoutes(r, aiHandler, authMiddleware)
	stats.RegisterRoutes(r, statsHandler, authMiddleware)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}
