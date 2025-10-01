package controllers

import (
	"h3-travel/config"
	"h3-travel/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// SignUp godoc
// @Summary Crée un nouveau compte utilisateur
// @Description Permet à un utilisateur de créer un compte
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body map[string]string true "Email et Password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /signup [post]
func SignUp(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 12)

	user := models.User{Email: input.Email, Password: string(hashed)}
	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email déjà utilisé"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compte créé"})
}

// Login godoc
// @Summary Connecte un utilisateur
// @Description Permet à un utilisateur de se connecter et récupérer un JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body map[string]string true "Email et Password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.First(&user, "email = ?", input.Email).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Identifiants invalides"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Identifiants invalides"})
		return
	}

	// Génération JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString(jwtSecret)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
