package tests

import (
	"bytes"
	"encoding/json"
	"h3-travel/controllers"
	"h3-travel/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// --- CREATE TRAVEL ---
func TestCreateTravelWithMockDB(t *testing.T) {
	mock, cleanup := SetupMockDB(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "travels"`).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			sqlmock.AnyArg(), // deleted_at
			"Découverte de Paris",
			"Visitez les monuments emblématiques de Paris en 3 jours.",
			299.99,
			10,
			true,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/travels", controllers.CreateTravel)

	payload := models.Travel{
		Title:       "Découverte de Paris",
		Description: "Visitez les monuments emblématiques de Paris en 3 jours.",
		Price:       299.99,
		Stock:       10,
		Active:      true,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/travels", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// --- GET ALL TRAVELS ---
func TestGetTravelsWithMockDB(t *testing.T) {
	mock, cleanup := SetupMockDB(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "title", "description", "price", "stock", "active"}).
		AddRow(1, "Découverte de Paris", "Visitez les monuments", 299.99, 10, true).
		AddRow(2, "Safari en Afrique", "Safari inoubliable", 1499.50, 5, true)

	mock.ExpectQuery(`SELECT \* FROM "travels"`).WillReturnRows(rows)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/travels", controllers.GetTravels)

	req := httptest.NewRequest("GET", "/travels", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var travels []models.Travel
	_ = json.Unmarshal(resp.Body.Bytes(), &travels)
	assert.Len(t, travels, 2)
}

// --- GET ONE TRAVEL ---
func TestGetTravelWithMockDB(t *testing.T) {
	mock, cleanup := SetupMockDB(t)
	defer cleanup()

	row := sqlmock.NewRows([]string{"id", "title", "description", "price", "stock", "active"}).
		AddRow(1, "Découverte de Paris", "Visitez les monuments", 299.99, 10, true)

	mock.ExpectQuery(`SELECT \* FROM "travels" WHERE "travels"\."id" = \$1 AND "travels"\."deleted_at" IS NULL ORDER BY "travels"\."id" LIMIT \$2`).
		WithArgs(int64(1), sqlmock.AnyArg()).
		WillReturnRows(row)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/travels/:id", controllers.GetTravel)

	req := httptest.NewRequest("GET", "/travels/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var travel models.Travel
	_ = json.Unmarshal(resp.Body.Bytes(), &travel)
	assert.Equal(t, "Découverte de Paris", travel.Title)
}

// --- UPDATE TRAVEL ---
func TestUpdateTravelWithMockDB(t *testing.T) {
	mock, cleanup := SetupMockDB(t)
	defer cleanup()

	row := sqlmock.NewRows([]string{"id", "title", "description", "price", "stock", "active"}).
		AddRow(1, "Découverte de Paris", "Visitez les monuments", 299.99, 10, true)
	mock.ExpectQuery(`SELECT \* FROM "travels" WHERE "travels"\."id" = \$1 AND "travels"\."deleted_at" IS NULL ORDER BY "travels"\."id" LIMIT \$2`).
		WithArgs(int64(1), sqlmock.AnyArg()).
		WillReturnRows(row)

	mock.ExpectExec(`UPDATE "travels"`).
		WithArgs(sqlmock.AnyArg(), "Paris by Night", sqlmock.AnyArg(), 10, true, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/travels/:id", controllers.UpdateTravel)

	payload := models.Travel{
		Title: "Paris by Night",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("PUT", "/travels/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// --- DELETE TRAVEL ---
func TestDeleteTravelWithMockDB(t *testing.T) {
	mock, cleanup := SetupMockDB(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "travels" SET "deleted_at"=\$1 WHERE "travels"."id" = \$2 AND "travels"."deleted_at" IS NULL`).
		WithArgs(sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/travels/:id", controllers.DeleteTravel)

	req := httptest.NewRequest("DELETE", "/travels/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}
