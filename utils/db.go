package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/sambeetpanda507/advance-search/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect() *gorm.DB {
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connected to database")
	db.AutoMigrate(&models.News{})
	return db
}
