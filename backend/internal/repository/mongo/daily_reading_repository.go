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

type DailyReadingRepository struct {
	coll *mongo.Collection
}

func NewDailyReadingRepository(db *mongo.Database) *DailyReadingRepository {
	return &DailyReadingRepository{
		coll: db.Collection("daily_readings"),
	}
}

func (r *DailyReadingRepository) FindByUserAndDate(
	ctx context.Context,
	userID primitive.ObjectID,
	date string,
) (*models.DailyReading, error) {
	var item models.DailyReading
	err := r.coll.FindOne(ctx, bson.M{"userId": userID, "date": date}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find daily reading: %w", err)
	}
	return &item, nil
}

func (r *DailyReadingRepository) Create(ctx context.Context, item *models.DailyReading) error {
	if item == nil {
		return errors.New("daily reading is nil")
	}
	_, err := r.coll.InsertOne(ctx, item)
	if err != nil {
		return fmt.Errorf("insert daily reading: %w", err)
	}
	return nil
}

func (r *DailyReadingRepository) UpdateSpreads(
	ctx context.Context,
	id primitive.ObjectID,
	spreadA []models.CardDraw,
	spreadB []models.CardDraw,
	deckThemeID *primitive.ObjectID,
) (*models.DailyReading, error) {
	now := time.Now().UTC()
	setFields := bson.M{
		"spreadA":               spreadA,
		"spreadB":               spreadB,
		"lastSpreadGeneratedAt": now,
		"updatedAt":             now,
	}

	if deckThemeID != nil {
		setFields["deckThemeId"] = *deckThemeID
	}

	update := bson.M{
		"$set": setFields,
		"$inc": bson.M{"spreadGenerationsCount": 1},
	}
	if deckThemeID == nil {
		update["$unset"] = bson.M{"deckThemeId": ""}
	}

	var out models.DailyReading
	err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update daily reading spreads: %w", err)
	}

	return &out, nil
}

func (r *DailyReadingRepository) UpdateCardOfDay(
	ctx context.Context,
	id primitive.ObjectID,
	cardOfDay models.CardDraw,
	deckThemeID *primitive.ObjectID,
) (*models.DailyReading, error) {
	now := time.Now().UTC()
	setFields := bson.M{
		"cardOfDay":             cardOfDay,
		"lastSpreadGeneratedAt": now,
		"updatedAt":             now,
	}

	if deckThemeID != nil {
		setFields["deckThemeId"] = *deckThemeID
	}

	update := bson.M{
		"$set": setFields,
		"$inc": bson.M{"spreadGenerationsCount": 1},
	}
	if deckThemeID == nil {
		update["$unset"] = bson.M{"deckThemeId": ""}
	}

	var out models.DailyReading
	err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update daily reading card of day: %w", err)
	}

	return &out, nil
}

func (r *DailyReadingRepository) UpdateSpreadA(
	ctx context.Context,
	id primitive.ObjectID,
	spreadA []models.CardDraw,
	deckThemeID *primitive.ObjectID,
) (*models.DailyReading, error) {
	now := time.Now().UTC()
	setFields := bson.M{
		"spreadA":               spreadA,
		"lastSpreadGeneratedAt": now,
		"updatedAt":             now,
	}

	if deckThemeID != nil {
		setFields["deckThemeId"] = *deckThemeID
	}

	update := bson.M{
		"$set": setFields,
		"$inc": bson.M{"spreadGenerationsCount": 1},
	}
	if deckThemeID == nil {
		update["$unset"] = bson.M{"deckThemeId": ""}
	}

	var out models.DailyReading
	err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update daily reading spread A: %w", err)
	}

	return &out, nil
}

func (r *DailyReadingRepository) UpdateSpreadB(
	ctx context.Context,
	id primitive.ObjectID,
	spreadB []models.CardDraw,
	deckThemeID *primitive.ObjectID,
) (*models.DailyReading, error) {
	now := time.Now().UTC()
	setFields := bson.M{
		"spreadB":               spreadB,
		"lastSpreadGeneratedAt": now,
		"updatedAt":             now,
	}

	if deckThemeID != nil {
		setFields["deckThemeId"] = *deckThemeID
	}

	update := bson.M{
		"$set": setFields,
		"$inc": bson.M{"spreadGenerationsCount": 1},
	}
	if deckThemeID == nil {
		update["$unset"] = bson.M{"deckThemeId": ""}
	}

	var out models.DailyReading
	err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&out)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update daily reading spread B: %w", err)
	}

	return &out, nil
}
