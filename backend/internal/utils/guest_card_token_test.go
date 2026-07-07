package utils

import (
	"testing"
	"time"
)

func TestGuestCardTokenRoundTrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	token, err := GenerateGuestCardToken("0f12c9cc-1ef2-4e78-a0cd-3a30e7fa4b1a", "the_high_priestess", "2026-07-03", "Europe/Moscow", "secret", now, now.Add(time.Hour))
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	claims, err := VerifyGuestCardToken(token, "secret")
	if err != nil {
		t.Fatalf("verify token: %v", err)
	}
	if claims.CardKey != "the_high_priestess" || claims.AccountingDate != "2026-07-03" || claims.Timezone != "Europe/Moscow" {
		t.Fatalf("unexpected claims: %#v", claims)
	}
}

func TestGuestCardTokenRejectsWrongSecretAndExpiredToken(t *testing.T) {
	now := time.Now().UTC()
	valid, err := GenerateGuestCardToken("0f12c9cc-1ef2-4e78-a0cd-3a30e7fa4b1a", "the_fool", "2026-07-03", "UTC", "secret", now, now.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := VerifyGuestCardToken(valid, "other-secret"); err == nil {
		t.Fatal("expected token signed with another secret to be rejected")
	}

	expired, err := GenerateGuestCardToken("0f12c9cc-1ef2-4e78-a0cd-3a30e7fa4b1a", "the_fool", "2026-07-02", "UTC", "secret", now.Add(-2*time.Hour), now.Add(-time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := VerifyGuestCardToken(expired, "secret"); err == nil {
		t.Fatal("expected expired token to be rejected")
	}
}
