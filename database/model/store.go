package model

import (
	"time"

	"gorm.io/gorm"
)

type Store struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"unique"`
	Slug      string         `json:"slug"`
	UserID    uint           `json:"userId"`
	Products  []Product      `json:"products" gorm:"foreignKey:StoreID"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
