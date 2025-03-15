package interfaces

// BaseRepositoryInterface defines the basic repository contract for checking database connectivity.
type BaseRepositoryInterface interface {
	// IsDown checks if the database is down.
	//
	// It attempts to get the underlying database object and pings it.
	// If any error occurs during this process, it assumes the database is down.
	//
	// Returns:
	//   - bool: true if the database is down, false otherwise
	IsDown() bool
}
