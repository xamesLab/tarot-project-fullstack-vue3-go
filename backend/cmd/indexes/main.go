package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	configpkg "tarotdaily/backend/internal/config"
)

type config struct {
	EnvFile   string
	MongoURI  string
	MongoDB   string
	Timeout   time.Duration
	ShowStats bool
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		fail(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		fail(fmt.Errorf("mongo connect: %w", err))
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		fail(fmt.Errorf("mongo ping: %w", err))
	}

	db := client.Database(cfg.MongoDB)
	if err := ensureIndexes(ctx, db); err != nil {
		fail(err)
	}

	fmt.Printf("indexes ensured for database: %s\n", cfg.MongoDB)
	if cfg.ShowStats {
		if err := printStats(ctx, db); err != nil {
			fail(err)
		}
	}
}

func parseConfig() (config, error) {
	cfg := config{}
	flag.StringVar(&cfg.EnvFile, "env-file", ".env", "path to env file")
	flag.StringVar(&cfg.MongoURI, "mongo-uri", "", "mongodb uri")
	flag.StringVar(&cfg.MongoDB, "mongo-db", "", "mongodb database name")
	flag.DurationVar(&cfg.Timeout, "timeout", 20*time.Second, "operation timeout")
	flag.BoolVar(&cfg.ShowStats, "stats", true, "print collection stats after index creation")
	flag.Parse()

	env, err := configpkg.LoadEnvMap(cfg.EnvFile)
	if err != nil {
		return cfg, err
	}

	if cfg.MongoURI == "" {
		cfg.MongoURI = env["MONGO_URI"]
	}
	if cfg.MongoDB == "" {
		cfg.MongoDB = env["MONGO_DB_NAME"]
	}
	if strings.TrimSpace(cfg.MongoURI) == "" {
		return cfg, errors.New("mongo-uri is required")
	}
	if strings.TrimSpace(cfg.MongoDB) == "" {
		return cfg, errors.New("mongo-db is required")
	}
	return cfg, nil
}

func ensureIndexes(ctx context.Context, db *mongo.Database) error {
	users := db.Collection("users")
	dailyReadings := db.Collection("daily_readings")
	cards := db.Collection("cards")
	deckThemes := db.Collection("deck_themes")
	decks := db.Collection("decks")
	rolePolicies := db.Collection("role_policies")

	_, err := users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "login", Value: 1}},
		Options: options.Index().SetName("uniq_users_login").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create index users.login: %w", err)
	}

	_, err = users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "role", Value: 1}, {Key: "status", Value: 1}, {Key: "createdAt", Value: -1}},
		Options: options.Index().SetName("idx_users_role_status_created_at"),
	})
	if err != nil {
		return fmt.Errorf("create index users(role,status,createdAt): %w", err)
	}

	_, err = users.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "displayName", Value: 1}},
		Options: options.Index().SetName("idx_users_display_name"),
	})
	if err != nil {
		return fmt.Errorf("create index users.displayName: %w", err)
	}

	_, err = dailyReadings.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "userId", Value: 1}, {Key: "date", Value: 1}},
		Options: options.Index().SetName("uniq_daily_readings_user_date").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create index daily_readings(userId,date): %w", err)
	}

	_, err = cards.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "code", Value: 1}},
		Options: options.Index().SetName("uniq_cards_code").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create index cards.code: %w", err)
	}

	_, err = deckThemes.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetName("uniq_deck_themes_name").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create index deck_themes.name: %w", err)
	}

	_, err = decks.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetName("uniq_decks_name").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create index decks.name: %w", err)
	}

	_, err = decks.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "status", Value: 1}, {Key: "createdAt", Value: -1}},
		Options: options.Index().SetName("idx_decks_status_created_at"),
	})
	if err != nil {
		return fmt.Errorf("create index decks(status,createdAt): %w", err)
	}

	_, err = rolePolicies.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "role", Value: 1}},
		Options: options.Index().SetName("uniq_role_policies_role").SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("create index role_policies.role: %w", err)
	}

	return nil
}

func printStats(ctx context.Context, db *mongo.Database) error {
	collections := []string{"users", "cards", "deck_themes", "decks", "daily_readings", "role_policies"}
	for _, name := range collections {
		coll := db.Collection(name)

		count, err := coll.CountDocuments(ctx, bson.D{})
		if err != nil {
			return fmt.Errorf("count %s: %w", name, err)
		}

		fmt.Printf("- %s: %d docs\n", name, count)
	}
	return nil
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "indexes error: %v\n", err)
	os.Exit(1)
}
