package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func CreateJWT(secret string, expiresAt time.Time) (string, error) {
	defaultExpiration := time.Now().Add(24 * time.Hour)

	if expiresAt.IsZero() || expiresAt.After(defaultExpiration) {
		expiresAt = defaultExpiration
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	})

	jwt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return jwt, nil
}
