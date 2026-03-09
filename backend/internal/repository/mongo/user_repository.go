package mongo

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"tarotdaily/backend/internal/domain/models"
)

type UserRepository struct {
	coll *mongo.Collection
}

type AdminUserListQuery struct {
	Page     int64
	PageSize int64
	Role     string
	Status   string
	Name     string
	SortAsc  bool
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		coll: db.Collection("users"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	_, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := r.coll.FindOne(ctx, bson.M{"login": login}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find user by login: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.User, error) {
	if len(ids) == 0 {
		return []models.User{}, nil
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cur, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find users by ids: %w", err)
	}
	defer cur.Close(ctx)

	out := make([]models.User, 0, len(ids))
	for cur.Next(ctx) {
		var user models.User
		if err := cur.Decode(&user); err != nil {
			return nil, fmt.Errorf("decode user by ids: %w", err)
		}
		out = append(out, user)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("iterate users by ids: %w", err)
	}

	return out, nil
}

func (r *UserRepository) UpdateSettings(
	ctx context.Context,
	id primitive.ObjectID,
	setFields bson.M,
	unsetFields bson.M,
) (*models.User, error) {
	update := bson.M{
		"$set": bson.M{
			"updatedAt": time.Now().UTC(),
		},
	}

	if len(setFields) > 0 {
		for k, v := range setFields {
			update["$set"].(bson.M)[k] = v
		}
	}
	if len(unsetFields) > 0 {
		update["$unset"] = unsetFields
	}

	var user models.User
	err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update user settings: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) UpdateRoleByLogin(
	ctx context.Context,
	login string,
	role models.UserRole,
) (*models.User, error) {
	var user models.User
	err := r.coll.FindOneAndUpdate(
		ctx,
		bson.M{"login": login},
		bson.M{
			"$set": bson.M{
				"role":      role,
				"updatedAt": time.Now().UTC(),
			},
		},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update user role: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) ListForAdmin(ctx context.Context, query AdminUserListQuery) ([]models.User, int64, error) {
	filter := buildAdminUsersFilter(query)

	total, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count admin users: %w", err)
	}

	sortOrder := int32(-1)
	if query.SortAsc {
		sortOrder = 1
	}

	findOpts := options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: sortOrder}}).
		SetSkip((query.Page - 1) * query.PageSize).
		SetLimit(query.PageSize)

	cur, err := r.coll.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("find admin users: %w", err)
	}
	defer cur.Close(ctx)

	users := make([]models.User, 0, query.PageSize)
	for cur.Next(ctx) {
		var user models.User
		if err := cur.Decode(&user); err != nil {
			return nil, 0, fmt.Errorf("decode admin user: %w", err)
		}
		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate admin users: %w", err)
	}

	return users, total, nil
}

func buildAdminUsersFilter(query AdminUserListQuery) bson.M {
	andConditions := make([]bson.M, 0, 4)

	if role := strings.TrimSpace(strings.ToLower(query.Role)); role != "" {
		andConditions = append(andConditions, bson.M{"role": role})
	}

	if status := strings.TrimSpace(strings.ToLower(query.Status)); status != "" {
		if status == "active" {
			andConditions = append(andConditions, bson.M{
				"$or": []bson.M{
					{"status": "active"},
					{"status": ""},
					{"status": bson.M{"$exists": false}},
				},
			})
		} else {
			andConditions = append(andConditions, bson.M{"status": status})
		}
	}

	if name := strings.TrimSpace(query.Name); name != "" {
		exactNameRegex := primitive.Regex{
			Pattern: "^" + regexp.QuoteMeta(name) + "$",
			Options: "i",
		}
		andConditions = append(andConditions, bson.M{
			"$or": []bson.M{
				{"displayName": exactNameRegex},
				{"login": exactNameRegex},
			},
		})
	}

	if len(andConditions) == 0 {
		return bson.M{}
	}
	if len(andConditions) == 1 {
		return andConditions[0]
	}

	return bson.M{"$and": andConditions}
}

func IsDuplicateKeyError(err error) bool {
	var writeErr mongo.WriteException
	if errors.As(err, &writeErr) {
		for _, e := range writeErr.WriteErrors {
			if e.Code == 11000 {
				return true
			}
		}
	}

	var commandErr mongo.CommandError
	if errors.As(err, &commandErr) {
		return commandErr.Code == 11000
	}
	return false
}

var ErrNotFound = errors.New("not found")
