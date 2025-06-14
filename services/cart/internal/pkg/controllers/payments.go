package controllers

// import (
// 	"context"

// 	"github.com/cinchprotocol/cinch-api/packages/core"
// 	"github.com/cinchprotocol/cinch-api/packages/core/pkg/tracing"
// 	pb "github.com/cinchprotocol/cinch-api/packages/proto/pkg/proto/assets/cart"
// 	"github.com/cinchprotocol/cinch-api/services/cart/internal/app"
// 	appservices "github.com/cinchprotocol/cinch-api/services/cart/internal/app/services"
// 	"github.com/cinchprotocol/cinch-api/services/cart/internal/pkg/services"
// )

// type CartController struct {
// 	pb.UnimplementedCartServiceServer
// }

// func NewCartController() *CartController {
// 	return &CartController{}
// }

// // CreatePayment implements the CreatePayment RPC method
// func (c *CartController) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
// 	// Start a new span for this function
// 	ctx, span := tracing.StartFunctionSpan(ctx, "cart", "CreatePayment")
// 	defer span.End()

// 	app.Logger.Info(ctx, "Creating payment",
// 		core.NewField("user_id", req.UserId),
// 		core.NewField("amount", req.Amount),
// 	)

// 	payment, err := appservices.PaymentService.CreatePayment(ctx, services.CreatePaymentParams{
// 		UserID: req.UserId,
// 		Amount: int32(req.Amount * 100),
// 	})
// 	if err != nil {
// 		app.Logger.Error(ctx, "Failed to create payment",
// 			core.NewField("error", err),
// 		)
// 		return nil, err
// 	}

// 	app.Logger.Info(ctx, "Payment created successfully",
// 		core.NewField("payment_id", payment.ID),
// 		core.NewField("status", payment.Status),
// 	)

// 	return &pb.CreatePaymentResponse{
// 		Status: payment.Status,
// 		Amount: float64(payment.Amount) / 100,
// 	}, nil
// }
