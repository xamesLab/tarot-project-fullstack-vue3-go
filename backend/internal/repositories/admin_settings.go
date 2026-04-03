package repositories

import (
	"context"
	"time"

	"tarot-daily/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const adminSettingsDocumentID = "global"

type AdminSettingsRepository struct {
	collection *mongo.Collection
}

func NewAdminSettingsRepository(db *mongo.Database) *AdminSettingsRepository {
	return &AdminSettingsRepository{collection: db.Collection("admin_settings")}
}

func (r *AdminSettingsRepository) EnsureDefault(ctx context.Context) error {
	now := time.Now().UTC()
	_, err := r.collection.UpdateByID(ctx, adminSettingsDocumentID, bson.M{
		"$setOnInsert": bson.M{
			"_id":                  adminSettingsDocumentID,
			"maintenance_mode":     false,
			"maintenance_message":  models.LocalizedText{EN: "Site under maintenance", RU: "Сайт на обслуживании"},
			"default_role":         "seeker",
			"default_deck_id":      "default",
			"default_back_id":      "default",
			"registration_enabled": true,
			"created_at":           now,
			"updated_at":           now,
		},
	}, options.Update().SetUpsert(true))
	return err
}

func (r *AdminSettingsRepository) Get(ctx context.Context) (*models.AdminSettings, error) {
	if err := r.EnsureDefault(ctx); err != nil {
		return nil, err
	}

	var settings models.AdminSettings
	if err := r.collection.FindOne(ctx, bson.M{"_id": adminSettingsDocumentID}).Decode(&settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

func (r *AdminSettingsRepository) Update(ctx context.Context, req models.AdminUpdateSettingsRequest) (*models.AdminSettings, error) {
	update := bson.M{"updated_at": time.Now().UTC()}
	if req.MaintenanceMode != nil {
		update["maintenance_mode"] = *req.MaintenanceMode
	}
	if req.MaintenanceMessage != nil {
		update["maintenance_message"] = req.MaintenanceMessage
	}
	if req.DefaultRole != nil {
		update["default_role"] = *req.DefaultRole
	}
	if req.DefaultDeckID != nil {
		update["default_deck_id"] = *req.DefaultDeckID
	}
	if req.DefaultBackID != nil {
		update["default_back_id"] = *req.DefaultBackID
	}
	if req.RegistrationEnabled != nil {
		update["registration_enabled"] = *req.RegistrationEnabled
	}

	if _, err := r.collection.UpdateByID(ctx, adminSettingsDocumentID, bson.M{"$set": update}); err != nil {
		return nil, err
	}

	return r.Get(ctx)
}
