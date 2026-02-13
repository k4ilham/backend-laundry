package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	ServiceID uint           `json:"service_id"`
	Service   Service        `gorm:"foreignKey:ServiceID" json:"service"`
	Quantity  float64        `json:"quantity"`
	Total     float64        `json:"total"`
	Status    string         `json:"status"` // pending, process, done
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
