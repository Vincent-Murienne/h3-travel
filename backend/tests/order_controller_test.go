package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"h3-travel/controllers"
	"h3-travel/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrderWithMockDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mock, cleanup := SetupMockDB(t)
	defer cleanup()

	userID := uint(1)
	travelID := uint(2)
	now := time.Now()

	mock.ExpectQuery(`SELECT \* FROM "travels" WHERE "travels"\."id" = \$1 AND "travels"\."deleted_at" IS NULL ORDER BY "travels"\."id" LIMIT \$2`).
		WithArgs(int64(travelID), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "stock", "active", "created_at", "updated_at"}).
			AddRow(travelID, "Test Trip", 10, true, now, now))

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "orders" .* RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			nil,              // deleted_at
			userID,           // user_id
			travelID,         // travel_id
			"paid",           // statut
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "travels"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := gin.Default()
	router.POST("/orders", func(c *gin.Context) {
		c.Set("user_id", userID)
		controllers.CreateOrder(c)
	})

	input := models.CreateOrderInput{
		TravelID: travelID,
		Card:     "4242424242424242",
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/orders", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var result models.Order
	_ = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, travelID, result.TravelID)
	assert.Equal(t, "paid", result.Statut)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// ----------------------
// TEST GET USER ORDERS
// ----------------------
func TestGetUserOrdersWithMockDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock, cleanup := SetupMockDB(t)
	defer cleanup()

	userID := uint(1)
	travelID := uint(2)
	now := time.Now()

	// Mock SELECT orders
	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "travel_id", "statut", "created_at", "updated_at"}).
			AddRow(1, userID, travelID, "paid", now, now).
			AddRow(2, userID, travelID, "paid", now, now))

	router := gin.Default()
	router.GET("/orders/user", func(c *gin.Context) {
		c.Set("user_id", userID)
		controllers.GetUserOrders(c)
	})

	req := httptest.NewRequest("GET", "/orders/user", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var orders []models.Order
	_ = json.Unmarshal(resp.Body.Bytes(), &orders)
	assert.Len(t, orders, 2)
	assert.Equal(t, userID, orders[0].UserID)
	assert.Equal(t, userID, orders[1].UserID)
}

// ----------------------
// TEST CANCEL ORDER
// ----------------------
func TestCancelOrderWithMockDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock, cleanup := SetupMockDB(t)
	defer cleanup()

	orderID := uint(1)
	userID := uint(1)
	travelID := uint(2)
	now := time.Now()

	mock.ExpectQuery(`SELECT \* FROM "orders" WHERE \(id = \$1 AND user_id = \$2\) AND "orders"\."deleted_at" IS NULL ORDER BY "orders"\."id" LIMIT \$3`).
		WithArgs("1", userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "travel_id", "statut", "created_at", "updated_at", "deleted_at"}).
			AddRow(orderID, userID, travelID, "paid", now, now, nil))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "orders" SET`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery(`SELECT \* FROM "travels" WHERE "travels"\."id" = \$1 AND "travels"\."deleted_at" IS NULL ORDER BY "travels"\."id" LIMIT \$2`).
		WithArgs(int64(travelID), 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "stock", "active", "created_at", "updated_at", "deleted_at"}).
			AddRow(travelID, "Test Trip", 9, true, now, now, nil))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "travels"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	router := gin.Default()
	router.PUT("/orders/:id/cancel", func(c *gin.Context) {
		c.Set("user_id", userID)
		controllers.CancelOrder(c)
	})

	req := httptest.NewRequest("PUT", "/orders/1/cancel", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result models.Order
	_ = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.Equal(t, "cancelled", result.Statut)
	assert.Equal(t, orderID, result.ID)
	assert.Equal(t, userID, result.UserID)

	assert.NoError(t, mock.ExpectationsWereMet())
}
