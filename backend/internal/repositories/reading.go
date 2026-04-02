package repositories

import (
	"context"
	"time"

	"tarot-daily/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReadingRepository struct {
	collection *mongo.Collection
}

func NewReadingRepository(db *mongo.Database) *ReadingRepository {
	return &ReadingRepository{collection: db.Collection("daily_readings")}
}

func (r *ReadingRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "accounting_date", Value: 1}, {Key: "type", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "generated_at", Value: -1}}},
		{Keys: bson.D{{Key: "expires_at", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0)},
	})
	return err
}

func (r *ReadingRepository) Create(ctx context.Context, reading *models.Reading) error {
	result, err := r.collection.InsertOne(ctx, reading)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		reading.ID = oid
	}

	return nil
}

func (r *ReadingRepository) FindByUserTypeAndAccountingDate(ctx context.Context, userID primitive.ObjectID, readingType, accountingDate string) (*models.Reading, error) {
	var reading models.Reading
	err := r.collection.FindOne(ctx, bson.M{
		"user_id":         userID,
		"type":            readingType,
		"accounting_date": accountingDate,
	}).Decode(&reading)
	if err != nil {
		return nil, err
	}
	return &reading, nil
}

func (r *ReadingRepository) CountByUserTypeAndAccountingDate(ctx context.Context, userID primitive.ObjectID, readingType, accountingDate string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{
		"user_id":         userID,
		"type":            readingType,
		"accounting_date": accountingDate,
	})
}

func (r *ReadingRepository) ListHistory(ctx context.Context, userID primitive.ObjectID, readingType string, skip, limit int64) ([]models.Reading, error) {
	filter := bson.M{"user_id": userID}
	if readingType != "" && readingType != "all" {
		filter["type"] = readingType
	}

	cursor, err := r.collection.Find(ctx, filter, options.Find().
		SetSort(bson.D{{Key: "generated_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var readings []models.Reading
	if err := cursor.All(ctx, &readings); err != nil {
		return nil, err
	}

	return readings, nil
}

func (r *ReadingRepository) CountHistory(ctx context.Context, userID primitive.ObjectID, readingType string) (int64, error) {
	filter := bson.M{"user_id": userID}
	if readingType != "" && readingType != "all" {
		filter["type"] = readingType
	}

	return r.collection.CountDocuments(ctx, filter)
}

func (r *ReadingRepository) CountSince(ctx context.Context, since time.Time) (int64, int64, int64, error) {
	total, err := r.collection.CountDocuments(ctx, bson.M{
		"generated_at": bson.M{"$gte": since},
	})
	if err != nil {
		return 0, 0, 0, err
	}

	cardOfDay, err := r.collection.CountDocuments(ctx, bson.M{
		"generated_at": bson.M{"$gte": since},
		"type":         models.ReadingTypeCardOfDay,
	})
	if err != nil {
		return 0, 0, 0, err
	}

	threeCard, err := r.collection.CountDocuments(ctx, bson.M{
		"generated_at": bson.M{"$gte": since},
		"type":         models.ReadingTypeThreeCard,
	})
	if err != nil {
		return 0, 0, 0, err
	}

	return total, cardOfDay, threeCard, nil
}
