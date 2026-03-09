package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	mongoURI string
	dbName   string
	appUser  string
	appPass  string
	timeout  time.Duration
}

func main() {
	cfg := parseFlags()
	if cfg.appPass == "" {
		fail("app password is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.mongoURI))
	if err != nil {
		failf("mongo connect failed: %v", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		failf("mongo ping failed: %v", err)
	}

	db := client.Database(cfg.dbName)

	// Ensure database exists (Mongo creates DB lazily on first write).
	if _, err := db.Collection("__bootstrap").InsertOne(ctx, bson.M{"createdAt": time.Now().UTC()}); err != nil {
		failf("bootstrap write failed: %v", err)
	}

	userExists, err := checkUserExists(ctx, db, cfg.appUser)
	if err != nil {
		failf("check user failed: %v", err)
	}

	if userExists {
		if err := updateUser(ctx, db, cfg.appUser, cfg.appPass, cfg.dbName); err != nil {
			failf("update user failed: %v", err)
		}
		fmt.Printf("updated app user: %s\n", cfg.appUser)
	} else {
		if err := createUser(ctx, db, cfg.appUser, cfg.appPass, cfg.dbName); err != nil {
			failf("create user failed: %v", err)
		}
		fmt.Printf("created app user: %s\n", cfg.appUser)
	}

	fmt.Printf("database ready: %s\n", cfg.dbName)
	fmt.Printf("connection string for app:\n")
	fmt.Printf("mongodb://%s:%s@127.0.0.1:27017/%s?authSource=%s\n", cfg.appUser, cfg.appPass, cfg.dbName, cfg.dbName)
}

func parseFlags() config {
	cfg := config{}
	flag.StringVar(&cfg.mongoURI, "mongo-uri", "mongodb://admin:331707@127.0.0.1:27017/?authSource=admin", "admin mongo uri")
	flag.StringVar(&cfg.dbName, "db-name", "tarot_daily", "project database name")
	flag.StringVar(&cfg.appUser, "app-user", "tarot_daily_app", "application mongo user")
	flag.StringVar(&cfg.appPass, "app-pass", "", "application mongo password")
	flag.DurationVar(&cfg.timeout, "timeout", 20*time.Second, "operation timeout")
	flag.Parse()
	return cfg
}

func checkUserExists(ctx context.Context, db *mongo.Database, username string) (bool, error) {
	var out bson.M
	cmd := bson.D{{Key: "usersInfo", Value: username}}
	if err := db.RunCommand(ctx, cmd).Decode(&out); err != nil {
		return false, err
	}

	rawUsers, ok := out["users"].(bson.A)
	if !ok {
		return false, nil
	}
	return len(rawUsers) > 0, nil
}

func createUser(ctx context.Context, db *mongo.Database, username, password, dbName string) error {
	cmd := bson.D{
		{Key: "createUser", Value: username},
		{Key: "pwd", Value: password},
		{Key: "roles", Value: bson.A{
			bson.D{{Key: "role", Value: "readWrite"}, {Key: "db", Value: dbName}},
		}},
	}
	return db.RunCommand(ctx, cmd).Err()
}

func updateUser(ctx context.Context, db *mongo.Database, username, password, dbName string) error {
	cmd := bson.D{
		{Key: "updateUser", Value: username},
		{Key: "pwd", Value: password},
		{Key: "roles", Value: bson.A{
			bson.D{{Key: "role", Value: "readWrite"}, {Key: "db", Value: dbName}},
		}},
	}
	return db.RunCommand(ctx, cmd).Err()
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, "dbinit error:", msg)
	os.Exit(1)
}

func failf(format string, args ...any) {
	fail(fmt.Sprintf(format, args...))
}
