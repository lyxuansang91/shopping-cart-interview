package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func NewMySQL(dsn string) *sql.DB {
	if dsn == "" {
		panic("DSN cannot be empty")
	}

	// Convert URL format DSN to Go-MySQL format
	mysqlDsn := strings.TrimPrefix(dsn, "mysql://")

	db, err := sql.Open("mysql", mysqlDsn)
	if err != nil {
		panic(fmt.Sprintf("Failed to open database connection: %v", err))
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify the connection
	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping database: %v", err))
	}

	return db
}
