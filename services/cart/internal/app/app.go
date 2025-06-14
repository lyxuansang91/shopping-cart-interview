package app

import (
	"context"

	"github.com/cinchprotocol/cinch-api/packages/core"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/config"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/tracing"
	svcconfig "github.com/cinchprotocol/cinch-api/services/cart/configs"
	apprepositories "github.com/cinchprotocol/cinch-api/services/cart/internal/app/repositories"
	appservices "github.com/cinchprotocol/cinch-api/services/cart/internal/app/services"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/sdk/trace"
)

var (
	Logger   core.Logger
	Config   *svcconfig.Config
	Cache    *redis.Client
	tp       *trace.TracerProvider
	Repos    *apprepositories.Repositories
	Services *appservices.Services
)

func Init() error {
	// Parse config
	cfg := &svcconfig.Config{}
	err := config.Parse(cfg)
	if err != nil {
		return err
	}

	// Initialize logger
	logger, err := core.NewZapLogger(true) // Use development mode for now
	if err != nil {
		return err
	}

	// Initialize tracing
	tp, err = tracing.InitTracer("cart", cfg.JaegerEndpoint)
	if err != nil {
		return err
	}

	// Initialize database connection
	db, err := core.NewMySQLDatabase(cfg.DatabaseURI)
	if err != nil {
		return err
	}

	// Initialize Redis connection
	cache, err := core.NewRedisClient(cfg.RedisURI)
	if err != nil {
		return err
	}

	// Initialize repositories
	Repos = apprepositories.NewRepositories(db)

	// Initialize services
	Services = appservices.NewServices(Repos, logger)

	// Bind constructs
	Logger = logger
	Config = cfg
	Cache = cache

	return nil
}

// Shutdown performs cleanup when the application is shutting down
func Shutdown() {
	if tp != nil {
		if err := tp.Shutdown(context.Background()); err != nil {
			Logger.Error(context.Background(), "Failed to shutdown tracer provider",
				core.NewField("error", err),
			)
		}
	}
}
