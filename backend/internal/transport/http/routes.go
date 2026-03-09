package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/decks"
	"tarotdaily/backend/internal/transport/http/middleware"
	"tarotdaily/backend/internal/transport/http/response"
)

func registerRootRoutes(engine *gin.Engine) {
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "tarot-daily-backend",
			"status":  "ok",
		})
	})
}

func registerStaticRoutes(engine *gin.Engine) {
	engine.Static("/uploads", decks.UploadsRootDir())
	engine.Static("/assets/cards", "../specs/temp/card-image")
}

func registerAPIRoutes(engine *gin.Engine, deps routeDependencies) {
	api := engine.Group("/api")

	registerHealthRoutes(api, deps)
	registerAuthRoutes(api, deps)
	registerProtectedRoutes(api, deps)
}

func registerHealthRoutes(api *gin.RouterGroup, deps routeDependencies) {
	api.GET("/health/live", deps.healthHandler.Live)
	api.GET("/health/ready", deps.healthHandler.Ready)
}

func registerAuthRoutes(api *gin.RouterGroup, deps routeDependencies) {
	authGroup := api.Group("/auth")
	authGroup.POST("/register", deps.authHandler.Register)
	authGroup.POST("/login", deps.authHandler.Login)
}

func registerProtectedRoutes(api *gin.RouterGroup, deps routeDependencies) {
	authProtected := api.Group("")
	authProtected.Use(middleware.AuthRequired(deps.jwtManager))

	authProtected.GET("/me", deps.meHandler.GetMe)
	authProtected.GET("/me/profile", deps.meHandler.GetMe)
	authProtected.PATCH("/me/settings", deps.settingsHandler.PatchMeSettings)
	authProtected.PATCH("/me/profile", deps.settingsHandler.PatchMeSettings)
	authProtected.GET("/decks/available", deps.decksHandler.ListAvailable)
	authProtected.GET("/deck-themes", deps.deckThemesHandler.List)
	authProtected.GET("/reading/today", deps.readingHandler.GetToday)
	authProtected.POST("/reading/generate", deps.readingHandler.GenerateSpreads)

	registerAdminRoutes(authProtected, deps)
}

func registerAdminRoutes(authProtected *gin.RouterGroup, deps routeDependencies) {
	adminGroup := authProtected.Group("/admin")
	adminGroup.Use(middleware.RoleRequired("admin"))

	adminGroup.GET("/access-check", deps.accessHandler.AdminAccess)
	adminGroup.GET("/settings", deps.adminSettingsHandler.Get)
	adminGroup.PATCH("/settings", deps.adminSettingsHandler.Patch)
	adminGroup.GET("/deck-themes", deps.deckThemesHandler.List)
	adminGroup.POST("/deck-themes", deps.deckThemesHandler.Create)
	adminGroup.GET("/decks", deps.adminDecksHandler.List)
	adminGroup.POST("/decks", deps.adminDecksHandler.Create)
	adminGroup.GET("/decks/:deckId", deps.adminDecksHandler.GetByID)
	adminGroup.POST("/decks/:deckId/assets/back", deps.adminDecksHandler.UploadBack)
	adminGroup.POST("/decks/:deckId/assets/cards/:cardKey", deps.adminDecksHandler.UploadCard)
	adminGroup.GET("/roles", deps.adminRolesHandler.List)
	adminGroup.POST("/roles", deps.adminRolesHandler.Upsert)
	adminGroup.GET("/users", deps.adminUsersHandler.List)
	adminGroup.PATCH("/users/role", deps.adminRolesHandler.Assign)
}

func registerNotFound(engine *gin.Engine) {
	engine.NoRoute(func(c *gin.Context) {
		response.JSONError(c, http.StatusNotFound, "not_found", "route not found")
	})
}
