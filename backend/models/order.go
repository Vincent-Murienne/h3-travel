package models

type Order struct {
	ID       uint   `json:"id"`
	UserID   uint   `json:"user_id"`
	VoyageID uint   `json:"voyage_id"`
	Statut   string `json:"statut"`
}

type CreateOrderInput struct {
	VoyageID uint   `json:"voyage_id" binding:"required"`
	Card     string `json:"card" binding:"required"`
}
