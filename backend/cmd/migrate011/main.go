package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	configpkg "tarotdaily/backend/internal/config"
)

type config struct {
	Apply            bool
	EnvFile          string
	MongoURI         string
	MongoDB          string
	CardsCatalogPath string
	DefaultThemeName string
	Timeout          time.Duration
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

type cardDoc struct {
	Code                 string `bson:"code"`
	Name                 string `bson:"name"`
	ArcanaType           string `bson:"arcanaType"`
	Description          string `bson:"description"`
	InterpretationOpen   string `bson:"interpretationOpen"`
	InterpretationClosed string `bson:"interpretationClosed"`
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		fail(err)
	}

	cards, assets, validCodes, err := loadCatalog(cfg.CardsCatalogPath)
	if err != nil {
		fail(err)
	}

	if !cfg.Apply {
		fmt.Println("migrate v0.1.1 dry-run:")
		fmt.Printf("- mongo db: %s\n", cfg.MongoDB)
		fmt.Printf("- cards catalog path: %s\n", cfg.CardsCatalogPath)
		fmt.Printf("- cards to upsert: %d\n", len(cards))
		fmt.Printf("- default theme name: %s\n", cfg.DefaultThemeName)
		fmt.Printf("- valid card codes: %d\n", len(validCodes))
		fmt.Println("- apply mode disabled; no DB writes were performed")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		fail(fmt.Errorf("connect mongo: %w", err))
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(ctx, nil); err != nil {
		fail(fmt.Errorf("ping mongo: %w", err))
	}

	db := client.Database(cfg.MongoDB)
	cardsColl := db.Collection("cards")
	dailyReadingsColl := db.Collection("daily_readings")
	deckThemesColl := db.Collection("deck_themes")

	var upsertedCards int64
	for _, card := range cards {
		res, err := cardsColl.UpdateOne(
			ctx,
			bson.M{"code": card.Code},
			bson.M{"$set": card},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			fail(fmt.Errorf("upsert card %s: %w", card.Code, err))
		}
		upsertedCards += res.UpsertedCount
	}

	deleteLegacyRes, err := cardsColl.DeleteMany(ctx, bson.M{"code": bson.M{"$nin": validCodes}})
	if err != nil {
		fail(fmt.Errorf("delete legacy cards: %w", err))
	}

	dailyRes, err := dailyReadingsColl.DeleteMany(ctx, bson.M{})
	if err != nil {
		fail(fmt.Errorf("clear daily_readings: %w", err))
	}

	themeUpdate := bson.M{
		"$set": bson.M{
			"isPlaceholder": true,
			"assets":        assets,
			"updatedAt":     time.Now().UTC(),
		},
		"$setOnInsert": bson.M{
			"name":      cfg.DefaultThemeName,
			"createdAt": time.Now().UTC(),
		},
	}
	themeRes, err := deckThemesColl.UpdateOne(
		ctx,
		bson.M{"name": cfg.DefaultThemeName},
		themeUpdate,
		options.Update().SetUpsert(true),
	)
	if err != nil {
		fail(fmt.Errorf("upsert default deck theme: %w", err))
	}

	cardsCount, err := cardsColl.CountDocuments(ctx, bson.M{})
	if err != nil {
		fail(fmt.Errorf("count cards: %w", err))
	}

	fmt.Println("migrate v0.1.1 applied:")
	fmt.Printf("- cards processed: %d\n", len(cards))
	fmt.Printf("- cards inserted: %d\n", upsertedCards)
	fmt.Printf("- legacy cards removed: %d\n", deleteLegacyRes.DeletedCount)
	fmt.Printf("- daily_readings removed: %d\n", dailyRes.DeletedCount)
	fmt.Printf("- default deck theme upserted: %t\n", themeRes.UpsertedCount > 0)
	fmt.Printf("- cards total in db: %d\n", cardsCount)
}

func parseConfig() (config, error) {
	cfg := config{}
	flag.BoolVar(&cfg.Apply, "apply", false, "apply migration")
	flag.StringVar(&cfg.EnvFile, "env-file", ".env", "path to env file")
	flag.StringVar(&cfg.MongoURI, "mongo-uri", "", "mongo uri")
	flag.StringVar(&cfg.MongoDB, "mongo-db", "", "mongo db")
	flag.StringVar(&cfg.CardsCatalogPath, "cards-catalog", "../specs/temp/mock/tarot-cards.json", "path to cards catalog")
	flag.StringVar(&cfg.DefaultThemeName, "default-theme-name", "Default Classic", "default deck theme name")
	flag.DurationVar(&cfg.Timeout, "timeout", 30*time.Second, "mongo operation timeout")
	flag.Parse()

	env, err := configpkg.LoadEnvMap(cfg.EnvFile)
	if err != nil {
		return cfg, err
	}

	if strings.TrimSpace(cfg.MongoURI) == "" {
		cfg.MongoURI = strings.TrimSpace(env["MONGO_URI"])
	}
	if strings.TrimSpace(cfg.MongoDB) == "" {
		cfg.MongoDB = strings.TrimSpace(env["MONGO_DB_NAME"])
	}
	if strings.TrimSpace(cfg.MongoURI) == "" {
		return cfg, errors.New("mongo-uri is required")
	}
	if strings.TrimSpace(cfg.MongoDB) == "" {
		return cfg, errors.New("mongo-db is required")
	}
	if strings.TrimSpace(cfg.CardsCatalogPath) == "" {
		return cfg, errors.New("cards-catalog is required")
	}
	if strings.TrimSpace(cfg.DefaultThemeName) == "" {
		return cfg, errors.New("default-theme-name is required")
	}

	return cfg, nil
}

func loadCatalog(path string) ([]cardDoc, map[string]string, []string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read cards catalog: %w", err)
	}

	var catalog cardCatalogFile
	if err := json.Unmarshal(raw, &catalog); err != nil {
		return nil, nil, nil, fmt.Errorf("decode cards catalog: %w", err)
	}
	if len(catalog.List) != 78 {
		return nil, nil, nil, fmt.Errorf("invalid cards catalog size: got %d, expected 78", len(catalog.List))
	}

	cards := make([]cardDoc, 0, len(catalog.List))
	validCodes := make([]string, 0, len(catalog.List))
	seenIDs := make(map[int]struct{}, len(catalog.List))
	assets := make(map[string]string, len(catalog.List))

	for _, item := range catalog.List {
		if !isValidCardID(item.ID) {
			return nil, nil, nil, fmt.Errorf("invalid card id: %d", item.ID)
		}
		if _, exists := seenIDs[item.ID]; exists {
			return nil, nil, nil, fmt.Errorf("duplicate card id: %d", item.ID)
		}
		seenIDs[item.ID] = struct{}{}

		code := strconv.Itoa(item.ID)
		name := strings.TrimSpace(item.Name)
		arcana := strings.TrimSpace(item.Arcana)
		interpretation := strings.TrimSpace(item.Interpretation)
		cardType := strings.TrimSpace(item.Type)
		if name == "" || arcana == "" || interpretation == "" {
			return nil, nil, nil, fmt.Errorf("invalid empty fields for card id: %d", item.ID)
		}

		cards = append(cards, cardDoc{
			Code:                 code,
			Name:                 name,
			ArcanaType:           arcana,
			Description:          cardType,
			InterpretationOpen:   interpretation,
			InterpretationClosed: interpretation,
		})
		validCodes = append(validCodes, code)
		assets[code] = fmt.Sprintf("/assets/cards/%s.jpg", code)
	}

	sort.Strings(validCodes)
	sort.Slice(cards, func(i, j int) bool {
		left, _ := strconv.Atoi(cards[i].Code)
		right, _ := strconv.Atoi(cards[j].Code)
		return left < right
	})

	return cards, assets, validCodes, nil
}

func isValidCardID(id int) bool {
	if id >= 22 && id <= 77 {
		return true
	}
	return id >= 100 && id <= 121
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "migrate011 error: %v\n", err)
	os.Exit(1)
}
