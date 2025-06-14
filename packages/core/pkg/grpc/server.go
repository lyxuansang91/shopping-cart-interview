package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/cinchprotocol/cinch-api/packages/core"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/middleware"
	"google.golang.org/grpc"
)

// RegisterFunction is a type for the server registration function
type RegisterFunction func(*grpc.Server)

// ServerConfig contains the configuration for the GRPC server
type ServerConfig struct {
	Port        string
	Logger      core.Logger
	ServiceName string
}

// StartServer initializes and starts a gRPC server with the given configuration
func StartServer(ctx context.Context, config ServerConfig, register RegisterFunction) error {
	// Create a new gRPC server with tracing interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.GRPCTracingInterceptor(config.ServiceName)),
	)

	// Register services using the provided function
	register(grpcServer)

	// Start gRPC server
	grpcAddr := fmt.Sprintf(":%s", config.Port)
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	config.Logger.Info(ctx, "gRPC Server starting",
		core.NewField("addr", grpcAddr),
	)

	// Handle context cancellation
	go func() {
		<-ctx.Done()
		config.Logger.Info(ctx, "Shutting down gRPC server gracefully...")
		grpcServer.GracefulStop()
	}()

	// Start server
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve gRPC: %v", err)
	}

	return nil
}
