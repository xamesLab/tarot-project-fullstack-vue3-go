package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppPort string
	AppEnv  string

	JWTSecret string
	JWTTTL    time.Duration

	MongoURI    string
	MongoDBName string

	GatewayBaseURL            string
	GatewayRequestPath        string
	GatewayTimeout            time.Duration
	BackendPublicBaseURL      string
	GatewayCallbackHMACSecret string
	GatewayCallbackAllowlist  []string

	HTTPReadHeaderTimeout time.Duration
	HTTPReadTimeout       time.Duration
	HTTPWriteTimeout      time.Duration
	HTTPIdleTimeout       time.Duration
	ShutdownTimeout       time.Duration
	MongoConnectTimeout   time.Duration
	MongoPingTimeout      time.Duration
}

func Load(envFile string) (Config, error) {
	if err := loadEnvFile(envFile); err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppPort:              getEnv("APP_PORT", "8080"),
		AppEnv:               getEnv("APP_ENV", "dev"),
		JWTSecret:            getEnv("JWT_SECRET", ""),
		JWTTTL:               getEnvDuration("JWT_TTL", 24*time.Hour),
		MongoURI:             getEnv("MONGO_URI", ""),
		MongoDBName:          getEnv("MONGO_DB_NAME", ""),
		GatewayBaseURL:       getEnv("GATEWAY_BASE_URL", "https://api.xmlab.ru"),
		GatewayRequestPath:   normalizeURLPath(getEnv("GATEWAY_REQUEST_PATH", "/request")),
		GatewayTimeout:       getEnvDuration("GATEWAY_TIMEOUT", 10*time.Second),
		BackendPublicBaseURL: getEnv("BACKEND_PUBLIC_BASE_URL", ""),
		GatewayCallbackHMACSecret: getEnv(
			"GATEWAY_CALLBACK_HMAC_SECRET",
			"",
		),
		GatewayCallbackAllowlist: getEnvCSV("GATEWAY_CALLBACK_ALLOWLIST"),
		HTTPReadHeaderTimeout:    getEnvDuration("HTTP_READ_HEADER_TIMEOUT", 5*time.Second),
		HTTPReadTimeout:          getEnvDuration("HTTP_READ_TIMEOUT", 15*time.Second),
		HTTPWriteTimeout:         getEnvDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
		HTTPIdleTimeout:          getEnvDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout:          getEnvDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
		MongoConnectTimeout:      getEnvDuration("MONGO_CONNECT_TIMEOUT", 10*time.Second),
		MongoPingTimeout:         getEnvDuration("MONGO_PING_TIMEOUT", 3*time.Second),
	}

	if err := validate(cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func loadEnvFile(path string) error {
	envMap, err := LoadEnvMap(path)
	if err != nil {
		return err
	}

	for key, value := range envMap {
		// Respect explicitly set process env values.
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}

func getEnv(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}

	// Support plain seconds ("15") and duration strings ("15s").
	if digitsOnly(v) {
		sec, err := strconv.Atoi(v)
		if err == nil {
			return time.Duration(sec) * time.Second
		}
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

func getEnvCSV(key string) []string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return nil
	}

	parts := strings.Split(v, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

func normalizeURLPath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return "/request"
	}
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

func digitsOnly(s string) bool {
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return s != ""
}

func validate(cfg Config) error {
	if strings.TrimSpace(cfg.AppPort) == "" {
		return errors.New("APP_PORT is required")
	}
	if strings.TrimSpace(cfg.MongoURI) == "" {
		return errors.New("MONGO_URI is required")
	}
	if strings.TrimSpace(cfg.MongoDBName) == "" {
		return errors.New("MONGO_DB_NAME is required")
	}
	if strings.TrimSpace(cfg.GatewayBaseURL) == "" {
		return errors.New("GATEWAY_BASE_URL is required")
	}
	if strings.TrimSpace(cfg.GatewayRequestPath) == "" {
		return errors.New("GATEWAY_REQUEST_PATH is required")
	}
	if cfg.GatewayTimeout <= 0 {
		return errors.New("GATEWAY_TIMEOUT must be > 0")
	}
	if strings.TrimSpace(cfg.JWTSecret) == "" {
		return errors.New("JWT_SECRET is required")
	}
	if cfg.JWTTTL <= 0 {
		return errors.New("JWT_TTL must be > 0")
	}
	return nil
}
