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

type AskRequestRepository struct {
	collection *mongo.Collection
}

func NewAskRequestRepository(db *mongo.Database) *AskRequestRepository {
	return &AskRequestRepository{collection: db.Collection("ask_requests")}
}

func (r *AskRequestRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "requested_at", Value: -1}}},
		{Keys: bson.D{{Key: "reading_id", Value: 1}, {Key: "requested_at", Value: -1}}},
		{Keys: bson.D{{Key: "status", Value: 1}, {Key: "updated_at", Value: -1}}},
	})
	return err
}

func (r *AskRequestRepository) Create(ctx context.Context, askRequest *models.AskRequest) error {
	result, err := r.collection.InsertOne(ctx, askRequest)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		askRequest.ID = oid
	}

	return nil
}

func (r *AskRequestRepository) FindByIDAndUser(ctx context.Context, requestID, userID primitive.ObjectID) (*models.AskRequest, error) {
	var askRequest models.AskRequest
	err := r.collection.FindOne(ctx, bson.M{
		"_id":     requestID,
		"user_id": userID,
	}).Decode(&askRequest)
	if err != nil {
		return nil, err
	}

	return &askRequest, nil
}

func (r *AskRequestRepository) FindLatestByReadingAndUser(ctx context.Context, readingID, userID primitive.ObjectID) (*models.AskRequest, error) {
	var askRequest models.AskRequest
	err := r.collection.FindOne(
		ctx,
		bson.M{
			"reading_id": readingID,
			"user_id":    userID,
		},
		options.FindOne().SetSort(bson.D{{Key: "requested_at", Value: -1}}),
	).Decode(&askRequest)
	if err != nil {
		return nil, err
	}

	return &askRequest, nil
}

func (r *AskRequestRepository) UpdateStatus(ctx context.Context, requestID primitive.ObjectID, status string, completedAt *time.Time, responseExcerpt, errorCode, errorMessage string, usage *models.AskUsage, responseData map[string]any) error {
	update := bson.M{
		"status":           status,
		"updated_at":       time.Now().UTC(),
		"error_code":       errorCode,
		"error_message":    errorMessage,
		"response_excerpt": responseExcerpt,
		"usage":            usage,
		"response_data":    responseData,
	}

	if completedAt != nil {
		update["completed_at"] = *completedAt
	}

	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": requestID}, bson.M{"$set": update})
	return err
}

func (r *AskRequestRepository) CountByUserSince(ctx context.Context, userID primitive.ObjectID, since time.Time) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{
		"user_id":      userID,
		"requested_at": bson.M{"$gte": since},
	})
}
