package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func Load() (*Config, func(), error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:                getEnv("APP_ENV", "dev"),
		AppPort:               getEnv("APP_PORT", "8888"),
		AIDirectEnabled:       getEnvBool("AI_DIRECT_ENABLED", true),
		JWTSecret:             getEnv("JWT_SECRET", "dev-jwt-secret-local-only"),
		JWTTTL:                getEnv("JWT_TTL", "24h"),
		FrontendURL:           getEnv("FRONTEND_URL", "http://localhost:5173"),
		BackendPublicBaseURL:  getEnv("BACKEND_PUBLIC_BASE_URL", "http://localhost:8888"),
		MongoURI:              getEnv("MONGO_URI", ""),
		MongoDBName:           getEnv("MONGO_DB_NAME", "tarot_daily"),
		GatewayBaseURL:        getEnv("GATEWAY_BASE_URL", ""),
		GatewayRequestPath:    getEnv("GATEWAY_REQUEST_PATH", ""),
		GatewayCallbackSecret: getEnv("GATEWAY_CALLBACK_HMAC_SECRET", ""),
		GatewayCallbackAllow:  getEnv("GATEWAY_CALLBACK_ALLOWLIST", ""),
		OpenRouterBaseURL:     getEnv("OPENROUTER_BASE_URL", "https://openrouter.ai/api/v1"),
		OpenRouterKey:         getEnv("OPENROUTER_KEY", ""),
		OpenRouterModelMain:   getEnv("OPENROUTER_MODEL_MAIN", ""),
		OpenRouterReferer:     getEnv("OPENROUTER_HTTP_REFERER", ""),
		OpenRouterTitle:       getEnv("OPENROUTER_APP_TITLE", ""),
	}

	if cfg.MongoURI == "" {
		return nil, nil, fmt.Errorf("missing required environment variable: MONGO_URI")
	}

	var err error

	if cfg.MongoConnectTimeout, err = time.ParseDuration(getEnv("MONGO_CONNECT_TIMEOUT", "10s")); err != nil {
		return nil, nil, fmt.Errorf("parse MONGO_CONNECT_TIMEOUT: %w", err)
	}
	if cfg.MongoPingTimeout, err = time.ParseDuration(getEnv("MONGO_PING_TIMEOUT", "3s")); err != nil {
		return nil, nil, fmt.Errorf("parse MONGO_PING_TIMEOUT: %w", err)
	}
	if cfg.HTTPReadHeaderTimeout, err = time.ParseDuration(getEnv("HTTP_READ_HEADER_TIMEOUT", "5s")); err != nil {
		return nil, nil, fmt.Errorf("parse HTTP_READ_HEADER_TIMEOUT: %w", err)
	}
	if cfg.HTTPReadTimeout, err = time.ParseDuration(getEnv("HTTP_READ_TIMEOUT", "15s")); err != nil {
		return nil, nil, fmt.Errorf("parse HTTP_READ_TIMEOUT: %w", err)
	}
	if cfg.HTTPWriteTimeout, err = time.ParseDuration(getEnv("HTTP_WRITE_TIMEOUT", "15s")); err != nil {
		return nil, nil, fmt.Errorf("parse HTTP_WRITE_TIMEOUT: %w", err)
	}
	if cfg.HTTPIdleTimeout, err = time.ParseDuration(getEnv("HTTP_IDLE_TIMEOUT", "60s")); err != nil {
		return nil, nil, fmt.Errorf("parse HTTP_IDLE_TIMEOUT: %w", err)
	}
	if cfg.ShutdownTimeout, err = time.ParseDuration(getEnv("SHUTDOWN_TIMEOUT", "10s")); err != nil {
		return nil, nil, fmt.Errorf("parse SHUTDOWN_TIMEOUT: %w", err)
	}
	if cfg.GatewayTimeout, err = time.ParseDuration(getEnv("GATEWAY_TIMEOUT", "10s")); err != nil {
		return nil, nil, fmt.Errorf("parse GATEWAY_TIMEOUT: %w", err)
	}
	if cfg.OpenRouterTimeout, err = time.ParseDuration(getEnv("OPENROUTER_TIMEOUT", "20s")); err != nil {
		return nil, nil, fmt.Errorf("parse OPENROUTER_TIMEOUT: %w", err)
	}

	if err := cfg.connectMongo(); err != nil {
		return nil, nil, err
	}

	return cfg, cfg.disconnectMongo, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		return fallback
	}

	switch value {
	case "1", "true", "TRUE", "yes", "YES", "on", "ON":
		return true
	case "0", "false", "FALSE", "no", "NO", "off", "OFF":
		return false
	default:
		return fallback
	}
}
