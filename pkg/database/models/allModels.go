package models

import (
	"time"
)

type RefreshToken struct {
	ID        uint `gorm:"primaryKey"`
	UserGUID  string
	TokenHash string
	Session   string `gorm:"unique"`
	IPAddress string
	Used      bool
	CreatedAt time.Time
}
