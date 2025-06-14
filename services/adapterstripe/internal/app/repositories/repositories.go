package repositories

import (
	"database/sql"

	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/db"
	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/repositories"
)

// Repositories holds all repository instances
type Repositories struct {
	PaymentMethod repositories.PaymentMethodRepository
}

// NewRepositories creates a new Repositories instance with all repositories initialized
func NewRepositories(sqlDB *sql.DB) *Repositories {
	queries := db.New(sqlDB)
	return &Repositories{
		PaymentMethod: repositories.NewPaymentMethodRepository(queries),
	}
}
