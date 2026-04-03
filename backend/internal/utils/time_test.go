package utils

import (
	"testing"
	"time"
)

func TestGetAccountingDateForTimeBeforeThreeAM(t *testing.T) {
	loc := time.FixedZone("UTC+2", 2*60*60)
	now := time.Date(2026, 3, 31, 2, 59, 0, 0, loc)

	got := getAccountingDateForTime(now)

	if got != "2026-03-30" {
		t.Fatalf("expected 2026-03-30, got %s", got)
	}
}

func TestGetAccountingDateForTimeAtThreeAM(t *testing.T) {
	loc := time.FixedZone("UTC+2", 2*60*60)
	now := time.Date(2026, 3, 31, 3, 0, 0, 0, loc)

	got := getAccountingDateForTime(now)

	if got != "2026-03-31" {
		t.Fatalf("expected 2026-03-31, got %s", got)
	}
}

func TestGetNextResetTimeForTimeBeforeThreeAM(t *testing.T) {
	loc := time.FixedZone("UTC+2", 2*60*60)
	now := time.Date(2026, 3, 31, 2, 59, 0, 0, loc)

	got := getNextResetTimeForTime(now)
	want := time.Date(2026, 3, 31, 3, 0, 0, 0, loc)

	if !got.Equal(want) {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestGetNextResetTimeForTimeAfterThreeAM(t *testing.T) {
	loc := time.FixedZone("UTC-5", -5*60*60)
	now := time.Date(2026, 3, 31, 11, 15, 0, 0, loc)

	got := getNextResetTimeForTime(now)
	want := time.Date(2026, 4, 1, 3, 0, 0, 0, loc)

	if !got.Equal(want) {
		t.Fatalf("expected %s, got %s", want, got)
	}
}

func TestValidateTimezone(t *testing.T) {
	if !ValidateTimezone("Europe/Amsterdam") {
		t.Fatal("expected Europe/Amsterdam to be valid")
	}

	if ValidateTimezone("Not/AZone") {
		t.Fatal("expected Not/AZone to be invalid")
	}
}
