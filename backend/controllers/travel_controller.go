package controllers

import (
	"h3-travel/config"
	"h3-travel/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// --- CREATE ---
// CreateTravel godoc
// @Summary Crée un travel
// @Description Permet à un admin de créer un nouveau travel
// @Tags Travels
// @Accept json
// @Produce json
// @Param travel body models.Travel true "Travel"
// @Success 200 {object} models.Travel
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /travels [post]
// @Security BearerAuth
func CreateTravel(c *gin.Context) {
	var travel models.Travel
	if err := c.ShouldBindJSON(&travel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&travel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, travel)
}

// --- READ ALL ---
// GetTravels godoc
// @Summary Récupère tous les travels
// @Description Permet à un admin de voir la liste de tous les travels
// @Tags Travels
// @Produce json
// @Success 200 {array} models.Travel
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /travels [get]
func GetTravels(c *gin.Context) {
	var travels []models.Travel
	config.DB.Find(&travels)
	c.JSON(http.StatusOK, travels)
}

// --- READ ONE ---
// GetTravel godoc
// @Summary Récupère un travel
// @Description Permet à un admin de récupérer un travel par son ID
// @Tags Travels
// @Produce json
// @Param id path int true "ID du travel"
// @Success 200 {object} models.Travel
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /travels/{id} [get]
func GetTravel(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var travel models.Travel
	if err := config.DB.First(&travel, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Travel non trouvé"})
		return
	}

	c.JSON(http.StatusOK, travel)
}

// --- UPDATE ---
// UpdateTravel godoc
// @Summary Met à jour un travel
// @Description Permet à un admin de mettre à jour un travel existant
// @Tags Travels
// @Accept json
// @Produce json
// @Param id path int true "ID du travel"
// @Param travel body models.Travel true "Travel à mettre à jour"
// @Success 200 {object} models.Travel
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /travels/{id} [put]
// @Security BearerAuth
func UpdateTravel(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	var travel models.Travel
	if err := config.DB.First(&travel, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Travel non trouvé"})
		return
	}

	// Bind JSON pour les updates
	var input models.Travel
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&travel).Updates(input)
	c.JSON(http.StatusOK, travel)
}

// --- DELETE ---
// DeleteTravel godoc
// @Summary Supprime un travel
// @Description Permet à un admin de supprimer un travel par son ID
// @Tags Travels
// @Produce json
// @Param id path int true "ID du travel"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /travels/{id} [delete]
// @Security BearerAuth
func DeleteTravel(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID invalide"})
		return
	}

	if err := config.DB.Delete(&models.Travel{}, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Travel supprimé"})
}
