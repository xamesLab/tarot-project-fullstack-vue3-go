package repositories

import (
	"context"
	"time"

	"tarot-daily/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CardBackRepository struct {
	collection *mongo.Collection
}

func NewCardBackRepository(db *mongo.Database) *CardBackRepository {
	return &CardBackRepository{collection: db.Collection("card_backs")}
}

func (r *CardBackRepository) CountAll(ctx context.Context) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{})
}

func (r *CardBackRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"status": status})
}

func (r *CardBackRepository) CountByStatuses(ctx context.Context, statuses []string) (int64, error) {
	filter := bson.M{}
	if len(statuses) > 0 {
		filter["status"] = bson.M{"$in": statuses}
	}

	return r.collection.CountDocuments(ctx, filter)
}

func (r *CardBackRepository) Create(ctx context.Context, cardBack *models.CardBack) error {
	_, err := r.collection.InsertOne(ctx, cardBack)
	return err
}

func (r *CardBackRepository) List(ctx context.Context, skip, limit int64) ([]models.CardBack, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().
		SetSort(bson.D{{Key: "updated_at", Value: -1}}).
		SetSkip(skip).
		SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cardBacks []models.CardBack
	if err := cursor.All(ctx, &cardBacks); err != nil {
		return nil, err
	}

	return cardBacks, nil
}

func (r *CardBackRepository) ListByStatuses(ctx context.Context, statuses []string, skip, limit int64) ([]models.CardBack, error) {
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

	var cardBacks []models.CardBack
	if err := cursor.All(ctx, &cardBacks); err != nil {
		return nil, err
	}

	return cardBacks, nil
}

func (r *CardBackRepository) ListAvailable(ctx context.Context) ([]models.CardBack, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"status": "done"}, options.Find().
		SetSort(bson.D{{Key: "updated_at", Value: -1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cardBacks []models.CardBack
	if err := cursor.All(ctx, &cardBacks); err != nil {
		return nil, err
	}

	return cardBacks, nil
}

func (r *CardBackRepository) FindByID(ctx context.Context, cardBackID string) (*models.CardBack, error) {
	var cardBack models.CardBack
	if err := r.collection.FindOne(ctx, bson.M{"_id": cardBackID}).Decode(&cardBack); err != nil {
		return nil, err
	}

	return &cardBack, nil
}

func (r *CardBackRepository) Update(ctx context.Context, cardBackID string, updates bson.M) (*models.CardBack, error) {
	if updates == nil {
		updates = bson.M{}
	}
	updates["updated_at"] = time.Now().UTC()

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var cardBack models.CardBack
	if err := r.collection.FindOneAndUpdate(ctx, bson.M{"_id": cardBackID}, bson.M{"$set": updates}, opts).Decode(&cardBack); err != nil {
		return nil, err
	}

	return &cardBack, nil
}

func (r *CardBackRepository) Delete(ctx context.Context, cardBackID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": cardBackID})
	return err
}
