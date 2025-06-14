package mappers

import (
	"time"

	models "github.com/cinchprotocol/cinch-api/packages/core/pkg/models"
	"github.com/cinchprotocol/cinch-api/packages/core/pkg/uuid"
	"github.com/cinchprotocol/cinch-api/packages/proto/pkg/proto/assets/cart"
	sqlcmodels "github.com/cinchprotocol/cinch-api/services/cart/internal/pkg/db"
)

// ToDomainPaymentMethod converts a SQLC payment method to a models payment method
func ToDomainPaymentMethod(pm *sqlcmodels.PaymentMethod) *models.PaymentMethod {
	if pm == nil {
		return nil
	}

	return &models.PaymentMethod{
		ID:              uuid.ID(pm.ID),
		PartnerID:       uuid.ID(pm.ID), // Using the same ID for now
		PartnerCode:     pm.PaymentMethodCode,
		PartnerMethodID: pm.PartnerPmType,
		Type:            models.PaymentMethodTypeCard, // Default to card type
		Status:          models.PaymentMethodStatus(pm.Status),
		CreatedAt:       &pm.CreatedAt,
		UpdatedAt:       &pm.UpdatedAt,
	}
}

// ToDomainPaymentMethods converts a slice of SQLC payment methods to models payment methods
func ToDomainPaymentMethods(pms []sqlcmodels.PaymentMethod) []*models.PaymentMethod {
	if len(pms) == 0 {
		return nil
	}

	result := make([]*models.PaymentMethod, len(pms))
	for i, pm := range pms {
		result[i] = ToDomainPaymentMethod(&pm)
	}
	return result
}

// MapProtoToDomainPaymentMethod maps a proto payment method to a domain payment method
func MapProtoToDomainPaymentMethod(paymentMethod *cart.PaymentMethod) (*models.PaymentMethod, error) {
	if paymentMethod == nil {
		return nil, nil
	}

	id, err := uuid.Parse(paymentMethod.Id)
	if err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, paymentMethod.CreatedAt)
	if err != nil {
		return nil, err
	}

	updatedAt, err := time.Parse(time.RFC3339, paymentMethod.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &models.PaymentMethod{
		ID:              id,
		PartnerID:       id, // Using the same ID for now
		PartnerCode:     paymentMethod.PaymentMethodCode,
		PartnerMethodID: paymentMethod.Name,
		Type:            models.PaymentMethodTypeCard, // Default to card type
		Status:          models.PaymentMethodStatus(paymentMethod.Status),
		CreatedAt:       &createdAt,
		UpdatedAt:       &updatedAt,
	}, nil
}

// MapDomainToProtoPaymentMethod maps a domain payment method to a proto payment method
func MapDomainToProtoPaymentMethod(paymentMethod *models.PaymentMethod) *cart.PaymentMethod {
	if paymentMethod == nil {
		return nil
	}

	return &cart.PaymentMethod{
		Id:                string(paymentMethod.ID),
		PaymentMethodCode: paymentMethod.PartnerCode,
		Name:              paymentMethod.PartnerMethodID,
		Description:       "", // Not used in domain model
		Status:            string(paymentMethod.Status),
		CreatedAt:         paymentMethod.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         paymentMethod.UpdatedAt.Format(time.RFC3339),
		DeletedAt:         "", // Not used in domain model
	}
}
