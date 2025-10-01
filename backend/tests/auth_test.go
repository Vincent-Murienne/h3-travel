package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"h3-travel/controllers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndLoginUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	SetupTestDB()

	router := gin.Default()
	router.POST("/users/register", controllers.SignUp)
	router.POST("/users/login", controllers.Login)

	userPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(userPayload)

	// Création
	req := httptest.NewRequest("POST", "/users/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	// Login
	reqLogin := httptest.NewRequest("POST", "/users/login", bytes.NewBuffer(body))
	reqLogin.Header.Set("Content-Type", "application/json")
	respLogin := httptest.NewRecorder()
	router.ServeHTTP(respLogin, reqLogin)
	assert.Equal(t, http.StatusOK, respLogin.Code)

	var result map[string]interface{}
	json.Unmarshal(respLogin.Body.Bytes(), &result)
	assert.NotEmpty(t, result["token"])
}
