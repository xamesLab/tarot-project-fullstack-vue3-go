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

	validate := utils.NewValidator()
	userRepository := repositories.NewUserRepository(cfg.Database)
	_ = repositories.NewCardRepository(cfg.Database)
	if err := userRepository.EnsureIndexes(context.Background()); err != nil {
		log.Fatalf("ensure user indexes: %v", err)
	}

	authService := services.NewAuthService(userRepository)
	userService := services.NewUserService(userRepository)

	healthHandler := handlers.NewHealthHandler(cfg)
	authHandler := handlers.NewAuthHandler(authService, validate, cfg)
	userHandler := handlers.NewUserHandler(userService, validate)

	api := router.Group("/api")
	healthHandler.Register(api)

	authRoutes := api.Group("/auth")
	authHandler.RegisterRoutes(authRoutes)

	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	userHandler.RegisterRoutes(protected)

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
