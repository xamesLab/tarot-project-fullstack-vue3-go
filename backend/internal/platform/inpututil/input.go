package inpututil

import (
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrEmptyString     = errors.New("empty string")
	ErrInvalidObjectID = errors.New("invalid object id")
)

func RequireTrimmed(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", ErrEmptyString
	}
	return trimmed, nil
}

func ParseObjectID(raw string) (primitive.ObjectID, error) {
	trimmed, err := RequireTrimmed(raw)
	if err != nil {
		return primitive.NilObjectID, err
	}

	oid, err := primitive.ObjectIDFromHex(trimmed)
	if err != nil {
		return primitive.NilObjectID, ErrInvalidObjectID
	}

	return oid, nil
}

func ParseOptionalObjectID(raw string) (*primitive.ObjectID, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	oid, err := primitive.ObjectIDFromHex(trimmed)
	if err != nil {
		return nil, ErrInvalidObjectID
	}

	return &oid, nil
}
