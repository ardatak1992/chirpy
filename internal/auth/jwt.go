package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	var mySigningKey = []byte("super-secret-key")

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy-access",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
			Subject:   userID.String(),
		})

	return token.SignedString(mySigningKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	var claims jwt.Claims

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte("AllYourBase"), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	userIDStr, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.UUID{}, err
	}

	return userID, nil
}
