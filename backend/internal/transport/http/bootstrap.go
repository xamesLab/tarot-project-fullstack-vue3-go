package http

import (
	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/adminsettings"
	"tarotdaily/backend/internal/adminusers"
	"tarotdaily/backend/internal/ask"
	"tarotdaily/backend/internal/auth"
	"tarotdaily/backend/internal/config"
	"tarotdaily/backend/internal/decks"
	"tarotdaily/backend/internal/deckthemes"
	"tarotdaily/backend/internal/gatewayclient"
	jwtplatform "tarotdaily/backend/internal/platform/jwt"
	mongoplatform "tarotdaily/backend/internal/platform/mongo"
	"tarotdaily/backend/internal/reading"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
	"tarotdaily/backend/internal/roles"
	"tarotdaily/backend/internal/transport/http/handlers"
	"tarotdaily/backend/internal/transport/http/middleware"
)

type routeDependencies struct {
	jwtManager           *jwtplatform.Manager
	gatewayCallbackAuth  gin.HandlerFunc
	healthHandler        *handlers.HealthHandler
	gatewayHealthHandler *handlers.GatewayHealthHandler
	authHandler          *handlers.AuthHandler
	meHandler            *handlers.MeHandler
	settingsHandler      *handlers.SettingsHandler
	decksHandler         *handlers.DecksHandler
	deckThemesHandler    *handlers.DeckThemesHandler
	askHandler           *handlers.AskHandler
	readingHandler       *handlers.ReadingHandler
	adminSettingsHandler *handlers.AdminSettingsHandler
	accessHandler        *handlers.AccessHandler
	adminDecksHandler    *handlers.AdminDecksHandler
	adminUsersHandler    *handlers.AdminUsersHandler
	adminRolesHandler    *handlers.AdminRolesHandler
}

func buildRouteDependencies(cfg config.Config, mongoClient *mongoplatform.Client) routeDependencies {
	userRepo := mongorepo.NewUserRepository(mongoClient.DB)
	adminSettingsRepo := mongorepo.NewAdminSettingsRepository(mongoClient.DB)
	deckRepo := mongorepo.NewDeckRepository(mongoClient.DB)
	deckThemeRepo := mongorepo.NewDeckThemeRepository(mongoClient.DB)
	cardRepo := mongorepo.NewCardRepository(mongoClient.DB)
	dailyReadingRepo := mongorepo.NewDailyReadingRepository(mongoClient.DB)
	rolePolicyRepo := mongorepo.NewRolePolicyRepository(mongoClient.DB)
	askRequestRepo := mongorepo.NewAskRequestRepository(mongoClient.DB)

	gatewayClient, err := gatewayclient.New(cfg.GatewayBaseURL, cfg.GatewayRequestPath, cfg.GatewayTimeout)
	if err != nil {
		panic(err)
	}

	jwtManager := jwtplatform.NewManager(cfg.JWTSecret, cfg.JWTTTL)

	adminSettingsService := adminsettings.NewService(adminSettingsRepo)
	authService := auth.NewService(userRepo, deckRepo, deckThemeRepo, jwtManager)
	adminUsersService := adminusers.NewService(userRepo)
	decksService := decks.NewService(deckRepo, userRepo)
	rolesService := roles.NewService(rolePolicyRepo, userRepo)
	askService := ask.NewService(userRepo, askRequestRepo, gatewayClient, cfg.BackendPublicBaseURL)
	readingService := reading.NewService(userRepo, cardRepo, deckRepo, dailyReadingRepo, deckThemeRepo, rolesService)
	deckThemesService := deckthemes.NewService(deckThemeRepo)

	return routeDependencies{
		jwtManager:           jwtManager,
		gatewayCallbackAuth:  middleware.GatewayCallbackVerify(cfg.GatewayCallbackHMACSecret, cfg.GatewayCallbackAllowlist),
		healthHandler:        handlers.NewHealthHandler(mongoClient.Ping, cfg.MongoPingTimeout),
		gatewayHealthHandler: handlers.NewGatewayHealthHandler(gatewayClient),
		authHandler:          handlers.NewAuthHandler(authService),
		meHandler:            handlers.NewMeHandler(authService),
		settingsHandler:      handlers.NewSettingsHandler(authService),
		decksHandler:         handlers.NewDecksHandler(decksService),
		deckThemesHandler:    handlers.NewDeckThemesHandler(deckThemesService),
		askHandler:           handlers.NewAskHandler(askService),
		readingHandler:       handlers.NewReadingHandler(readingService),
		adminSettingsHandler: handlers.NewAdminSettingsHandler(adminSettingsService),
		accessHandler:        handlers.NewAccessHandler(),
		adminDecksHandler:    handlers.NewAdminDecksHandler(decksService),
		adminUsersHandler:    handlers.NewAdminUsersHandler(adminUsersService),
		adminRolesHandler:    handlers.NewAdminRolesHandler(rolesService),
	}
}
