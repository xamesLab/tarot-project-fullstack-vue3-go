package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	AppEnv                string
	AppPort               string
	AIDirectEnabled       bool
	JWTSecret             string
	JWTTTL                string
	FrontendURL           string
	BackendPublicBaseURL  string
	MongoURI              string
	MongoDBName           string
	MongoConnectTimeout   time.Duration
	MongoPingTimeout      time.Duration
	HTTPReadHeaderTimeout time.Duration
	HTTPReadTimeout       time.Duration
	HTTPWriteTimeout      time.Duration
	HTTPIdleTimeout       time.Duration
	ShutdownTimeout       time.Duration
	GatewayBaseURL        string
	GatewayRequestPath    string
	GatewayTimeout        time.Duration
	GatewayCallbackSecret string
	GatewayCallbackAllow  string
	OpenRouterBaseURL     string
	OpenRouterKey         string
	OpenRouterModelMain   string
	OpenRouterReferer     string
	OpenRouterTitle       string
	OpenRouterTimeout     time.Duration
	MongoClient           *mongo.Client
	Database              *mongo.Database
}

func (c *Config) connectMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), c.MongoConnectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(c.MongoURI))
	if err != nil {
		return fmt.Errorf("connect mongo: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), c.MongoPingTimeout)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return fmt.Errorf("ping mongo: %w", err)
	}

	c.MongoClient = client
	c.Database = client.Database(c.MongoDBName)
	return nil
}

func (c *Config) disconnectMongo() {
	if c.MongoClient == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = c.MongoClient.Disconnect(ctx)
}
