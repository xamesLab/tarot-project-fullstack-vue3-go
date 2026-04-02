package repositories

import (
	"context"
	"time"

	"tarot-daily/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PolicyRepository struct {
	collection *mongo.Collection
}

func NewPolicyRepository(db *mongo.Database) *PolicyRepository {
	return &PolicyRepository{collection: db.Collection("role_policies")}
}

func (r *PolicyRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "role", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func (r *PolicyRepository) FindByRole(ctx context.Context, role string) (*models.RolePolicy, error) {
	var policy models.RolePolicy
	if err := r.collection.FindOne(ctx, bson.M{"role": role}).Decode(&policy); err != nil {
		return nil, err
	}
	return &policy, nil
}

func (r *PolicyRepository) ListAll(ctx context.Context) ([]models.RolePolicy, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var policies []models.RolePolicy
	if err := cursor.All(ctx, &policies); err != nil {
		return nil, err
	}

	return policies, nil
}

func (r *PolicyRepository) UpdateByRole(ctx context.Context, role string, threeCardLimit int64, historyRetentionDays int) (*models.RolePolicy, error) {
	if err := r.collection.FindOneAndUpdate(ctx, bson.M{"role": role}, bson.M{
		"$set": bson.M{
			"three_card_limit":       threeCardLimit,
			"history_retention_days": historyRetentionDays,
			"updated_at":             time.Now().UTC(),
		},
	}).Err(); err != nil {
		return nil, err
	}

	return r.FindByRole(ctx, role)
}
