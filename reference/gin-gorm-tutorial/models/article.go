package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title string
	Slug string `gorm:"unique_index"`
	Desc string `sql:"type:text;"`
	Tag string
	UserID uint
	CreatedAt time.Time
}