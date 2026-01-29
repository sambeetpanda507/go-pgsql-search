package models

import "gorm.io/gorm"

type News struct {
	gorm.Model
	Title       string `json:"title"`
	Description string `json:"description"`
}
