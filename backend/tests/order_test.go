package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"h3-travel/config"
	"h3-travel/controllers"
	"h3-travel/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestOrderCRU(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()

	router := gin.Default()
	router.POST("/orders", controllers.CreateOrder)
	router.GET("/orders/user", controllers.GetUserOrders)
	router.PUT("/orders/:id/cancel", controllers.CancelOrder)

	// Créer voyage pour passer l'ID
	travel := models.Voyage{Title: "Test Trip", Stock: 10, Active: true}
	config.DB.Create(&travel)

	// Créer commande
	orderPayload := map[string]interface{}{
		"voyage_id": travel.ID,
		"card":      "4111111111111111",
	}
	body, _ := json.Marshal(orderPayload)
	req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var order models.Order
	json.Unmarshal(resp.Body.Bytes(), &order)

	// Récupérer commandes
	reqGet := httptest.NewRequest("GET", "/orders/user", nil)
	respGet := httptest.NewRecorder()
	router.ServeHTTP(respGet, reqGet)
	assert.Equal(t, http.StatusOK, respGet.Code)

	// Annuler commande
	reqCancel := httptest.NewRequest("PUT", "/orders/"+string(rune(order.ID))+"/cancel", nil)
	respCancel := httptest.NewRecorder()
	router.ServeHTTP(respCancel, reqCancel)
	assert.Equal(t, http.StatusOK, respCancel.Code)
}
