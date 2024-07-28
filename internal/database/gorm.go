package database

import (
	"github.com/nimyab/anonymous-chat/internal/config"
	"github.com/nimyab/anonymous-chat/internal/database/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectAndMigrateDatabase(cfg *config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&models.User{}, &models.Message{}, &models.Chat{})
	if err != nil {
		panic("failed to migrate database")
	}
	return db
}
