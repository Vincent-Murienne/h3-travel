package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID   uint   `json:"user_id"`
	TravelID uint   `json:"travel_id"`
	Statut   string `json:"statut"`
}

type CreateOrderInput struct {
	TravelID uint   `json:"travel_id" binding:"required"`
	Card     string `json:"card" binding:"required"`
}
