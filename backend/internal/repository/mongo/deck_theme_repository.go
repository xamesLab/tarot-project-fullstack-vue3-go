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

type DeckThemeRepository struct {
	coll *mongo.Collection
}

func NewDeckThemeRepository(db *mongo.Database) *DeckThemeRepository {
	return &DeckThemeRepository{
		coll: db.Collection("deck_themes"),
	}
}

func (r *DeckThemeRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.DeckTheme, error) {
	var theme models.DeckTheme
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&theme)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find deck theme by id: %w", err)
	}
	return &theme, nil
}

func (r *DeckThemeRepository) List(ctx context.Context) ([]models.DeckTheme, error) {
	cur, err := r.coll.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("find deck themes: %w", err)
	}
	defer cur.Close(ctx)

	out := make([]models.DeckTheme, 0)
	for cur.Next(ctx) {
		var item models.DeckTheme
		if err := cur.Decode(&item); err != nil {
			return nil, fmt.Errorf("decode deck theme: %w", err)
		}
		out = append(out, item)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("iterate deck themes: %w", err)
	}

	return out, nil
}

func (r *DeckThemeRepository) Create(
	ctx context.Context,
	name string,
	createdBy primitive.ObjectID,
) (*models.DeckTheme, error) {
	now := time.Now().UTC()
	item := &models.DeckTheme{
		ID:            primitive.NewObjectID(),
		Name:          name,
		IsPlaceholder: true,
		Assets:        map[string]string{},
		CreatedBy:     createdBy,
		CreatedAt:     now,
	}

	_, err := r.coll.InsertOne(ctx, item)
	if err != nil {
		return nil, fmt.Errorf("insert deck theme: %w", err)
	}

	return item, nil
}
