package jwt

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = os.Getenv("JWT_SECRET")

const expirationTime = time.Minute * 120

type TokenClaims struct {
	UserID string `json:"sub"`
	jwt.RegisteredClaims
}

// GenerateToken crea un token JWT para un usuario
func GenerateToken(userID string) (string, error) {

	claims := TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func validateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("not a valid claims object")

	}
	return claims, nil
}

func GetUserFromTokenStr(tokenString string) (string, error) {
	claims, err := validateToken(tokenString)
	if err != nil {
		return "", errors.New("not a valid token")
	}
	return claims.UserID, nil
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", errors.New("token not found")
	}
	tokenArray := strings.Split(token, " ")

	if len(tokenArray) != 2 || tokenArray[0] != "Bearer" {
		return "", errors.New("authorization header format must be 'Bearer <token>'")
	}

	return tokenArray[1], nil
}
