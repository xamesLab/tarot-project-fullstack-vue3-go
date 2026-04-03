package repositories

import (
	"context"
	"sort"

	"tarot-daily/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CardRepository struct {
	collection *mongo.Collection
}

func NewCardRepository(db *mongo.Database) *CardRepository {
	return &CardRepository{collection: db.Collection("cards")}
}

func (r *CardRepository) FindByKey(ctx context.Context, key string) (*models.Card, error) {
	var card models.Card
	if err := r.collection.FindOne(ctx, bson.M{"key": key}).Decode(&card); err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *CardRepository) FindByKeys(ctx context.Context, keys []string) ([]models.Card, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"key": bson.M{"$in": keys}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cards []models.Card
	if err := cursor.All(ctx, &cards); err != nil {
		return nil, err
	}

	return cards, nil
}

func (r *CardRepository) ListAllKeys(ctx context.Context) ([]string, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().
		SetSort(bson.D{{Key: "key", Value: 1}}).
		SetProjection(bson.M{"key": 1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	keys := make([]string, 0, 78)
	for cursor.Next(ctx) {
		var card models.Card
		if err := cursor.Decode(&card); err != nil {
			return nil, err
		}
		keys = append(keys, card.Key)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	sort.Strings(keys)
	return keys, nil
}
