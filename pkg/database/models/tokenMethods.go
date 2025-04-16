package models

import (
	"auth/pkg/database"
	"fmt"
)

func (t RefreshToken) Create() error {
	db := database.GetDB()

	return db.Create(&t).Error
}

func FindRefreshBySession(session string) (*RefreshToken, error) {
	db := database.GetDB()
	token := new(RefreshToken)

	err := db.Where(&RefreshToken{Session: session}).Find(&token).Error

	if token.ID == 0 {
		return nil, fmt.Errorf("токен не найден")
	}

	return token, err
}

func (t *RefreshToken) Marked() error {
	db := database.GetDB()

	return db.Model(&t).Update("used", true).Error
}
