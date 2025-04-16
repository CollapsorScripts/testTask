package token

import (
	"auth/pkg/config"
	"github.com/google/uuid"
	"strings"
	"testing"
)

const configPathLocal = "../../config/local.yaml"

func TestCreate(t *testing.T) {
	cfg := config.MustLoadByPath(configPathLocal)
	guid := uuid.New().String()
	ip := "192.168.1.106"
	token, session, err := Create(cfg, guid, ip)
	if err != nil {
		t.Errorf("Ошибка при создании токена: %v", err)
		return
	}

	if token == "" || session == "" {
		t.Error("токен или сессия пустые")
		return
	}

	if !strings.Contains(token, ".") {
		t.Error("некорректный формат JWT")
		return
	}

	t.Logf("Токен успешно сгенерирован: \nToken: %s\nSession: %s", token, session)
}

func TestCreateRefresh(t *testing.T) {
	refreshToken, hash, err := CreateRefresh()
	if err != nil {
		t.Errorf("Ошибка при создании рефреш токена: %v", err)
		return
	}

	t.Logf("Рефреш токен успешно сгенерирован: \nRefreshToken: %s\nHash: %s", refreshToken, hash)
}

func TestValidateRefresh_HappyPath(t *testing.T) {
	token, hash, err := CreateRefresh()
	if err != nil {
		t.Errorf("Ошибка при создании рефреш токена: %v", err)
		return
	}
	isValid := ValidateRefresh(token, hash)
	if !isValid {
		t.Logf("Токен не валиден")
		return
	}

	t.Logf("Токен валиден")
}
