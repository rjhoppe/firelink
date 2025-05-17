package database

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rjhoppe/firelink/models"
	"github.com/rjhoppe/firelink/ntfy"

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

func BackupDB(filename string) error {
	user := os.Getenv("POSTGRES_USER")
	db := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	password := os.Getenv("POSTGRES_PASSWORD")

	cmd := exec.Command("pg_dump", "-U", user, "-h", host, db)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+password)

	out, err := cmd.Output()
	if err != nil {
			return fmt.Errorf("pg_dump failed: %w", err)
	}

	fileLoc := filepath.Join("/app/database", filename)

	// Write output to file
	if err := os.WriteFile(fileLoc, out, 0644); err != nil {
			return fmt.Errorf("failed to write backup file: %w", err)
	}

	ntfy.NtfyDBBackup(fileLoc)

	return nil
}
