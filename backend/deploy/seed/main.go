package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type sourceCardFile struct {
	List []sourceCard `json:"list"`
}

type sourceCard struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	NameEN           string `json:"name_en"`
	Key              string `json:"key"`
	Arcana           string `json:"arcana"`
	Interpretation   string `json:"interpretation"`
	InterpretationEN string `json:"interpretation_en"`
	Type             string `json:"type"`
	TypeEN           string `json:"type_en"`
	Tags             string `json:"tags"`
	TagsEN           string `json:"tags_en"`
}

func main() {
	_ = godotenv.Load(filepath.Join("..", "..", ".env"))
	_ = godotenv.Load()

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")

	if mongoURI == "" {
		log.Fatal("MONGO_URI is required")
	}
	if dbName == "" {
		dbName = "tarot_daily"
	}

	client, db, err := connectMongo(mongoURI, dbName)
	if err != nil {
		log.Fatalf("connect mongo: %v", err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	cardsPath, err := resolveCardsPath()
	if err != nil {
		log.Fatalf("resolve cards path: %v", err)
	}
	cards, err := loadCards(cardsPath)
	if err != nil {
		log.Fatalf("load cards: %v", err)
	}

	if len(cards) != 78 {
		log.Fatalf("expected 78 cards, got %d", len(cards))
	}

	if err := seedCards(db, cards); err != nil {
		log.Fatalf("seed cards: %v", err)
	}
	if err := seedRolePolicies(db); err != nil {
		log.Fatalf("seed role policies: %v", err)
	}
	if err := seedAdminSettings(db); err != nil {
		log.Fatalf("seed admin settings: %v", err)
	}
	if err := seedDefaultDeck(db); err != nil {
		log.Fatalf("seed default deck: %v", err)
	}
	if err := seedDefaultCardBack(db); err != nil {
		log.Fatalf("seed default card back: %v", err)
	}

	if err := ensureIndexes(db); err != nil {
		log.Fatalf("ensure indexes: %v", err)
	}

	log.Println("seed completed successfully")
}

func connectMongo(uri, dbName string) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, nil, err
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer pingCancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, nil, err
	}

	return client, client.Database(dbName), nil
}

func loadCards(path string) ([]sourceCard, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var file sourceCardFile
	if err := json.Unmarshal(content, &file); err != nil {
		return nil, err
	}

	if len(file.List) == 0 {
		return nil, errors.New("card file is empty")
	}

	return file.List, nil
}

func resolveCardsPath() (string, error) {
	candidates := []string{
		filepath.Join("..", "specs", "temp", "mock", "tarot-cards.json"),
		filepath.Join("..", "..", "specs", "temp", "mock", "tarot-cards.json"),
		filepath.Join("..", "..", "..", "specs", "temp", "mock", "tarot-cards.json"),
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("tarot-cards.json not found in expected locations")
}

func seedCards(db *mongo.Database, cards []sourceCard) error {
	collection := db.Collection("cards")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	docs := make([]any, 0, len(cards))
	now := time.Now().UTC()

	for _, card := range cards {
		docs = append(docs, bson.M{
			"_id":    card.Key,
			"id":     card.ID,
			"key":    card.Key,
			"arcana": card.Arcana,
			"type": bson.M{
				"en": card.TypeEN,
				"ru": card.Type,
			},
			"names": bson.M{
				"en": card.NameEN,
				"ru": card.Name,
			},
			"interpretations": bson.M{
				"en": card.InterpretationEN,
				"ru": card.Interpretation,
			},
			"tags": bson.M{
				"en": card.TagsEN,
				"ru": card.Tags,
			},
			"created_at": now,
			"updated_at": now,
		})
	}

	if err := collection.Drop(ctx); err != nil {
		return fmt.Errorf("drop cards collection: %w", err)
	}

	_, err := collection.InsertMany(ctx, docs)
	return err
}

func seedRolePolicies(db *mongo.Database) error {
	collection := db.Collection("role_policies")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now().UTC()
	docs := []any{
		bson.M{"_id": "seeker", "role": "seeker", "display_name": bson.M{"en": "Seeker", "ru": "Искатель"}, "three_card_limit": 3, "history_retention_days": 1, "created_at": now, "updated_at": now},
		bson.M{"_id": "adept", "role": "adept", "display_name": bson.M{"en": "Adept", "ru": "Адепт"}, "three_card_limit": 5, "history_retention_days": 7, "created_at": now, "updated_at": now},
		bson.M{"_id": "master", "role": "master", "display_name": bson.M{"en": "Master", "ru": "Мастер"}, "three_card_limit": 10, "history_retention_days": 15, "created_at": now, "updated_at": now},
		bson.M{"_id": "keeper", "role": "keeper", "display_name": bson.M{"en": "Keeper", "ru": "Смотритель"}, "three_card_limit": 999999, "history_retention_days": 30, "created_at": now, "updated_at": now},
	}

	if err := collection.Drop(ctx); err != nil {
		return fmt.Errorf("drop role_policies collection: %w", err)
	}

	_, err := collection.InsertMany(ctx, docs)
	return err
}

func seedAdminSettings(db *mongo.Database) error {
	collection := db.Collection("admin_settings")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := collection.Drop(ctx); err != nil {
		return fmt.Errorf("drop admin_settings collection: %w", err)
	}

	_, err := collection.InsertOne(ctx, bson.M{
		"_id":                  "global",
		"maintenance_mode":     false,
		"maintenance_message":  bson.M{"en": "Site under maintenance", "ru": "Сайт на обслуживании"},
		"default_role":         "seeker",
		"default_deck_id":      "default",
		"default_back_id":      "default",
		"registration_enabled": true,
		"created_at":           time.Now().UTC(),
		"updated_at":           time.Now().UTC(),
	})
	return err
}

func seedDefaultDeck(db *mongo.Database) error {
	collection := db.Collection("decks")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := collection.Drop(ctx); err != nil {
		return fmt.Errorf("drop decks collection: %w", err)
	}

	_, err := collection.InsertOne(ctx, bson.M{
		"_id":                  "default",
		"name":                 bson.M{"en": "Tarot Classic", "ru": "Таро Классика"},
		"description":          bson.M{"en": "Default classic tarot deck.", "ru": "Базовая классическая колода таро."},
		"author":               "Pamela Colman Smith",
		"year":                 1909,
		"tags":                 []string{"traditional", "classic"},
		"status":               "done",
		"card_count":           78,
		"back_image_path":      "/uploads/decks/default/back.png",
		"card_image_base_path": "/uploads/decks/default/cards/",
		"created_at":           time.Now().UTC(),
		"updated_at":           time.Now().UTC(),
	})
	return err
}

func seedDefaultCardBack(db *mongo.Database) error {
	collection := db.Collection("card_backs")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := collection.Drop(ctx); err != nil {
		return fmt.Errorf("drop card_backs collection: %w", err)
	}

	_, err := collection.InsertOne(ctx, bson.M{
		"_id":         "default",
		"name":        bson.M{"en": "Classic Back", "ru": "Классическая рубашка"},
		"description": bson.M{"en": "Default classic card back.", "ru": "Базовая классическая рубашка карт."},
		"tags":        []string{"traditional", "classic"},
		"status":      "done",
		"image_path":  "/uploads/decks/default/back.png",
		"created_at":  time.Now().UTC(),
		"updated_at":  time.Now().UTC(),
	})
	return err
}

func ensureIndexes(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if _, err := db.Collection("cards").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "key", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "arcana", Value: 1}}},
	}); err != nil {
		return err
	}

	if _, err := db.Collection("role_policies").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "role", Value: 1}},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return err
	}

	if _, err := db.Collection("decks").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "tags", Value: 1}}},
	}); err != nil {
		return err
	}

	if _, err := db.Collection("card_backs").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "tags", Value: 1}}},
	}); err != nil {
		return err
	}

	return nil
}
