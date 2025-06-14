package test

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// SeedData represents test data to be seeded into the database
type SeedData struct {
	Table   string
	Columns []string
	Values  [][]interface{}
}

// SeedDatabase seeds the test database with the provided data
func SeedDatabase(t *testing.T, db *sql.DB, data []SeedData) {
	ctx := context.Background()

	// Get the absolute path to the workspace root
	workspaceRoot, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	// Move up to the workspace root
	for i := 0; i < 3; i++ {
		workspaceRoot = filepath.Dir(workspaceRoot)
	}

	// Construct the absolute path to the seed file
	absPath := filepath.Join(workspaceRoot, data[0].Table)
	t.Logf("Loading seed file from: %s", absPath)
	content, err := os.ReadFile(absPath)
	if err != nil {
		t.Fatalf("Failed to read seed file %s: %v", absPath, err)
	}

	// Split the content into individual statements, being careful with UUID_TO_BIN
	var statements []string
	currentStmt := ""
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		// Handle full-line comments
		if strings.HasPrefix(line, "--") {
			continue
		}

		// Remove inline comments
		if idx := strings.Index(line, "--"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}

		currentStmt += " " + line
		if strings.HasSuffix(line, ";") {
			statements = append(statements, strings.TrimSpace(currentStmt))
			currentStmt = ""
		}
	}

	// Execute each statement
	for _, stmt := range statements {
		t.Logf("Executing seed statement: %s", stmt)
		_, err := db.ExecContext(ctx, stmt)
		if err != nil {
			t.Fatalf("Failed to execute seed statement: %v\nStatement: %s", err, stmt)
		}
	}
}

// LoadSeedData is kept for backward compatibility but no longer used
func LoadSeedData(t *testing.T, seedPath string) []SeedData {
	return []SeedData{{Table: seedPath}} // Return a single SeedData with the path
}

// helper function to join strings with a separator
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for _, s := range strs[1:] {
		result += sep + s
	}
	return result
}
