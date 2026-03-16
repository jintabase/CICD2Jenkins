package model

import "time"

type Article struct {
	ID         string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Title      string    `gorm:"size:255;not null" json:"title"`
	Summary    string    `gorm:"type:text" json:"summary"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	Tags       []string  `gorm:"serializer:json;type:text;not null" json:"tags"`
	Published  bool      `gorm:"not null;default:false" json:"published"`
	AuthorID   string    `gorm:"type:varchar(36);not null;index" json:"author_id"`
	AuthorName string    `gorm:"size:64;not null" json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
