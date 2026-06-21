package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"cau-used-goods-app/backend/internal/admin"
	"cau-used-goods-app/backend/internal/ai"
	"cau-used-goods-app/backend/internal/appeal"
	"cau-used-goods-app/backend/internal/auth"
	"cau-used-goods-app/backend/internal/chat"
	"cau-used-goods-app/backend/internal/config"
	"cau-used-goods-app/backend/internal/db"
	"cau-used-goods-app/backend/internal/favorite"
	"cau-used-goods-app/backend/internal/message"
	"cau-used-goods-app/backend/internal/middleware"
	"cau-used-goods-app/backend/internal/order"
	"cau-used-goods-app/backend/internal/product"
	"cau-used-goods-app/backend/internal/report"
	"cau-used-goods-app/backend/internal/review"
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

	sensitiveRepo := sensitive.NewRepository(db.DB())
	sensitiveService := sensitive.NewService(sensitiveRepo)

	adminRepo := admin.NewRepository(db.DB())
	adminService := admin.NewService(adminRepo)
	adminHandler := admin.NewHandler(adminService)
	sensitiveService.SetAdminLogger(adminService)
	sensitiveHandler := sensitive.NewHandler(sensitiveService)

	messageRepo := message.NewRepository(db.DB())
	messageService := message.NewService(messageRepo)
	messageHandler := message.NewHandler(messageService)

	productRepo := product.NewRepository(db.DB())
	productService := product.NewService(productRepo, sensitiveService, adminService, messageService)
	productHandler := product.NewHandler(productService)

	chatRepo := chat.NewRepository(db.DB())
	chatService := chat.NewService(chatRepo)
	chatHandler := chat.NewHandler(chatService)

	orderRepo := order.NewRepository(db.DB())
	orderService := order.NewService(orderRepo, chatService, adminService)
	orderHandler := order.NewHandler(orderService)

	userRepo := user.NewRepository(db.DB())
	userService := user.NewService(userRepo, productService, orderService, messageService, adminService)
	userHandler := user.NewHandler(userService)

	authRepo := auth.NewRepository(db.DB())
	authService := auth.NewService(authRepo, cfg.JWT, cfg.Wechat, userService)
	authHandler := auth.NewHandler(authService)

	favoriteRepo := favorite.NewRepository(db.DB())
	favoriteService := favorite.NewService(favoriteRepo)
	favoriteHandler := favorite.NewHandler(favoriteService)

	reviewRepo := review.NewRepository(db.DB())
	reviewService := review.NewService(reviewRepo, messageService)
	reviewHandler := review.NewHandler(reviewService)

	reportRepo := report.NewRepository(db.DB())
	reportService := report.NewService(reportRepo, sensitiveService, messageService, adminService)
	reportHandler := report.NewHandler(reportService)

	appealRepo := appeal.NewRepository(db.DB())
	appealService := appeal.NewService(appealRepo, adminService, messageService)
	appealHandler := appeal.NewHandler(appealService)

	uploadService := upload.NewService()
	uploadHandler := upload.NewHandler(uploadService)

	statsRepo := stats.NewRepository(db.DB())
	statsService := stats.NewService(statsRepo)
	statsHandler := stats.NewHandler(statsService)

	aiService := ai.NewService(cfg.AI.APIKey)
	aiHandler := ai.NewHandler(aiService)

	r := gin.Default()
	r.Static("/uploads", "./uploads")

	authMiddleware := middleware.Auth(db.DB(), cfg.JWT.Secret)
	optionalAuthMiddleware := middleware.OptionalAuth(db.DB(), cfg.JWT.Secret)
	readableMiddleware := middleware.ReadableAccount(db.DB())
	normalMiddleware := middleware.NormalAccount(db.DB())
	verifiedMiddleware := middleware.Verified(db.DB())
	adminMiddleware := middleware.Admin(db.DB())
	superAdminMiddleware := middleware.SuperAdmin(db.DB())

	auth.RegisterRoutes(r, authHandler, authMiddleware, cfg.Server.Env == "dev")
	user.RegisterRoutes(r, userHandler, authMiddleware, readableMiddleware, normalMiddleware)
	user.RegisterAdminRoutes(r, userHandler, authMiddleware, adminMiddleware, superAdminMiddleware)
	order.RegisterRoutes(r, orderHandler, authMiddleware, readableMiddleware, verifiedMiddleware, adminMiddleware)
	favorite.RegisterRoutes(r, favoriteHandler, authMiddleware, verifiedMiddleware)
	review.RegisterRoutes(r, reviewHandler, authMiddleware, verifiedMiddleware)
	report.RegisterRoutes(r, reportHandler, authMiddleware, verifiedMiddleware, adminMiddleware)
	message.RegisterRoutes(r, messageHandler, authMiddleware, readableMiddleware)
	chat.RegisterRoutes(r, chatHandler, authMiddleware, readableMiddleware, verifiedMiddleware)
	appeal.RegisterRoutes(r, appealHandler, authMiddleware, readableMiddleware, verifiedMiddleware, adminMiddleware)
	admin.RegisterRoutes(r, adminHandler, authMiddleware, adminMiddleware)
	sensitive.RegisterAdminRoutes(r, sensitiveHandler, authMiddleware, adminMiddleware)

	product.RegisterRoutes(r, productHandler, authMiddleware, optionalAuthMiddleware, verifiedMiddleware, adminMiddleware)
	upload.RegisterRoutes(r, uploadHandler, authMiddleware)
	ai.RegisterRoutes(r, aiHandler, authMiddleware)
	stats.RegisterRoutes(r, statsHandler, authMiddleware, adminMiddleware)
	startOrderCleanupJob(context.Background(), orderService, time.Minute)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}

func startOrderCleanupJob(ctx context.Context, orderService *order.Service, interval time.Duration) {
	if interval <= 0 {
		return
	}
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				count, err := orderService.CancelExpiredOrders(ctx)
				if err != nil {
					log.Printf("cleanup expired orders failed: %v", err)
				}
				if count > 0 {
					log.Printf("cleanup expired orders: canceled %d order(s)", count)
				}
			}
		}
	}()
}
