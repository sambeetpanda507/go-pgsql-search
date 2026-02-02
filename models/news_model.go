package models

import (
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Embedding   pgvector.Vector `gorm:"type:vector(384)" json:"embedding"`
}
