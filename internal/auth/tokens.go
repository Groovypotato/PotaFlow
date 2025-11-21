package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims holds the user fields embedded in JWTs.
type Claims struct {
	UserID string
	Email  string
}

// GenerateToken issues a signed JWT containing the user ID and email.
func GenerateToken(userID, email string, secret []byte, expiry time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"exp":   time.Now().Add(expiry).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken validates and extracts claims from a JWT.
func ParseToken(tokenStr string, secret []byte) (Claims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secret, nil
	})
	if err != nil {
		return Claims{}, err
	}

	if !token.Valid {
		return Claims{}, jwt.ErrTokenInvalidClaims
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, jwt.ErrTokenInvalidClaims
	}

	sub, _ := mapClaims["sub"].(string)
	email, _ := mapClaims["email"].(string)
	if sub == "" || email == "" {
		return Claims{}, jwt.ErrTokenInvalidClaims
	}

	return Claims{UserID: sub, Email: email}, nil
}
