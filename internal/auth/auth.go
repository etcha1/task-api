package auth

import (
	"os"

	"github.com/go-chi/jwtauth/v5"
)

var TokenAuth *jwtauth.JWTAuth

func Initialize() {
	InitializeWithSecret(getSecret())
}

func InitializeWithSecret(secret string) {
	if secret == "" {
		secret = "dev-secret-change-me"
	}

	TokenAuth = jwtauth.New("HS256", []byte(secret), nil)
}

func getSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}

	return ""
}

func NewToken(userID int) (string, error) {
	if TokenAuth == nil {
		Initialize()
	}

	_, tokenString, err := TokenAuth.Encode(map[string]interface{}{"user_id": userID})
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
