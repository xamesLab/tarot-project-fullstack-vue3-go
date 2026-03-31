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

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{collection: db.Collection("users")}
}

func (r *UserRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "role", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
	})
	return err
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}

	return nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	if err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, id primitive.ObjectID, username, timezone string) (*models.User, error) {
	update := bson.M{"updated_at": time.Now().UTC()}
	if username != "" {
		update["username"] = username
	}
	if timezone != "" {
		update["timezone"] = timezone
	}

	if err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
	).Err(); err != nil {
		return nil, err
	}

	return r.FindByID(ctx, id)
}

func (r *UserRepository) UpdateSettings(ctx context.Context, id primitive.ObjectID, req models.UpdateSettingsRequest) (*models.User, error) {
	update := bson.M{"updated_at": time.Now().UTC()}
	if req.Language != "" {
		update["language"] = req.Language
	}
	if req.Theme != "" {
		update["theme"] = req.Theme
	}
	if req.SelectedDeckID != "" {
		update["selected_deck_id"] = req.SelectedDeckID
	}
	if req.SelectedBackID != "" {
		update["selected_back_id"] = req.SelectedBackID
	}

	if err := r.collection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": update},
	).Err(); err != nil {
		return nil, err
	}

	return r.FindByID(ctx, id)
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, id primitive.ObjectID, lastLogin time.Time) error {
	_, err := r.collection.UpdateByID(ctx, id, bson.M{
		"$set": bson.M{
			"last_login_at": lastLogin,
			"updated_at":    lastLogin,
		},
	})
	return err
}
