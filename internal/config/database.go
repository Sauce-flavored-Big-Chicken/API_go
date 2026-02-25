package config

import (
	"digital-community/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(dbPath string) error {
	var err error

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Rotation{},
		&models.PressCategory{},
		&models.PressNews{},
		&models.Notice{},
		&models.FriendlyNeighbor{},
		&models.FNComment{},
		&models.Activity{},
		&models.ActivityCategory{},
		&models.Registration{},
		&models.Comment{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}
