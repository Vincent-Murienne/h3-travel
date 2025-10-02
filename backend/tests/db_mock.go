package tests

import (
	"h3-travel/config"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Impossible de créer sqlmock: %v", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Impossible d'ouvrir GORM avec sqlmock: %v", err)
	}

	config.DB = gormDB

	cleanup := func() {
		db.Close()
	}

	return mock, cleanup
}
