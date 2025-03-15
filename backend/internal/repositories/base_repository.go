package repositories

import "gorm.io/gorm"

// BaseRepository is a struct that holds a database connection.
// This struct serves as a base for other repositories,
// providing a common database connection for all.
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new BaseRepository.
//
// Parameters:
//   - db (*gorm.DB): The database instance.
//
// Returns:
//   - *BaseRepository: A pointer to the new BaseRepository.
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// IsDown checks if the database is down.
//
// It attempts to get the underlying database object and pings it.
// If any error occurs during this process, it assumes the database is down.
//
// Returns:
//   - bool: true if the database is down, false otherwise
func (n *BaseRepository) IsDown() bool {
	db, err := n.db.DB()
	if err != nil {
		return true // Assume the database is down if we can't get the underlying DB object
	}

	err = db.Ping()
	return err != nil // If Ping returns an error, the database is down
}
