package config

import (
	"backend/internal/models"
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// ConnectDB establishes a database connection with proper production handling
func ConnectDB(isTest bool) *gorm.DB {
	var db *gorm.DB
	var errConnect error

	if isTest {
		// In-memory SQLite for tests
		db, errConnect = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	} else {
		// Get proper database path
		dbPath := getDbPath()

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			log.Fatalf("Failed to create database directory: %v", err)
		}

		// Production SQLite configuration
		connString := fmt.Sprintf("file:%s?cache=shared&_busy_timeout=5000&_journal_mode=WAL", dbPath)
		db, errConnect = gorm.Open(sqlite.Open(connString), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})

		// Configure connection pool
		sqlDB, _ := db.DB()
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	if errConnect != nil {
		log.Fatalf("Failed to connect to database: %v", errConnect)
	}

	// Initialize database (now with migrations)
	initializeDatabase(db)
	return db
}

// getDbPath returns the appropriate database path based on environment
func getDbPath() string {
	if os.Getenv("DEV") == "false" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			log.Printf("Warning: Could not get config dir, using current directory: %v", err)
			return "novels.db"
		}
		return filepath.Join(configDir, "PeakNovelGo", "novels.db")
	}
	return "novels.db" // Dev/test path
}

// initializeDatabase handles both schema migration and data seeding
func initializeDatabase(db *gorm.DB) {
	// First run schema migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
}

// runMigrations executes database schema changes in controlled way
func runMigrations(db *gorm.DB) error {
	// For production, replace AutoMigrate with proper migrations
	if os.Getenv("DEV") != "false" {
		// Development mode - use AutoMigrate for convenience

		autoMigrate(db)
		return nil
	}

	// Production mode - use proper migrations
	// Implementation using golang-migrate:
	// 1. Create migration files in your project (e.g., /migrations/*.sql)
	// 2. Use the migration tool to apply them

	// Example migration setup (pseudo-code):
	/*
		m, err := migrate.New(
			"file:///migrations",
			fmt.Sprintf("sqlite3://%s", getDbPath()))
		if err != nil {
			return err
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return err
		}
	*/

	// For now, we'll fall back to AutoMigrate in production too
	// TODO: Replace with real migrations before production deployment
	autoMigrate(db)
	return nil
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
