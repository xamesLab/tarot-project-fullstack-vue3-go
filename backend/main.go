package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tarot-daily/config"
	"tarot-daily/internal/handlers"
	"tarot-daily/internal/middleware"
	"tarot-daily/internal/repositories"
	"tarot-daily/internal/services"
	"tarot-daily/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, cleanup, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	defer cleanup()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), middleware.CORS(cfg.FrontendURL))
	router.MaxMultipartMemory = 256 << 20
	router.Static("/uploads", "./uploads")

	validate := utils.NewValidator()
	userRepository := repositories.NewUserRepository(cfg.Database)
	cardRepository := repositories.NewCardRepository(cfg.Database)
	policyRepository := repositories.NewPolicyRepository(cfg.Database)
	readingRepository := repositories.NewReadingRepository(cfg.Database)
	askRequestRepository := repositories.NewAskRequestRepository(cfg.Database)
	aiInterpretationCacheRepository := repositories.NewAIInterpretationCacheRepository(cfg.Database)
	adminSettingsRepository := repositories.NewAdminSettingsRepository(cfg.Database)
	deckRepository := repositories.NewDeckRepository(cfg.Database)
	cardBackRepository := repositories.NewCardBackRepository(cfg.Database)
	if err := userRepository.EnsureIndexes(context.Background()); err != nil {
		log.Fatalf("ensure user indexes: %v", err)
	}
	if err := policyRepository.EnsureIndexes(context.Background()); err != nil {
		log.Fatalf("ensure policy indexes: %v", err)
	}
	if err := readingRepository.EnsureIndexes(context.Background()); err != nil {
		log.Fatalf("ensure reading indexes: %v", err)
	}
	if err := askRequestRepository.EnsureIndexes(context.Background()); err != nil {
		log.Fatalf("ensure ask request indexes: %v", err)
	}
	if err := aiInterpretationCacheRepository.EnsureIndexes(context.Background()); err != nil {
		log.Fatalf("ensure ai interpretation cache indexes: %v", err)
	}
	if err := adminSettingsRepository.EnsureDefault(context.Background()); err != nil {
		log.Fatalf("ensure admin settings: %v", err)
	}

	authService := services.NewAuthService(userRepository, adminSettingsRepository)
	userService := services.NewUserService(userRepository, deckRepository, cardBackRepository)
	readingService := services.NewReadingService(userRepository, cardRepository, deckRepository, cardBackRepository, readingRepository, policyRepository, adminSettingsRepository, cfg.JWTSecret)
	askService := services.NewAskService(cfg, userRepository, readingRepository, cardRepository, askRequestRepository, aiInterpretationCacheRepository, policyRepository, adminSettingsRepository)
	adminService := services.NewAdminService(userRepository, cardRepository, readingRepository, policyRepository, adminSettingsRepository, deckRepository, cardBackRepository)
	catalogService := services.NewCatalogService(deckRepository, cardBackRepository, adminSettingsRepository)

	healthHandler := handlers.NewHealthHandler(cfg)
	authHandler := handlers.NewAuthHandler(authService, validate, cfg)
	userHandler := handlers.NewUserHandler(userService, validate)
	readingHandler := handlers.NewReadingHandler(readingService, validate)
	askHandler := handlers.NewAskHandler(askService, validate)
	adminHandler := handlers.NewAdminHandler(adminService, validate)
	catalogHandler := handlers.NewCatalogHandler(catalogService)

	api := router.Group("/api")
	healthHandler.Register(api)

	authRoutes := api.Group("/auth")
	authHandler.RegisterRoutes(authRoutes)
	readingHandler.RegisterPublicRoutes(api)
	askHandler.RegisterPublicRoutes(api)

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	userHandler.RegisterRoutes(protected)
	readingHandler.RegisterRoutes(protected)
	askHandler.RegisterRoutes(protected)
	catalogHandler.RegisterRoutes(protected)

	adminRoutes := protected.Group("/admin")
	adminRoutes.Use(middleware.RequireRole("keeper"))
	adminHandler.RegisterRoutes(adminRoutes)

	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTPReadHeaderTimeout,
		ReadTimeout:       cfg.HTTPReadTimeout,
		WriteTimeout:      cfg.HTTPWriteTimeout,
		IdleTimeout:       cfg.HTTPIdleTimeout,
	}

	go func() {
		log.Printf("backend listening on http://localhost:%s", cfg.AppPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}

	log.Println("backend stopped")
	time.Sleep(100 * time.Millisecond)
}
