package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"h3-travel/controllers"
	"h3-travel/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTravelCRUD(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()

	router := gin.Default()
	router.POST("/voyages", controllers.CreateVoyage)
	router.GET("/voyages/:id", controllers.GetVoyage)
	router.PUT("/voyages/:id", controllers.UpdateVoyage)
	router.DELETE("/voyages/:id", controllers.DeleteVoyage)

	// Création
	payload := map[string]interface{}{
		"title":  "Paris Trip",
		"stock":  10,
		"active": true,
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest("POST", "/voyages", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var travel models.Voyage
	json.Unmarshal(resp.Body.Bytes(), &travel)

	// Récupération
	reqGet := httptest.NewRequest("GET", "/voyages/"+string(rune(travel.ID)), nil)
	respGet := httptest.NewRecorder()
	router.ServeHTTP(respGet, reqGet)
	assert.Equal(t, http.StatusOK, respGet.Code)

	// Update
	updatePayload := map[string]interface{}{"stock": 5}
	bodyUpdate, _ := json.Marshal(updatePayload)
	reqUpdate := httptest.NewRequest("PUT", "/voyages/"+string(rune(travel.ID)), bytes.NewBuffer(bodyUpdate))
	reqUpdate.Header.Set("Content-Type", "application/json")
	respUpdate := httptest.NewRecorder()
	router.ServeHTTP(respUpdate, reqUpdate)
	assert.Equal(t, http.StatusOK, respUpdate.Code)
}
