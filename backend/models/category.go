package models

import (
	"time"
)

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique;not null" json:"name"`
	Scripts   []Script  `gorm:"foreignKey:CategoryID" json:"-"` // One-to-Many relationship
	CreatedAt time.Time `json:"created_at"`
}
