package repositories

import (
	"context"
	"time"

	"tarot-daily/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeckRepository struct {
	collection *mongo.Collection
}

func NewDeckRepository(db *mongo.Database) *DeckRepository {
	return &DeckRepository{collection: db.Collection("decks")}
}

func (r *DeckRepository) CountAll(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

func (r *DeckRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

func (r *DeckRepository) CountByStatuses(ctx context.Context, statuses []string) (int64, error) {
	filter := bson.M{}
	if len(statuses) > 0 {
		filter["status"] = bson.M{"$in": statuses}
	}

	return r.collection.CountDocuments(ctx, filter)
}

func (r *DeckRepository) Create(ctx context.Context, deck *models.Deck) error {
	_, err := r.collection.InsertOne(ctx, deck)
	return err
}

func (r *DeckRepository) List(ctx context.Context, skip, limit int64) ([]models.Deck, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().
		SetSort(bson.D{{Key: "updated_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var decks []models.Deck
	if err := cursor.All(ctx, &decks); err != nil {
		return nil, err
	}

	return decks, nil
}

func (r *DeckRepository) ListByStatuses(ctx context.Context, statuses []string, skip, limit int64) ([]models.Deck, error) {
	filter := bson.M{}
	if len(statuses) > 0 {
		filter["status"] = bson.M{"$in": statuses}
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().
		SetSort(bson.D{{Key: "updated_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var decks []models.Deck
	if err := cursor.All(ctx, &decks); err != nil {
		return nil, err
	}

	return decks, nil
}

func (r *DeckRepository) ListAvailable(ctx context.Context) ([]models.Deck, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"status": "done"}, options.Find().
		SetSort(bson.D{{Key: "updated_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var decks []models.Deck
	if err := cursor.All(ctx, &decks); err != nil {
		return nil, err
	}

	return decks, nil
}

func (r *DeckRepository) FindByID(ctx context.Context, deckID string) (*models.Deck, error) {
	var deck models.Deck
	if err := r.collection.FindOne(ctx, bson.M{"_id": deckID}).Decode(&deck); err != nil {
		return nil, err
	}

	return &deck, nil
}

func (r *DeckRepository) Update(ctx context.Context, deckID string, updates bson.M) (*models.Deck, error) {
	if updates == nil {
		updates = bson.M{}
	}
	updates["updated_at"] = time.Now().UTC()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var deck models.Deck
	if err := r.collection.FindOneAndUpdate(ctx, bson.M{"_id": deckID}, bson.M{"$set": updates}, opts).Decode(&deck); err != nil {
		return nil, err
	}

	return &deck, nil
}

func (r *DeckRepository) Delete(ctx context.Context, deckID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": deckID})
	return err
}
