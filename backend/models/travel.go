package models

import "gorm.io/gorm"

type Travel struct {
	gorm.Model
	Title       string  `gorm:"not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"not null"`
	Stock       int     `gorm:"not null"`
	Active      bool    `gorm:"default:true"`
}
