package models

import "testing"

func TestNormalizeRole(t *testing.T) {
	got := NormalizeRole("  AdMiN  ")
	if got != "admin" {
		t.Fatalf("NormalizeRole() = %q, want %q", got, "admin")
	}
}

func TestIsValidRole(t *testing.T) {
	if !IsValidRole("Tarot_guide") {
		t.Fatalf("IsValidRole() = false, want true")
	}
	if IsValidRole("1admin") {
		t.Fatalf("IsValidRole() = true, want false")
	}
}

func TestDisplayNameOrLogin(t *testing.T) {
	if got := DisplayNameOrLogin("  Mystic Reader  ", "reader42"); got != "Mystic Reader" {
		t.Fatalf("DisplayNameOrLogin() = %q, want %q", got, "Mystic Reader")
	}
	if got := DisplayNameOrLogin("   ", "reader42"); got != "reader42" {
		t.Fatalf("DisplayNameOrLogin() = %q, want %q", got, "reader42")
	}
}

func TestNormalizeUserStatus(t *testing.T) {
	if got := NormalizeUserStatus(""); got != UserStatusActive {
		t.Fatalf("NormalizeUserStatus(empty) = %q, want %q", got, UserStatusActive)
	}
	if got := NormalizeUserStatus(UserStatusBanned); got != UserStatusBanned {
		t.Fatalf("NormalizeUserStatus(banned) = %q, want %q", got, UserStatusBanned)
	}
}
