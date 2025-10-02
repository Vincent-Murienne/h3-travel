package tests

import (
	"bytes"
	"encoding/json"
	"h3-travel/controllers"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestSignUpWithMockDB(t *testing.T) {
	os.Setenv("JWT_SECRET", "secretfortest")

	mock, cleanup := SetupMockDB(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(
			sqlmock.AnyArg(),   // created_at
			sqlmock.AnyArg(),   // updated_at
			sqlmock.AnyArg(),   // deleted_at (NULL aussi accepté)
			"test@example.com", // email
			sqlmock.AnyArg(),   // password hash
			sqlmock.AnyArg(),   // role (souvent "user")
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/signup", controllers.SignUp)

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestLoginWithMockDB(t *testing.T) {
	os.Setenv("JWT_SECRET", "secretfortest")

	mock, cleanup := SetupMockDB(t)
	defer cleanup()

	// --- Cas succès ---
	password := "password123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs("test@example.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "role"}).
			AddRow(1, "test@example.com", string(hashed), "user"))

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/login", controllers.Login)

	payload := map[string]string{
		"email":    "test@example.com",
		"password": password,
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var result map[string]interface{}
	_ = json.Unmarshal(resp.Body.Bytes(), &result)
	assert.NotEmpty(t, result["token"])

	// --- Cas échec : mauvais mot de passe ---
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs("test@example.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "role"}).
			AddRow(1, "test@example.com", string(hashed), "user"))

	badPayload := map[string]string{
		"email":    "test@example.com",
		"password": "wrongpass",
	}
	badBody, _ := json.Marshal(badPayload)

	req2 := httptest.NewRequest("POST", "/login", bytes.NewBuffer(badBody))
	req2.Header.Set("Content-Type", "application/json")
	resp2 := httptest.NewRecorder()

	router.ServeHTTP(resp2, req2)

	assert.Equal(t, http.StatusUnauthorized, resp2.Code)
}
