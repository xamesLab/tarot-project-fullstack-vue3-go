package inpututil

import (
	"errors"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRequireTrimmed(t *testing.T) {
	got, err := RequireTrimmed("  tarot  ")
	if err != nil {
		t.Fatalf("RequireTrimmed() error = %v", err)
	}
	if got != "tarot" {
		t.Fatalf("RequireTrimmed() = %q, want %q", got, "tarot")
	}

	_, err = RequireTrimmed("   ")
	if !errors.Is(err, ErrEmptyString) {
		t.Fatalf("RequireTrimmed(empty) error = %v, want %v", err, ErrEmptyString)
	}
}

func TestParseObjectID(t *testing.T) {
	want := primitive.NewObjectID()
	got, err := ParseObjectID("  " + want.Hex() + "  ")
	if err != nil {
		t.Fatalf("ParseObjectID() error = %v", err)
	}
	if got != want {
		t.Fatalf("ParseObjectID() = %v, want %v", got, want)
	}

	_, err = ParseObjectID("bad-id")
	if !errors.Is(err, ErrInvalidObjectID) {
		t.Fatalf("ParseObjectID(invalid) error = %v, want %v", err, ErrInvalidObjectID)
	}
}

func TestParseOptionalObjectID(t *testing.T) {
	got, err := ParseOptionalObjectID("   ")
	if err != nil {
		t.Fatalf("ParseOptionalObjectID(empty) error = %v", err)
	}
	if got != nil {
		t.Fatalf("ParseOptionalObjectID(empty) = %v, want nil", got)
	}

	want := primitive.NewObjectID()
	got, err = ParseOptionalObjectID(want.Hex())
	if err != nil {
		t.Fatalf("ParseOptionalObjectID() error = %v", err)
	}
	if got == nil || *got != want {
		t.Fatalf("ParseOptionalObjectID() = %v, want %v", got, want)
	}
}
