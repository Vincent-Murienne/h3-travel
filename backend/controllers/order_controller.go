package controllers

import (
	"h3-travel/config"
	"h3-travel/models"
	"h3-travel/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- CREATE ORDER ---
// CreateOrder godoc
// @Summary Crée une commande
// @Description Permet à un utilisateur de créer une commande pour un voyage
// @Tags Orders
// @Accept json
// @Produce json
// @Param input body models.CreateOrderInput true "Informations pour la commande"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	var input models.CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Vérifie la validité de la carte
	if !utils.ValidateCardNumber(input.Card) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Numéro de carte invalide"})
		return
	}

	// Vérifie le stock
	var voyage models.Voyage
	if err := config.DB.First(&voyage, input.VoyageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voyage non trouvé"})
		return
	}

	if voyage.Stock <= 0 || !voyage.Active {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Voyage indisponible"})
		return
	}

	// Crée la commande
	order := models.Order{
		UserID:   c.GetUint("user_id"),
		VoyageID: input.VoyageID,
		Statut:   "paid",
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Décrémente le stock
	config.DB.Model(&voyage).Update("Stock", voyage.Stock-1)

	c.JSON(http.StatusOK, order)
}

// --- LIST USER ORDERS ---
// GetUserOrders godoc
// @Summary Liste des commandes d'un utilisateur
// @Description Récupère toutes les commandes de l'utilisateur connecté
// @Tags Orders
// @Produce json
// @Success 200 {object} models.Order
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /orders/user [get]
func GetUserOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	var orders []models.Order
	config.DB.Where("user_id = ?", userID.(uint)).Find(&orders)
	c.JSON(http.StatusOK, orders)
}

// --- CANCEL ORDER ---
// CancelOrder godoc
// @Summary Annule une commande
// @Description Permet à un utilisateur d'annuler une commande si elle est encore payée
// @Tags Orders
// @Produce json
// @Param id path int true "ID de la commande"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /orders/{id}/cancel [put]
func CancelOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	orderID := c.Param("id")

	var order models.Order
	if err := config.DB.First(&order, "id = ? AND user_id = ?", orderID, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Commande non trouvée"})
		return
	}

	if order.Statut != "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Impossible d'annuler"})
		return
	}

	order.Statut = "cancelled"
	config.DB.Save(&order)

	// Restock le voyage
	var voyage models.Voyage
	config.DB.First(&voyage, order.VoyageID)
	config.DB.Model(&voyage).Update("Stock", voyage.Stock+1)

	c.JSON(http.StatusOK, order)
}
