package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"tarotdaily/backend/internal/domain/models"
)

type CardRepository struct {
	coll *mongo.Collection
}

func NewCardRepository(db *mongo.Database) *CardRepository {
	return &CardRepository{
		coll: db.Collection("cards"),
	}
}

func (r *CardRepository) Sample(ctx context.Context, n int) ([]models.Card, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: n}}}},
	}

	cur, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("sample cards: %w", err)
	}
	defer cur.Close(ctx)

	out := make([]models.Card, 0, n)
	for cur.Next(ctx) {
		var card models.Card
		if err := cur.Decode(&card); err != nil {
			return nil, fmt.Errorf("decode sampled card: %w", err)
		}
		out = append(out, card)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("iterate sampled cards: %w", err)
	}

	return out, nil
}

func (r *CardRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.Card, error) {
	if len(ids) == 0 {
		return []models.Card{}, nil
	}

	cur, err := r.coll.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, fmt.Errorf("find cards by ids: %w", err)
	}
	defer cur.Close(ctx)

	out := make([]models.Card, 0, len(ids))
	for cur.Next(ctx) {
		var card models.Card
		if err := cur.Decode(&card); err != nil {
			return nil, fmt.Errorf("decode card: %w", err)
		}
		out = append(out, card)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("iterate cards by ids: %w", err)
	}
	return out, nil
}
