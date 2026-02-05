package models

import (
	"time"

	"gorm.io/gorm"
)

type Script struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Tags       string         `gorm:"index" json:"tags"`        // Comma separated tags
	CategoryID *uint          `gorm:"index" json:"category_id"` // Pointer to allow NULL (uncategorized)
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// ScriptFTS represents the virtual table for Full Text Search
// In GORM we manage FTS tables usually via raw SQL or specific migrations
// for simplicity we will rely on GORM's Like for basic search or configure FTS5 manually
