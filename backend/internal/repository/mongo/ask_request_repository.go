package mongo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"tarotdaily/backend/internal/domain/models"
)

type AskRequestRepository struct {
	coll *mongo.Collection
}

func NewAskRequestRepository(db *mongo.Database) *AskRequestRepository {
	return &AskRequestRepository{
		coll: db.Collection("ask_requests"),
	}
}

func (r *AskRequestRepository) Create(ctx context.Context, item *models.AskRequest) error {
	if item == nil {
		return errors.New("ask request is nil")
	}
	if strings.TrimSpace(item.RequestID) == "" {
		return errors.New("request id is required")
	}
	if item.UserID.IsZero() {
		return errors.New("user id is required")
	}
	if strings.TrimSpace(item.Text) == "" {
		return errors.New("text is required")
	}

	now := time.Now().UTC()
	if item.ID.IsZero() {
		item.ID = primitive.NewObjectID()
	}
	if item.Status == "" {
		item.Status = models.AskRequestStatusAccepted
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	item.UpdatedAt = now

	_, err := r.coll.InsertOne(ctx, item)
	if err != nil {
		return fmt.Errorf("insert ask request: %w", err)
	}
	return nil
}

func (r *AskRequestRepository) FindByRequestID(ctx context.Context, requestID string) (*models.AskRequest, error) {
	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		return nil, errors.New("request id is required")
	}

	var item models.AskRequest
	err := r.coll.FindOne(ctx, bson.M{"requestId": requestID}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find ask request by request id: %w", err)
	}
	return &item, nil
}

func (r *AskRequestRepository) UpdateStatus(
	ctx context.Context,
	requestID string,
	status models.AskRequestStatus,
) (*models.AskRequest, error) {
	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		return nil, errors.New("request id is required")
	}
	if status == "" {
		return nil, errors.New("status is required")
	}

	var item models.AskRequest
	err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"requestId": requestID},
		bson.M{
			"$set": bson.M{
				"status":    status,
				"updatedAt": time.Now().UTC(),
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update ask request status: %w", err)
	}

	return &item, nil
}

func (r *AskRequestRepository) UpdateResult(
	ctx context.Context,
	requestID string,
	resultText string,
) (*models.AskRequest, error) {
	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		return nil, errors.New("request id is required")
	}
	if strings.TrimSpace(resultText) == "" {
		return nil, errors.New("result text is required")
	}

	now := time.Now().UTC()

	var item models.AskRequest
	err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"requestId": requestID},
		bson.M{
			"$set": bson.M{
				"status":      models.AskRequestStatusSuccess,
				"resultText":  resultText,
				"updatedAt":   now,
				"completedAt": now,
			},
			"$unset": bson.M{
				"errorCode":    "",
				"errorMessage": "",
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update ask request result: %w", err)
	}

	return &item, nil
}

func (r *AskRequestRepository) UpdateError(
	ctx context.Context,
	requestID string,
	errorCode string,
	errorMessage string,
) (*models.AskRequest, error) {
	requestID = strings.TrimSpace(requestID)
	if requestID == "" {
		return nil, errors.New("request id is required")
	}
	if strings.TrimSpace(errorCode) == "" {
		return nil, errors.New("error code is required")
	}
	if strings.TrimSpace(errorMessage) == "" {
		return nil, errors.New("error message is required")
	}

	now := time.Now().UTC()

	var item models.AskRequest
	err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"requestId": requestID},
		bson.M{
			"$set": bson.M{
				"status":       models.AskRequestStatusError,
				"errorCode":    errorCode,
				"errorMessage": errorMessage,
				"updatedAt":    now,
				"completedAt":  now,
			},
			"$unset": bson.M{
				"resultText": "",
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update ask request error: %w", err)
	}

	return &item, nil
}
