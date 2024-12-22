package database

import (
	"cookbook/database/models"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New() *gorm.DB {

	var connString = fmt.Sprintf("user=%s database=%s password=%s host=%s port=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("Failed to connect to database")
	}

	db.AutoMigrate(
		&models.Recepie{},
		&models.Ingredient{},
		&models.Step{},
		&models.Role{},
		&models.User{},
		&models.Session{},
	)

	return db
}
