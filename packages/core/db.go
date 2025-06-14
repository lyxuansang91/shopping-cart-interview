package core

import (
	"database/sql"
	"fmt"
	"strings"
)

// NewMySQLDatabase creates and validates a MySQL database connection
func NewMySQLDatabase(databaseURI string) (*sql.DB, error) {
	fmt.Println("=== Database Connection Debug ===")
	fmt.Printf("Raw databaseURI: %q\n", databaseURI)

	if databaseURI == "" {
		return nil, fmt.Errorf("database URI is empty - check environment variable DATABASE_URI")
	}

	// Format database URI for MySQL driver
	dbURI := strings.TrimPrefix(databaseURI, "mysql://")
	fmt.Printf("After removing prefix: %q\n", dbURI)

	// Convert URI format to MySQL DSN format
	// Format: user:pass@tcp(host:port)/dbname
	// We need to handle cases where the format is already correct
	if !strings.Contains(dbURI, "@tcp(") {
		// Split into user:pass and host:port/dbname
		parts := strings.SplitN(dbURI, "@", 2)
		fmt.Printf("Split parts: %#v\n", parts)

		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid DSN format: expected user:pass@host:port/dbname")
		}

		// Split host:port/dbname
		hostAndDB := strings.SplitN(parts[1], "/", 2)
		fmt.Printf("Host and DB parts: %#v\n", hostAndDB)

		if len(hostAndDB) != 2 {
			return nil, fmt.Errorf("invalid DSN format: missing database name")
		}

		// Reconstruct in MySQL DSN format
		dbURI = fmt.Sprintf("%s@tcp(%s)/%s", parts[0], hostAndDB[0], hostAndDB[1])
		fmt.Printf("Final DSN: %q\n", dbURI)
	}

	// Add parsing parameters for MySQL timestamps
	if !strings.Contains(dbURI, "?") {
		dbURI += "?"
	} else {
		dbURI += "&"
	}
	dbURI += "parseTime=true"

	// Initialize database connection
	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
