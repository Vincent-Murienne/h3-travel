package controllers

import (
	"h3-travel/config"
	"h3-travel/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- CREATE ---
// CreateVoyage godoc
// @Summary Crée un voyage
// @Description Permet à un admin de créer un nouveau voyage
// @Tags Travels
// @Accept json
// @Produce json
// @Param voyage body models.Voyage true "Voyage"
// @Success 200 {object} models.Voyage
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /travels [post]
// @Security BearerAuth
func CreateVoyage(c *gin.Context) {
	var voyage models.Voyage
	if err := c.ShouldBindJSON(&voyage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&voyage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, voyage)
}

// --- READ ALL ---
// GetTravels godoc
// @Summary Récupère tous les travels
// @Description Permet à un admin de voir la liste de tous les travels
// @Tags Travels
// @Produce json
// @Success 200 {array} models.Voyage
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /travels [get]
func GetTravels(c *gin.Context) {
	var travels []models.Voyage
	config.DB.Find(&travels)
	c.JSON(http.StatusOK, travels)
}

// --- READ ONE ---
// GetVoyage godoc
// @Summary Récupère un voyage
// @Description Permet à un admin de récupérer un voyage par son ID
// @Tags Travels
// @Produce json
// @Param id path int true "ID du voyage"
// @Success 200 {object} models.Voyage
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /travels/{id} [get]
func GetVoyage(c *gin.Context) {
	id := c.Param("id")
	var voyage models.Voyage
	if err := config.DB.First(&voyage, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voyage non trouvé"})
		return
	}
	c.JSON(http.StatusOK, voyage)
}

// --- UPDATE ---
// UpdateVoyage godoc
// @Summary Met à jour un voyage
// @Description Permet à un admin de mettre à jour un voyage existant
// @Tags Travels
// @Accept json
// @Produce json
// @Param id path int true "ID du voyage"
// @Param voyage body models.Voyage true "Voyage à mettre à jour"
// @Success 200 {object} models.Voyage
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /travels/{id} [put]
// @Security BearerAuth
func UpdateVoyage(c *gin.Context) {
	id := c.Param("id")
	var voyage models.Voyage
	if err := config.DB.First(&voyage, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voyage non trouvé"})
		return
	}

	var input models.Voyage
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&voyage).Updates(input)
	c.JSON(http.StatusOK, voyage)
}

// --- DELETE ---
// DeleteVoyage godoc
// @Summary Supprime un voyage
// @Description Permet à un admin de supprimer un voyage par son ID
// @Tags Travels
// @Produce json
// @Param id path int true "ID du voyage"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /travels/{id} [delete]
// @Security BearerAuth
func DeleteVoyage(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Voyage{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Voyage supprimé"})
}
