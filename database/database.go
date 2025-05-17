package database

import (
	"fmt"
	"log"
	"os"

	"github.com/rjhoppe/firelink/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Example: use environment variables for config
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	// sslmode := os.Getenv("POSTGRES_SSLMODE") // usually "disable" for local dev

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, dbname, port, // sslmode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	// Migrate the schema
	DB.AutoMigrate(&models.Dinner{}, &models.Drink{})
}

func GetDB() *gorm.DB {
	return DB
}

func SaveToDB[T any](db *gorm.DB, value *T) error {
	return db.Create(value).Error
}

func CheckRecordExists[T any](db *gorm.DB, value *T) (bool, error) {
	var count int64
	if err := db.Model(value).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetRecord(record interface{}) error {
	return DB.First(record).Error
}