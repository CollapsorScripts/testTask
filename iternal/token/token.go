package token

import (
	"auth/pkg/config"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func Create(cfg *config.Config, guid, ip string) (string, string, error) {
	uniqueSession := uuid.New().String()

	claims := jwt.MapClaims{
		"sub":     guid,
		"exp":     time.Now().Add(time.Minute * 30).Unix(),
		"iat":     time.Now().Unix(),
		"session": uniqueSession,
		"ip":      ip,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signed, err := token.SignedString([]byte(cfg.JwtSecret))
	if err != nil {
		return "", "", err
	}

	return signed, uniqueSession, nil
}

func CreateRefresh() (string, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}

	token := base64.StdEncoding.EncodeToString(bytes)
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	return token, string(hash), nil // Возвращаем токен и его хеш
}

func ValidateRefresh(token, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(token)) == nil
}

func Parse(cfg *config.Config, tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != "HS512" {
			return nil, fmt.Errorf("неверный алгоритм подписи токена: %v", token.Header["alg"])
		}
		return []byte(cfg.JwtSecret), nil
	})
}
