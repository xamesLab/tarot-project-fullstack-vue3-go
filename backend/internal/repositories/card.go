package repositories

import (
	"context"

	"tarot-daily/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
