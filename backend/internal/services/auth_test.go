package services

import "testing"

func TestBuildUsernameFromEmail(t *testing.T) {
	if got := buildUsernameFromEmail("alice@example.com"); got != "alice" {
		t.Fatalf("expected alice, got %s", got)
	}

	if got := buildUsernameFromEmail("alice.smith@example.com"); got != "alice.smith" {
		t.Fatalf("expected alice.smith, got %s", got)
	}

	if got := buildUsernameFromEmail("a@example.com"); got != "a01" {
		t.Fatalf("expected a01, got %s", got)
	}
}
