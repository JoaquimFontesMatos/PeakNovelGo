package config

import (
	"backend/internal/models"
	"fmt"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB connects to the database, either PostgreSQL or SQLite in-memory, based on an environment variable.
func ConnectDB(isTest bool) *gorm.DB {
	// Load environment variables (for PostgreSQL config)
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}

	var db *gorm.DB
	var dsn string
	var errConnect error

	if isTest {
		// In-memory SQLite connection for testing
		db, errConnect = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
		if errConnect != nil {
			log.Fatalf("Failed to connect to in-memory database: %v", errConnect)
		}
		fmt.Println("Connected to in-memory SQLite database for testing.")
	} else {
		// Read environment variables or configuration for PostgreSQL
		dbHost := os.Getenv("DB_HOST")
		dbPort := os.Getenv("DB_PORT")
		dbUser := os.Getenv("DB_USER")
		dbPassword := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		sslMode := os.Getenv("DB_SSL_MODE")

		// Build DSN (Data Source Name) for PostgreSQL
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			dbHost, dbUser, dbPassword, dbName, dbPort, sslMode)

		// Open database connection using GORM for PostgreSQL
		db, errConnect = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if errConnect != nil {
			log.Fatalf("Failed to connect to PostgreSQL database: %v", errConnect)
		}

		// Test PostgreSQL connection (optional)
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to initialize PostgreSQL database connection: %v", err)
		}
		if err := sqlDB.Ping(); err != nil {
			log.Fatalf("Database is unreachable: %v", err)
		}

		fmt.Println("Connected to PostgreSQL database.")
	}

	// Automatically migrate database schema for models
	autoMigrate(db)

	return db
}

func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.RevokedToken{}, &models.Novel{},
		&models.Chapter{}, &models.Tag{}, &models.NovelTag{}, &models.NovelAuthor{}, &models.Author{},
		&models.BookmarkedNovel{}, &models.NovelGenre{}, &models.Genre{}, &models.LogEntry{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
	fmt.Println("Database schema migrated successfully!")
}
