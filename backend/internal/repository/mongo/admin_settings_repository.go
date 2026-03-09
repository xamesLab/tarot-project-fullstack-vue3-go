package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"tarotdaily/backend/internal/domain/models"
)

type AdminSettingsRepository struct {
	coll *mongo.Collection
}

func NewAdminSettingsRepository(db *mongo.Database) *AdminSettingsRepository {
	return &AdminSettingsRepository{
		coll: db.Collection("admin_settings"),
	}
}

func (r *AdminSettingsRepository) Get(ctx context.Context) (*models.AdminSettings, error) {
	var item models.AdminSettings
	err := r.coll.FindOne(ctx, bson.M{"_id": models.GlobalAdminSettingsID}).Decode(&item)
	if err == nil {
		return &item, nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, fmt.Errorf("find admin settings: %w", err)
	}

	defaultDoc := defaultAdminSettings()
	if _, err := r.coll.InsertOne(ctx, defaultDoc); err != nil {
		if !IsDuplicateKeyError(err) {
			return nil, fmt.Errorf("insert default admin settings: %w", err)
		}
	}

	if err := r.coll.FindOne(ctx, bson.M{"_id": models.GlobalAdminSettingsID}).Decode(&item); err != nil {
		return nil, fmt.Errorf("find admin settings after insert: %w", err)
	}
	return &item, nil
}

func (r *AdminSettingsRepository) Update(ctx context.Context, setFields bson.M) (*models.AdminSettings, error) {
	now := time.Now().UTC()
	updateSet := bson.M{
		"updatedAt": now,
	}
	for key, value := range setFields {
		updateSet[key] = value
	}

	update := bson.M{
		"$set": updateSet,
		"$setOnInsert": bson.M{
			"_id":            models.GlobalAdminSettingsID,
			"createdAt":      now,
			"projectTitle":   "Tarot Daily",
			"supportEmail":   "",
			"maintenanceMode": false,
			"dayResetHour":   3,
			"maxUploadMb":    8,
		},
	}

	var out models.AdminSettings
	if err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"_id": models.GlobalAdminSettingsID},
		update,
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	).Decode(&out); err != nil {
		return nil, fmt.Errorf("upsert admin settings: %w", err)
	}

	return &out, nil
}

func defaultAdminSettings() *models.AdminSettings {
	now := time.Now().UTC()
	return &models.AdminSettings{
		ID:              models.GlobalAdminSettingsID,
		ProjectTitle:    "Tarot Daily",
		SupportEmail:    "",
		MaintenanceMode: false,
		DayResetHour:    3,
		MaxUploadMB:     8,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}
