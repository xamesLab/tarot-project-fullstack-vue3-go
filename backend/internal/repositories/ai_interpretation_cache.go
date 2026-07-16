package repositories

import (
	"context"
	"time"

	"tarot-daily/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AIInterpretationCacheRepository struct {
	collection *mongo.Collection
}

func NewAIInterpretationCacheRepository(db *mongo.Database) *AIInterpretationCacheRepository {
	return &AIInterpretationCacheRepository{collection: db.Collection("ai_interpretation_cache")}
}

func (r *AIInterpretationCacheRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "scope", Value: 1},
				{Key: "card_key", Value: 1},
				{Key: "language", Value: 1},
				{Key: "prompt_version", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "status", Value: 1}, {Key: "updated_at", Value: -1}}},
	})
	return err
}

func (r *AIInterpretationCacheRepository) FindCompleted(ctx context.Context, scope, cardKey, language, promptVersion string) (*models.AIInterpretationCache, error) {
	var entry models.AIInterpretationCache
	err := r.collection.FindOne(ctx, bson.M{
		"scope":          scope,
		"card_key":       cardKey,
		"language":       language,
		"prompt_version": promptVersion,
		"status":         models.AIInterpretationCacheStatusCompleted,
	}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (r *AIInterpretationCacheRepository) UpsertCompleted(ctx context.Context, entry *models.AIInterpretationCache) error {
	now := time.Now().UTC()
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = now
	}
	if entry.UpdatedAt.IsZero() {
		entry.UpdatedAt = now
	}

	filter := bson.M{
		"scope":          entry.Scope,
		"card_key":       entry.CardKey,
		"language":       entry.Language,
		"prompt_version": entry.PromptVersion,
	}
	update := bson.M{
		"$set": bson.M{
			"provider":       entry.Provider,
			"model":          entry.Model,
			"status":         entry.Status,
			"interpretation": entry.Interpretation,
			"usage":          entry.Usage,
			"response_data":  entry.ResponseData,
			"updated_at":     entry.UpdatedAt,
			"completed_at":   entry.CompletedAt,
		},
		"$setOnInsert": bson.M{
			"created_at": entry.CreatedAt,
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}
