package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	configpkg "tarotdaily/backend/internal/config"
)

type cardSeed struct {
	Code                 string `bson:"code" json:"code"`
	Name                 string `bson:"name" json:"name"`
	ArcanaType           string `bson:"arcanaType" json:"arcanaType"`
	Description          string `bson:"description" json:"description"`
	InterpretationOpen   string `bson:"interpretationOpen" json:"interpretationOpen"`
	InterpretationClosed string `bson:"interpretationClosed" json:"interpretationClosed"`
}

type deckThemeSeed struct {
	Name          string            `bson:"name" json:"name"`
	IsPlaceholder bool              `bson:"isPlaceholder" json:"isPlaceholder"`
	Assets        map[string]string `bson:"assets" json:"assets"`
}

type seedPayload struct {
	GeneratedAt string        `json:"generatedAt"`
	Cards       []cardSeed    `json:"cards"`
	DeckTheme   deckThemeSeed `json:"deckTheme"`
	AdminLogin  string        `json:"adminLogin"`
}

type appConfig struct {
	Apply         bool
	Timeout       time.Duration
	EnvFile       string
	ExportJSON    string
	CardsCatalog  string
	MongoURI      string
	MongoDBName   string
	AdminLogin    string
	AdminPassword string
	DeckName      string
}

type cardCatalogFile struct {
	List []cardCatalogItem `json:"list"`
}

type cardCatalogItem struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Key            string `json:"key"`
	Arcana         string `json:"arcana"`
	Interpretation string `json:"interpretation"`
	Type           string `json:"type"`
}

type adminDoc struct {
	ID primitive.ObjectID `bson:"_id"`
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		fail(err)
	}

	cards, err := buildCards(cfg.CardsCatalog)
	if err != nil {
		fail(err)
	}

	deck := buildDefaultDeckTheme(cfg.DeckName, cards)
	payload := seedPayload{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Cards:       cards,
		DeckTheme:   deck,
		AdminLogin:  cfg.AdminLogin,
	}

	if cfg.ExportJSON != "" {
		if err := exportJSON(cfg.ExportJSON, payload); err != nil {
			fail(err)
		}
		fmt.Printf("seed payload exported: %s\n", cfg.ExportJSON)
	}

	if !cfg.Apply {
		printDryRun(payload)
		return
	}

	if cfg.MongoURI == "" || cfg.MongoDBName == "" {
		fail(errors.New("MONGO_URI and MONGO_DB_NAME are required for --apply mode"))
	}

	if err := applySeeds(cfg, cards, deck); err != nil {
		fail(err)
	}
}

func parseConfig() (appConfig, error) {
	cfg := appConfig{}
	flag.BoolVar(&cfg.Apply, "apply", false, "apply seeds to MongoDB")
	flag.DurationVar(&cfg.Timeout, "timeout", 30*time.Second, "mongo operation timeout")
	flag.StringVar(&cfg.EnvFile, "env-file", ".env", "path to env file")
	flag.StringVar(&cfg.ExportJSON, "export-json", "", "optional path to export generated seed payload")
	flag.StringVar(&cfg.CardsCatalog, "cards-catalog", "../specs/temp/mock/tarot-cards.json", "path to tarot cards catalog JSON")
	flag.StringVar(&cfg.MongoURI, "mongo-uri", "", "override Mongo URI")
	flag.StringVar(&cfg.MongoDBName, "mongo-db", "", "override Mongo database name")
	flag.StringVar(&cfg.AdminLogin, "admin-login", "admin", "admin login for seed user")
	flag.StringVar(&cfg.AdminPassword, "admin-password", "admin12345", "admin password for seed user")
	flag.StringVar(&cfg.DeckName, "deck-name", "Default Classic", "seed deck theme name")
	flag.Parse()

	env, err := configpkg.LoadEnvMap(cfg.EnvFile)
	if err != nil {
		return cfg, err
	}

	if cfg.MongoURI == "" {
		cfg.MongoURI = env["MONGO_URI"]
	}
	if cfg.MongoDBName == "" {
		cfg.MongoDBName = env["MONGO_DB_NAME"]
	}

	if strings.TrimSpace(cfg.AdminLogin) == "" {
		return cfg, errors.New("admin login cannot be empty")
	}
	if len(cfg.AdminPassword) < 8 {
		return cfg, errors.New("admin password must contain at least 8 characters")
	}
	if strings.TrimSpace(cfg.DeckName) == "" {
		return cfg, errors.New("deck name cannot be empty")
	}
	if strings.TrimSpace(cfg.CardsCatalog) == "" {
		return cfg, errors.New("cards catalog path cannot be empty")
	}

	return cfg, nil
}

func buildCards(catalogPath string) ([]cardSeed, error) {
	raw, err := os.ReadFile(catalogPath)
	if err != nil {
		return nil, fmt.Errorf("read cards catalog: %w", err)
	}

	var catalog cardCatalogFile
	if err := json.Unmarshal(raw, &catalog); err != nil {
		return nil, fmt.Errorf("decode cards catalog: %w", err)
	}

	if len(catalog.List) != 78 {
		return nil, fmt.Errorf("invalid cards catalog size: got %d, expected 78", len(catalog.List))
	}

	out := make([]cardSeed, 0, len(catalog.List))
	seenIDs := make(map[int]struct{}, len(catalog.List))
	for _, item := range catalog.List {
		if !isValidCatalogCardID(item.ID) {
			return nil, fmt.Errorf("invalid card id in catalog: %d", item.ID)
		}
		if _, exists := seenIDs[item.ID]; exists {
			return nil, fmt.Errorf("duplicate card id in catalog: %d", item.ID)
		}
		seenIDs[item.ID] = struct{}{}

		name := strings.TrimSpace(item.Name)
		arcana := strings.TrimSpace(item.Arcana)
		interpretation := strings.TrimSpace(item.Interpretation)
		cardType := strings.TrimSpace(item.Type)
		if name == "" || arcana == "" || interpretation == "" {
			return nil, fmt.Errorf("invalid empty fields for card id: %d", item.ID)
		}

		out = append(out, cardSeed{
			Code:                 strconv.Itoa(item.ID),
			Name:                 name,
			ArcanaType:           arcana,
			Description:          cardType,
			InterpretationOpen:   interpretation,
			InterpretationClosed: interpretation,
		})
	}

	sort.Slice(out, func(i, j int) bool {
		left, _ := strconv.Atoi(out[i].Code)
		right, _ := strconv.Atoi(out[j].Code)
		return left < right
	})

	return out, nil
}

func buildDefaultDeckTheme(name string, cards []cardSeed) deckThemeSeed {
	assets := make(map[string]string, len(cards))
	for _, c := range cards {
		assets[c.Code] = fmt.Sprintf("/assets/cards/%s.jpg", c.Code)
	}

	return deckThemeSeed{
		Name:          name,
		IsPlaceholder: true,
		Assets:        assets,
	}
}

func applySeeds(cfg appConfig, cards []cardSeed, deck deckThemeSeed) error {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return fmt.Errorf("connect mongo: %w", err)
	}
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	db := client.Database(cfg.MongoDBName)
	users := db.Collection("users")
	cardsColl := db.Collection("cards")
	deckThemes := db.Collection("deck_themes")

	now := time.Now().UTC()
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	adminFilter := bson.M{"login": cfg.AdminLogin}
	adminUpdate := bson.M{
		"$set": bson.M{
			"passwordHash": string(passwordHash),
			"role":         "admin",
			"status":       "active",
			"displayName":  "Administrator",
			"timezone":     "UTC+00:00",
			"updatedAt":    now,
		},
		"$setOnInsert": bson.M{
			"createdAt": now,
		},
	}

	var admin adminDoc
	err = users.FindOneAndUpdate(
		ctx,
		adminFilter,
		adminUpdate,
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&admin)
	if err != nil {
		return fmt.Errorf("upsert admin user: %w", err)
	}

	upsertedCards := int64(0)
	for _, card := range cards {
		filter := bson.M{"code": card.Code}
		update := bson.M{"$set": card}
		res, err := cardsColl.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return fmt.Errorf("upsert card %s: %w", card.Code, err)
		}
		upsertedCards += res.UpsertedCount
	}

	deckFilter := bson.M{"name": deck.Name}
	deckUpdate := bson.M{
		"$set": bson.M{
			"isPlaceholder": deck.IsPlaceholder,
			"assets":        deck.Assets,
			"createdBy":     admin.ID,
		},
		"$setOnInsert": bson.M{
			"createdAt": now,
		},
	}

	deckRes, err := deckThemes.UpdateOne(ctx, deckFilter, deckUpdate, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("upsert deck theme: %w", err)
	}

	fmt.Println("seed applied successfully:")
	fmt.Printf("- cards total processed: %d\n", len(cards))
	fmt.Printf("- cards inserted: %d\n", upsertedCards)
	fmt.Printf("- admin login: %s\n", cfg.AdminLogin)
	fmt.Printf("- deck theme: %s (upserted=%t)\n", deck.Name, deckRes.UpsertedCount > 0)

	return nil
}

func exportJSON(path string, payload seedPayload) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("export path is empty")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create export dir: %w", err)
	}

	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write payload file: %w", err)
	}
	return nil
}

func printDryRun(payload seedPayload) {
	fmt.Println("seed dry-run (safe mode):")
	fmt.Printf("- generatedAt: %s\n", payload.GeneratedAt)
	fmt.Printf("- cards: %d\n", len(payload.Cards))
	fmt.Printf("- deck theme: %s\n", payload.DeckTheme.Name)
	fmt.Printf("- admin login: %s\n", payload.AdminLogin)
	fmt.Println("- apply mode disabled; no DB writes were performed")
}

func isValidCatalogCardID(id int) bool {
	if id >= 22 && id <= 77 {
		return true
	}
	return id >= 100 && id <= 121
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "seed error: %v\n", err)
	os.Exit(1)
}
