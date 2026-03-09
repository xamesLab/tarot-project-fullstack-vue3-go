package http

import (
	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/config"
	mongoplatform "tarotdaily/backend/internal/platform/mongo"
	"tarotdaily/backend/internal/transport/http/middleware"
)

func NewRouter(cfg config.Config, mongoClient *mongoplatform.Client) *gin.Engine {
	if cfg.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(middleware.RequestID())
	engine.Use(middleware.Recovery())
	engine.Use(middleware.ErrorHandler())

	deps := buildRouteDependencies(cfg, mongoClient)

	registerRootRoutes(engine)
	registerStaticRoutes(engine)
	registerAPIRoutes(engine, deps)
	registerNotFound(engine)

	return engine
}
