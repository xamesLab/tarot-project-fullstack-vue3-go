package utils

import (
	"fmt"
	"time"

	"tarot-daily/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

const GuestCardTokenPurpose = "guest_card_of_day"

func GenerateGuestCardToken(guestID, cardKey, accountingDate, timezone, secret string, issuedAt, expiresAt time.Time) (string, error) {
	claims := models.GuestCardClaims{
		Purpose:        GuestCardTokenPurpose,
		GuestID:        guestID,
		CardKey:        cardKey,
		AccountingDate: accountingDate,
		Timezone:       timezone,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt.UTC()),
			ExpiresAt: jwt.NewNumericDate(expiresAt.UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func VerifyGuestCardToken(tokenString, secret string) (*models.GuestCardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.GuestCardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.GuestCardClaims)
	if !ok || !token.Valid || claims.Purpose != GuestCardTokenPurpose || claims.GuestID == "" || claims.CardKey == "" || claims.AccountingDate == "" || claims.Timezone == "" {
		return nil, fmt.Errorf("invalid guest card claims")
	}
	return claims, nil
}
