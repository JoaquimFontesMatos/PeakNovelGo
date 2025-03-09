package config

import (
	"backend/internal/models"
	"fmt"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB establishes a database connection based on the provided boolean flag.  If `isTest` is true, it uses an in-memory
// SQLite database for testing; otherwise, it connects to a PostgreSQL database using environment variables.
//
// Parameters:
//   - isTest (bool): Indicates whether to connect to an in-memory SQLite database (true) or a PostgreSQL database (false).
//
// Returns:
//   - *gorm.DB: A pointer to the established GORM database connection.
//
// Error types:
//   - error: If connection to the database fails, a fatal error is logged, and the program terminates.
func ConnectDB(isTest bool) *gorm.DB {
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

// autoMigrate performs database migrations for all defined models.
//
// This function uses GORM's AutoMigrate function to create or update database tables
// based on the Go structs defined in the models package. It migrates tables for Users,
// RevokedTokens, Novels, Chapters, Tags, NovelTags, NovelAuthors, Authors, BookmarkedNovels,
// NovelGenres, Genres, and LogEntries. Failure to migrate results in a fatal error.
//
// Parameters:
//   - db (*gorm.DB): A pointer to a GORM database connection.
//
// Returns:
//   - : No explicit return value.
//
// Error types:
//   - error:  A fatal error is logged and the program exits if database migration fails.
func autoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.RevokedToken{}, &models.Novel{},
		&models.Chapter{}, &models.Tag{}, &models.NovelTag{}, &models.NovelAuthor{}, &models.Author{},
		&models.BookmarkedNovel{}, &models.NovelGenre{}, &models.Genre{}, &models.LogEntry{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}
	fmt.Println("Database schema migrated successfully!")
}
