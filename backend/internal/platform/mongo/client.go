package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"tarotdaily/backend/internal/config"
)

type Client struct {
	Raw *mongo.Client
	DB  *mongo.Database
}

func Connect(ctx context.Context, cfg config.Config) (*Client, error) {
	raw, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, cfg.MongoPingTimeout)
	defer pingCancel()
	if err := raw.Ping(pingCtx, nil); err != nil {
		_ = raw.Disconnect(context.Background())
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &Client{
		Raw: raw,
		DB:  raw.Database(cfg.MongoDBName),
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.Raw.Ping(ctx, nil)
}

func (c *Client) Disconnect(ctx context.Context) error {
	return c.Raw.Disconnect(ctx)
}
