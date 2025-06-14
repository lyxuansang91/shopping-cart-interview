package services

import (
	"github.com/cinchprotocol/cinch-api/packages/core"
	"github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/app/repositories"
	workflowinterfaces "github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/interfaces"
	pkgservices "github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/services"
	pkginterfaces "github.com/cinchprotocol/cinch-api/services/adapterstripe/internal/pkg/services/interfaces"
)

// Services holds all service instances
type Services struct {
	PaymentMethod   pkginterfaces.IPaymentMethodService
	Payment         pkginterfaces.IPaymentService
	Refund          pkginterfaces.IRefundService
	WorkflowPayment workflowinterfaces.PaymentService
}

// NewServices creates a new instance of Services
func NewServices(repos *repositories.Repositories, logger core.Logger) *Services {
	return &Services{
		PaymentMethod:   pkgservices.NewPaymentMethodService(repos.PaymentMethod),
		Payment:         pkgservices.NewPaymentService(logger),
		Refund:          pkgservices.NewRefundService(logger),
		WorkflowPayment: pkgservices.NewWorkflowPaymentService(logger),
	}
}
