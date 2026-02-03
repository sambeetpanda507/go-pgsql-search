package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"

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

	// Run the migration for 000_create_vector_extension.sql
	err = createVectorExtension(db)
	if err != nil {
		log.Fatal(err)
	}

	// Run automigrations for the models
	db.AutoMigrate(&models.News{})

	// Run migrations manually
	if err := runMigrations(db); err != nil {
		log.Fatal("Failed to run migrations: ", err.Error())
	}

	return db
}

func createVectorExtension(db *gorm.DB) error {
	sqlFile, err := os.Open("migrations/000_create_vector_extension.sql")
	if err != nil {
		log.Fatalf("Failed to open 000_create_vector_extension.sql, %v\n", err)
		return err
	}

	sqlBytes, err := io.ReadAll(sqlFile)
	if err != nil {
		log.Fatalf("failed to read 000_create_vector_extension.sql, %v\n", err)
		return err
	}

	if err := db.Exec(string(sqlBytes)).Error; err != nil {
		log.Fatalf("failed to execute 000_create_vector_extension.sql, %v\n", err)
		return err
	}

	fmt.Println("Successfully run migration for 000_create_vector_extension.sql")
	return nil
}

func runMigrations(db *gorm.DB) error {
	// Read files form migrations dir
	files, err := os.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, file := range files {
		// Check if the current file is not a directory and current file extension is .sql
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			migrationFiles = append(migrationFiles, filepath.Join("migrations", file.Name()))
		}
	}

	// Sort the file names in asc order
	sort.Strings(migrationFiles)
	for _, migrationFile := range migrationFiles {
		// Open the file
		file, err := os.Open(migrationFile)
		if err != nil {
			return fmt.Errorf("failed to open %s: %w", migrationFile, err)
		}

		defer file.Close()

		// Read the bytes
		sqlBytes, err := io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", migrationFile, err)
		}

		// Execute the sql
		if err := db.Exec(string(sqlBytes)).Error; err != nil {
			return fmt.Errorf("failed to execute %s: %w", migrationFile, err)
		}

		fmt.Printf("Migration applied: %s\n", migrationFile)
	}

	return nil
}
