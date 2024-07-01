package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Article []Article `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Username string
	FullName string
	Email string
	SocialID string
	Provider string
	Avatar string
	Role bool `gorm:default:0`
}