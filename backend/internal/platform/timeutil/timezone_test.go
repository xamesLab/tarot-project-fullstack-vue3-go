package timeutil

import (
	"testing"
	"time"
)

func TestBusinessDateAtReset_UTCOffsetBeforeReset(t *testing.T) {
	date, err := BusinessDateAtReset(
		time.Date(2026, 3, 9, 0, 30, 0, 0, time.UTC), // 03:30 local for UTC+03:00
		"UTC+03:00",
		3,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if date != "2026-03-09" {
		t.Fatalf("expected 2026-03-09, got %s", date)
	}
}

func TestBusinessDateAtReset_UTCOffsetEarlyMorning(t *testing.T) {
	date, err := BusinessDateAtReset(
		time.Date(2026, 3, 8, 23, 30, 0, 0, time.UTC), // 02:30 local for UTC+03:00
		"UTC+03:00",
		3,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if date != "2026-03-08" {
		t.Fatalf("expected 2026-03-08, got %s", date)
	}
}

func TestNormalizeTimezone_IANA(t *testing.T) {
	out, err := NormalizeTimezone("Europe/Moscow")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "Europe/Moscow" {
		t.Fatalf("expected Europe/Moscow, got %s", out)
	}
}

func TestNormalizeTimezone_Invalid(t *testing.T) {
	_, err := NormalizeTimezone("UTC+99:00")
	if err == nil {
		t.Fatalf("expected error for invalid timezone")
	}
}
