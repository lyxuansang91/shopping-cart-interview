package test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDB represents a test database instance
type TestDB struct {
	container testcontainers.Container
	DB        *sql.DB
	URI       string
}

// NewTestDB creates a new test database instance
func NewTestDB(t *testing.T, migrationsPath string) *TestDB {
	ctx := context.Background()

	// Start MySQL container
	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "root",
			"MYSQL_DATABASE":      "test",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL"),
	}

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start MySQL container: %v", err)
	}

	// Get container port
	port, err := mysqlC.MappedPort(ctx, "3306")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	// Build connection string
	uri := fmt.Sprintf("root:root@tcp(localhost:%s)/test?parseTime=true", port.Port())

	// Wait for MySQL to be ready
	var db *sql.DB
	for i := 0; i < 10; i++ {
		db, err = sql.Open("mysql", uri)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		t.Fatalf("Failed to connect to MySQL: %v", err)
	}

	// Run migrations if path is provided
	if migrationsPath != "" {
		// Get the absolute path to the workspace root
		cwd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get current working directory: %v", err)
		}
		workspaceRoot, err := findWorkspaceRoot(cwd)
		if err != nil {
			t.Fatalf("Failed to find workspace root: %v", err)
		}

		// Find all .up.sql files in the migrations directory
		migrationsDir := filepath.Join(workspaceRoot, migrationsPath)
		t.Logf("Looking for migrations in: %s", migrationsDir)
		files, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
		if err != nil {
			t.Fatalf("Failed to find migration files: %v", err)
		}

		// Sort files by name to ensure correct order
		sort.Strings(files)
		t.Logf("Found migration files: %v", files)

		// Execute each migration file
		for _, file := range files {
			t.Logf("Reading migration file: %s", file)
			content, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("Failed to read migration file %s: %v", file, err)
			}

			// Split the content into individual statements
			var statements []string
			currentStmt := ""
			for _, line := range strings.Split(string(content), "\n") {
				line = strings.TrimSpace(line)
				if line == "" {
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
				t.Logf("Executing statement: %s", stmt)
				_, err = db.Exec(stmt)
				if err != nil {
					t.Fatalf("Failed to execute migration statement in %s: %v\nStatement: %s", file, err, stmt)
				}
			}
		}

		// Source the db_functions.sql script
		dbFunctionsPath := filepath.Join(workspaceRoot, "scripts", "db_functions.sql")
		t.Logf("Loading database functions from: %s", dbFunctionsPath)
		content, err := os.ReadFile(dbFunctionsPath)
		if err != nil {
			t.Fatalf("Failed to read database functions file: %v", err)
		}

		// Split the content into individual statements
		var statements []string
		currentStmt := ""
		for _, line := range strings.Split(string(content), "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
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
			t.Logf("Executing database function statement: %s", stmt)
			_, err = db.Exec(stmt)
			if err != nil {
				t.Fatalf("Failed to execute database function statement: %v\nStatement: %s", err, stmt)
			}
		}
	}

	return &TestDB{
		container: mysqlC,
		DB:        db,
		URI:       uri,
	}
}

// Close closes the test database connection and container
func (tdb *TestDB) Close(t *testing.T) {
	if tdb.DB != nil {
		tdb.DB.Close()
	}
	if tdb.container != nil {
		if err := tdb.container.Terminate(context.Background()); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}
}

// RunWithDB runs a test function with a test database
func RunWithDB(t *testing.T, migrationsPath string, testFunc func(*testing.T, *TestDB)) {
	tdb := NewTestDB(t, migrationsPath)
	defer tdb.Close(t)
	testFunc(t, tdb)
}

// // MustExec executes a SQL statement and fails the test if it errors
// func (tdb *TestDB) MustExec(t *testing.T, query string, args ...interface{}) {
// 	_, err := tdb.DB.ExecContext(context.Background(), query, args...)
// 	if err != nil {
// 		t.Fatalf("Failed to execute query: %v\nQuery: %s\nArgs: %v", err, query, args)
// 	}
// }

// // MustQuery executes a SQL query and returns the result, failing the test if it errors
// func (tdb *TestDB) MustQuery(t *testing.T, query string, args ...interface{}) *sql.Rows {
// 	rows, err := tdb.DB.QueryContext(context.Background(), query, args...)
// 	if err != nil {
// 		t.Fatalf("Failed to execute query: %v\nQuery: %s\nArgs: %v", err, query, args)
// 	}
// 	return rows
// }

// // MustQueryRow executes a SQL query and returns a single row, failing the test if it errors
// func (tdb *TestDB) MustQueryRow(t *testing.T, query string, args ...interface{}) *sql.Row {
// 	return tdb.DB.QueryRowContext(context.Background(), query, args...)
// }

// // MustBegin starts a transaction and fails the test if it errors
// func (tdb *TestDB) MustBegin(t *testing.T) *sql.Tx {
// 	tx, err := tdb.DB.Begin()
// 	if err != nil {
// 		t.Fatalf("Failed to begin transaction: %v", err)
// 	}
// 	return tx
// }

// // MustCommit commits a transaction and fails the test if it errors
// func (tdb *TestDB) MustCommit(t *testing.T, tx *sql.Tx) {
// 	if err := tx.Commit(); err != nil {
// 		t.Fatalf("Failed to commit transaction: %v", err)
// 	}
// }

// // MustRollback rolls back a transaction and fails the test if it errors
// func (tdb *TestDB) MustRollback(t *testing.T, tx *sql.Tx) {
// 	if err := tx.Rollback(); err != nil {
// 		t.Fatalf("Failed to rollback transaction: %v", err)
// 	}
// }

// // MustExecContext executes a SQL statement with context and fails the test if it errors
// func (tdb *TestDB) MustExecContext(t *testing.T, ctx context.Context, query string, args ...interface{}) {
// 	_, err := tdb.DB.ExecContext(ctx, query, args...)
// 	if err != nil {
// 		t.Fatalf("Failed to execute query with context: %v\nQuery: %s\nArgs: %v", err, query, args)
// 	}
// }

// // MustQueryContext executes a SQL query with context and returns the result, failing the test if it errors
// func (tdb *TestDB) MustQueryContext(t *testing.T, ctx context.Context, query string, args ...interface{}) *sql.Rows {
// 	rows, err := tdb.DB.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		t.Fatalf("Failed to execute query with context: %v\nQuery: %s\nArgs: %v", err, query, args)
// 	}
// 	return rows
// }

// // MustQueryRowContext executes a SQL query with context and returns a single row, failing the test if it errors
// func (tdb *TestDB) MustQueryRowContext(t *testing.T, ctx context.Context, query string, args ...interface{}) *sql.Row {
// 	return tdb.DB.QueryRowContext(ctx, query, args...)
// }

// // MustBeginTx starts a transaction with context and fails the test if it errors
// func (tdb *TestDB) MustBeginTx(t *testing.T, ctx context.Context) *sql.Tx {
// 	tx, err := tdb.DB.BeginTx(ctx, nil)
// 	if err != nil {
// 		t.Fatalf("Failed to begin transaction with context: %v", err)
// 	}
// 	return tx
// }

// // MustCommitTx commits a transaction with context and fails the test if it errors
// func (tdb *TestDB) MustCommitTx(t *testing.T, ctx context.Context, tx *sql.Tx) {
// 	if err := tx.Commit(); err != nil {
// 		t.Fatalf("Failed to commit transaction with context: %v", err)
// 	}
// }

// // MustRollbackTx rolls back a transaction with context and fails the test if it errors
// func (tdb *TestDB) MustRollbackTx(t *testing.T, ctx context.Context, tx *sql.Tx) {
// 	if err := tx.Rollback(); err != nil {
// 		t.Fatalf("Failed to rollback transaction with context: %v", err)
// 	}
// }

// Helper to find the workspace root by searching for a marker file
func findWorkspaceRoot(startDir string) (string, error) {
	markerFiles := []string{"Makefile", ".git"}
	current := startDir
	for i := 0; i < 10; i++ { // avoid infinite loop
		for _, marker := range markerFiles {
			if _, err := os.Stat(filepath.Join(current, marker)); err == nil {
				return current, nil
			}
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", fmt.Errorf("workspace root not found from %s", startDir)
}
