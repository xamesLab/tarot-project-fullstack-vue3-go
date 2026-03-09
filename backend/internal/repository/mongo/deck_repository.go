package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"tarotdaily/backend/internal/domain/models"
)

type DeckRepository struct {
	coll *mongo.Collection
}

func NewDeckRepository(db *mongo.Database) *DeckRepository {
	return &DeckRepository{
		coll: db.Collection("decks"),
	}
}

func (r *DeckRepository) Create(ctx context.Context, name string, createdByUserID primitive.ObjectID) (*models.Deck, error) {
	now := time.Now().UTC()
	deck := &models.Deck{
		ID:                 primitive.NewObjectID(),
		Name:               name,
		Status:             models.DeckStatusProcess,
		CreatedByUserID:    createdByUserID,
		BackImageURL:       "",
		Cards:              map[string]string{},
		CardsUploadedCount: 0,
		RequiredCardsCount: models.DefaultRequiredCardsCount,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if _, err := r.coll.InsertOne(ctx, deck); err != nil {
		return nil, fmt.Errorf("insert deck: %w", err)
	}
	return deck, nil
}

func (r *DeckRepository) List(ctx context.Context) ([]models.Deck, error) {
	cur, err := r.coll.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, fmt.Errorf("find decks: %w", err)
	}
	defer cur.Close(ctx)

	out := make([]models.Deck, 0)
	for cur.Next(ctx) {
		var deck models.Deck
		if err := cur.Decode(&deck); err != nil {
			return nil, fmt.Errorf("decode deck: %w", err)
		}
		out = append(out, deck)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("iterate decks: %w", err)
	}

	return out, nil
}

func (r *DeckRepository) ListByStatus(ctx context.Context, status models.DeckStatus) ([]models.Deck, error) {
	cur, err := r.coll.Find(
		ctx,
		bson.M{"status": status},
		options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}),
	)
	if err != nil {
		return nil, fmt.Errorf("find decks by status: %w", err)
	}
	defer cur.Close(ctx)

	out := make([]models.Deck, 0)
	for cur.Next(ctx) {
		var deck models.Deck
		if err := cur.Decode(&deck); err != nil {
			return nil, fmt.Errorf("decode deck by status: %w", err)
		}
		out = append(out, deck)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("iterate decks by status: %w", err)
	}

	return out, nil
}

func (r *DeckRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.Deck, error) {
	var deck models.Deck
	if err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&deck); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find deck by id: %w", err)
	}
	return &deck, nil
}

func (r *DeckRepository) UpdateAssets(
	ctx context.Context,
	id primitive.ObjectID,
	backImageURL string,
	cards map[string]string,
	status models.DeckStatus,
	cardsUploadedCount int,
) (*models.Deck, error) {
	setFields := bson.M{
		"backImageUrl":       backImageURL,
		"cards":              cards,
		"status":             status,
		"cardsUploadedCount": cardsUploadedCount,
		"updatedAt":          time.Now().UTC(),
	}

	var deck models.Deck
	if err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": setFields},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&deck); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update deck assets: %w", err)
	}

	return &deck, nil
}
