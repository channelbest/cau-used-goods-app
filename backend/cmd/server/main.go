package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"cau-used-goods-app/backend/internal/auth"
	"cau-used-goods-app/backend/internal/config"
	"cau-used-goods-app/backend/internal/db"
	"cau-used-goods-app/backend/internal/favorite"
	"cau-used-goods-app/backend/internal/middleware"
	"cau-used-goods-app/backend/internal/order"
	"cau-used-goods-app/backend/internal/report"
	"cau-used-goods-app/backend/internal/review"
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

	orderRepo := order.NewRepository(db.DB())
	orderService := order.NewService(orderRepo)
	orderHandler := order.NewHandler(orderService)

	favoriteRepo := favorite.NewRepository(db.DB())
	favoriteService := favorite.NewService(favoriteRepo)
	favoriteHandler := favorite.NewHandler(favoriteService)

	reviewRepo := review.NewRepository(db.DB())
	reviewService := review.NewService(reviewRepo)
	reviewHandler := review.NewHandler(reviewService)

	reportRepo := report.NewRepository(db.DB())
	reportService := report.NewService(reportRepo)
	reportHandler := report.NewHandler(reportService)

	r := gin.Default()
	r.Static("/uploads", "./uploads")

	authMiddleware := middleware.Auth(cfg.JWT.Secret)
	verifiedMiddleware := middleware.Verified(db.DB())

	auth.RegisterRoutes(r, authHandler, authMiddleware, cfg.Server.Env == "dev")
	user.RegisterRoutes(r, userHandler, authMiddleware)
	user.RegisterAdminRoutes(r, userHandler, authMiddleware, middleware.Admin())
	order.RegisterRoutes(r, orderHandler, authMiddleware, verifiedMiddleware)
	favorite.RegisterRoutes(r, favoriteHandler, authMiddleware, verifiedMiddleware)
	review.RegisterRoutes(r, reviewHandler, authMiddleware, verifiedMiddleware)
	report.RegisterRoutes(r, reportHandler, authMiddleware, verifiedMiddleware)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}
