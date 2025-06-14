package internal

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cinchprotocol/cinch-api/packages/core/pkg/database"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	dbURI    string
	services []string
)

// MigrateCmd represents the migrate command
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run migrations for all services",
	Long:  `Run database migrations for all services or specified services.`,
}

// MigrateUpCmd represents the migrate up command
var MigrateUpCmd = &cobra.Command{
	Use:   "up [revision]",
	Short: "Run migrations up",
	Long:  `Run migrations up to the latest version or to a specific revision.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var targetRevision string
		if len(args) > 0 {
			targetRevision = args[0]
		}
		runMigrationsUp(targetRevision)
	},
}

// MigrateDownCmd represents the migrate down command
var MigrateDownCmd = &cobra.Command{
	Use:   "down <revision>",
	Short: "Run migrations down",
	Long:  `Run migrations down to a specific revision.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runMigrationsDown(args[0])
	},
}

// SeedCmd represents the seed command
var SeedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Run seeders for all services",
	Long:  `Run database seeders for all services or specified services.`,
	Run: func(cmd *cobra.Command, args []string) {
		runSeeders()
	},
}

func init() {
	// Add subcommands to MigrateCmd
	MigrateCmd.AddCommand(MigrateUpCmd)
	MigrateCmd.AddCommand(MigrateDownCmd)

	// Add flags to MigrateUpCmd
	MigrateUpCmd.Flags().StringVarP(&dbURI, "uri", "u", "", "Database URI (required)")
	MigrateUpCmd.Flags().StringSliceVarP(&services, "services", "s", []string{}, "Services to migrate (comma-separated, defaults to all)")
	MigrateUpCmd.MarkFlagRequired("uri")

	// Add flags to MigrateDownCmd
	MigrateDownCmd.Flags().StringVarP(&dbURI, "uri", "u", "", "Database URI (required)")
	MigrateDownCmd.Flags().StringSliceVarP(&services, "services", "s", []string{}, "Services to migrate (comma-separated, defaults to all)")
	MigrateDownCmd.MarkFlagRequired("uri")

	// Add flags to SeedCmd
	SeedCmd.Flags().StringVarP(&dbURI, "uri", "u", "", "Database URI (required)")
	SeedCmd.Flags().StringSliceVarP(&services, "services", "s", []string{}, "Services to seed (comma-separated, defaults to all)")
	SeedCmd.MarkFlagRequired("uri")
}

// runMigrationsUp runs migrations up for all services or specified services
func runMigrationsUp(targetRevision string) {
	// Get the workspace root directory
	workspaceRoot, err := getWorkspaceRoot()
	if err != nil {
		log.Fatalf("Failed to get workspace root: %v", err)
	}

	// Get all services if none specified
	if len(services) == 0 {
		services, err = getServices(workspaceRoot)
		if err != nil {
			log.Fatalf("Failed to get services: %v", err)
		}
	}

	// Parse the database URI
	dbURI := dbURI
	if strings.HasPrefix(dbURI, "mysql://") {
		dbURI = strings.TrimPrefix(dbURI, "mysql://")
	}

	// Run migrations for each service
	for _, service := range services {
		migrationsPath := filepath.Join(workspaceRoot, "services", service, "database", "migrations")

		// Check if migrations directory exists
		if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
			log.Printf("Skipping %s: migrations directory not found", service)
			continue
		}

		log.Printf("Running migrations up for %s service...", service)

		m, err := migrate.New("file://"+migrationsPath, "mysql://"+dbURI)
		if err != nil {
			log.Fatalf("Migration init error for %s: %v", service, err)
		}

		if targetRevision != "" {
			// Get current version
			version, _, err := m.Version()
			if err != nil && err != migrate.ErrNilVersion {
				log.Fatalf("Failed to get current version for %s: %v", service, err)
			}

			// Parse target version
			targetVersion := parseRevisionNumber(targetRevision)
			if targetVersion < version {
				stepsDown := version - targetVersion
				if err := m.Steps(-int(stepsDown)); err != nil && err != migrate.ErrNoChange {
					log.Fatalf("Migration down error for %s: %v", service, err)
				}
			} else if targetVersion > version {
				stepsUp := targetVersion - version
				if err := m.Steps(int(stepsUp)); err != nil && err != migrate.ErrNoChange {
					log.Fatalf("Migration up error for %s: %v", service, err)
				}
			}
		} else {
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Migration up error for %s: %v", service, err)
			}
		}

		log.Printf("Migrations completed for %s service", service)
	}
}

// runMigrationsDown runs migrations down for all services or specified services
func runMigrationsDown(targetRevision string) {
	// Get the workspace root directory
	workspaceRoot, err := getWorkspaceRoot()
	if err != nil {
		log.Fatalf("Failed to get workspace root: %v", err)
	}

	// Get all services if none specified
	if len(services) == 0 {
		services, err = getServices(workspaceRoot)
		if err != nil {
			log.Fatalf("Failed to get services: %v", err)
		}
	}

	// Parse the database URI
	dbURI := dbURI
	if strings.HasPrefix(dbURI, "mysql://") {
		dbURI = strings.TrimPrefix(dbURI, "mysql://")
	}

	// Run migrations for each service
	for _, service := range services {
		migrationsPath := filepath.Join(workspaceRoot, "services", service, "database", "migrations")

		// Check if migrations directory exists
		if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
			log.Printf("Skipping %s: migrations directory not found", service)
			continue
		}

		log.Printf("Running migrations down for %s service...", service)

		m, err := migrate.New("file://"+migrationsPath, "mysql://"+dbURI)
		if err != nil {
			log.Fatalf("Migration init error for %s: %v", service, err)
		}

		// Get current version
		version, _, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Fatalf("Failed to get current version for %s: %v", service, err)
		}

		// Parse target version
		targetVersion := parseRevisionNumber(targetRevision)
		if targetVersion < version {
			stepsDown := version - targetVersion
			if err := m.Steps(-int(stepsDown)); err != nil && err != migrate.ErrNoChange {
				log.Fatalf("Migration down error for %s: %v", service, err)
			}
		}

		log.Printf("Migrations completed for %s service", service)
	}
}

// parseRevisionNumber extracts the revision number from a migration filename
func parseRevisionNumber(revision string) uint {
	// Remove file extension if present
	revision = strings.TrimSuffix(revision, ".up.sql")
	revision = strings.TrimSuffix(revision, ".down.sql")

	// Extract the numeric part
	var version uint
	fmt.Sscanf(revision, "%d", &version)
	return version
}

// runSeeders runs seeders for all services or specified services
func runSeeders() {
	// Get the workspace root directory
	workspaceRoot, err := getWorkspaceRoot()
	if err != nil {
		log.Fatalf("Failed to get workspace root: %v", err)
	}

	// Get all services if none specified
	if len(services) == 0 {
		services, err = getServices(workspaceRoot)
		if err != nil {
			log.Fatalf("Failed to get services: %v", err)
		}
	}

	// Parse the database URI
	// The format is: mysql://user:password@tcp(host:port)/database
	// We need to convert it to: user:password@tcp(host:port)/database
	dbURI := dbURI
	if strings.HasPrefix(dbURI, "mysql://") {
		dbURI = strings.TrimPrefix(dbURI, "mysql://")
	}

	// Connect to the database
	db, err := sql.Open("mysql", dbURI)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run seeders for each service
	for _, service := range services {
		seedPath := filepath.Join(workspaceRoot, "services", service, "database", "seed.sql")

		// Check if seed file exists
		if _, err := os.Stat(seedPath); os.IsNotExist(err) {
			log.Printf("Skipping %s: seed file not found", service)
			continue
		}

		log.Printf("Running seed for %s service...", service)
		database.RunSeed(db, seedPath)
		log.Printf("Seed completed for %s service", service)
	}
}

// getWorkspaceRoot returns the root directory of the workspace
func getWorkspaceRoot() (string, error) {
	// Start from the current directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up until we find the go.work file in the root
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find workspace root (looking for go.work file)")
		}
		dir = parent
	}
}

// getServices returns a list of all services in the workspace
func getServices(workspaceRoot string) ([]string, error) {
	servicesDir := filepath.Join(workspaceRoot, "services")

	// Check if services directory exists
	if _, err := os.Stat(servicesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("services directory not found")
	}

	// Get all directories in the services directory
	entries, err := os.ReadDir(servicesDir)
	if err != nil {
		return nil, err
	}

	var services []string
	for _, entry := range entries {
		if entry.IsDir() {
			services = append(services, entry.Name())
		}
	}

	return services, nil
}
