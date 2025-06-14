package jobs

import (
	"context"
	"fmt"
	"net"

	"github.com/cinchprotocol/cinch-api/packages/core"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/tracing"
	pb "github.com/cinchprotocol/cinch-api/packages/proto/pkg/proto/assets/cart"
	"github.com/cinchprotocol/cinch-api/services/cart/internal/app"
	"github.com/cinchprotocol/cinch-api/services/cart/internal/app/controllers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// StartGRPCServer starts the gRPC server
func StartGRPCServer(ctx context.Context) error {
	// Start a new span for this function
	ctx, span := tracing.StartFunctionSpan(ctx, "cart", "StartGRPCServer")
	defer span.End()

	// Start the server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", app.Config.GrpcPort))
	if err != nil {
		app.Logger.Error(ctx, "Failed to listen",
			core.NewField("error", err),
		)
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	// Register services
	cartController := controllers.NewCartController(
		app.Services.PaymentMethod,
		app.Services.Payment,
		app.Services.Refund,
	)
	pb.RegisterCartServiceServer(grpcServer, cartController)

	// Register reflection service on gRPC server
	reflection.Register(grpcServer)

	app.Logger.Info(ctx, "Starting gRPC server on port "+app.Config.GrpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		app.Logger.Error(ctx, "Failed to serve",
			core.NewField("error", err),
		)
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}
