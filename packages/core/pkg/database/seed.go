package database

import (
	"database/sql"
	"io/ioutil"
	"log"
	"strings"
)

func RunSeed(db *sql.DB, path string) {
	log.Println("Seeding using:", path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read seed file: %v", err)
	}

	// Split the SQL file into individual statements
	statements := strings.Split(string(data), ";")

	// Execute each statement
	for _, stmt := range statements {
		// Skip empty statements
		if strings.TrimSpace(stmt) == "" {
			continue
		}

		if _, err := db.Exec(stmt); err != nil {
			log.Fatalf("Seed execution failed: %v", err)
		}
	}

	log.Println("Seed completed.")
}
