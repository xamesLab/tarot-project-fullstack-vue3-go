package mongo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"tarotdaily/backend/internal/domain/models"
)

type RolePolicyRepository struct {
	coll *mongo.Collection
}

func NewRolePolicyRepository(db *mongo.Database) *RolePolicyRepository {
	return &RolePolicyRepository{
		coll: db.Collection("role_policies"),
	}
}

func (r *RolePolicyRepository) FindByRole(ctx context.Context, role string) (*models.RolePolicy, error) {
	var item models.RolePolicy
	err := r.coll.FindOne(ctx, bson.M{"role": strings.TrimSpace(role)}).Decode(&item)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find role policy by role: %w", err)
	}
	return &item, nil
}

func (r *RolePolicyRepository) List(ctx context.Context) ([]models.RolePolicy, error) {
	cur, err := r.coll.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: "role", Value: 1}}))
	if err != nil {
		return nil, fmt.Errorf("find role policies: %w", err)
	}
	defer cur.Close(ctx)

	out := make([]models.RolePolicy, 0)
	for cur.Next(ctx) {
		var item models.RolePolicy
		if err := cur.Decode(&item); err != nil {
			return nil, fmt.Errorf("decode role policy: %w", err)
		}
		out = append(out, item)
	}
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("iterate role policies: %w", err)
	}
	return out, nil
}

func (r *RolePolicyRepository) Upsert(ctx context.Context, role string, dailySpreadLimit int) (*models.RolePolicy, error) {
	now := time.Now().UTC()
	filter := bson.M{"role": strings.TrimSpace(role)}
	update := bson.M{
		"$set": bson.M{
			"dailySpreadLimit": dailySpreadLimit,
			"updatedAt":        now,
		},
		"$setOnInsert": bson.M{
			"createdAt": now,
		},
	}

	var out models.RolePolicy
	err := r.coll.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&out)
	if err != nil {
		return nil, fmt.Errorf("upsert role policy: %w", err)
	}

	return &out, nil
}
