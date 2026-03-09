package http

import (
	"tarotdaily/backend/internal/adminsettings"
	"tarotdaily/backend/internal/adminusers"
	"tarotdaily/backend/internal/auth"
	"tarotdaily/backend/internal/config"
	"tarotdaily/backend/internal/decks"
	"tarotdaily/backend/internal/deckthemes"
	jwtplatform "tarotdaily/backend/internal/platform/jwt"
	mongoplatform "tarotdaily/backend/internal/platform/mongo"
	"tarotdaily/backend/internal/reading"
	mongorepo "tarotdaily/backend/internal/repository/mongo"
	"tarotdaily/backend/internal/roles"
	"tarotdaily/backend/internal/transport/http/handlers"
)

type routeDependencies struct {
	jwtManager           *jwtplatform.Manager
	healthHandler        *handlers.HealthHandler
	authHandler          *handlers.AuthHandler
	meHandler            *handlers.MeHandler
	settingsHandler      *handlers.SettingsHandler
	decksHandler         *handlers.DecksHandler
	deckThemesHandler    *handlers.DeckThemesHandler
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

	jwtManager := jwtplatform.NewManager(cfg.JWTSecret, cfg.JWTTTL)

	adminSettingsService := adminsettings.NewService(adminSettingsRepo)
	authService := auth.NewService(userRepo, deckRepo, deckThemeRepo, jwtManager)
	adminUsersService := adminusers.NewService(userRepo)
	decksService := decks.NewService(deckRepo, userRepo)
	rolesService := roles.NewService(rolePolicyRepo, userRepo)
	readingService := reading.NewService(userRepo, cardRepo, deckRepo, dailyReadingRepo, deckThemeRepo, rolesService)
	deckThemesService := deckthemes.NewService(deckThemeRepo)

	return routeDependencies{
		jwtManager:           jwtManager,
		healthHandler:        handlers.NewHealthHandler(mongoClient.Ping, cfg.MongoPingTimeout),
		authHandler:          handlers.NewAuthHandler(authService),
		meHandler:            handlers.NewMeHandler(authService),
		settingsHandler:      handlers.NewSettingsHandler(authService),
		decksHandler:         handlers.NewDecksHandler(decksService),
		deckThemesHandler:    handlers.NewDeckThemesHandler(deckThemesService),
		readingHandler:       handlers.NewReadingHandler(readingService),
		adminSettingsHandler: handlers.NewAdminSettingsHandler(adminSettingsService),
		accessHandler:        handlers.NewAccessHandler(),
		adminDecksHandler:    handlers.NewAdminDecksHandler(decksService),
		adminUsersHandler:    handlers.NewAdminUsersHandler(adminUsersService),
		adminRolesHandler:    handlers.NewAdminRolesHandler(rolesService),
	}
}
