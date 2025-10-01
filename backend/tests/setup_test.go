package tests

import (
	"h3-travel/config"
	"h3-travel/models"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() {
	var err error
	config.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect test db: %v", err)
	}

	err = config.DB.AutoMigrate(&models.User{}, &models.Voyage{}, &models.Order{})
	if err != nil {
		log.Fatalf("failed to migrate test db: %v", err)
	}
}
