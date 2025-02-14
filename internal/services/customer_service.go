package services

import (
	"context"

	"github.com/sean-miningah/sil-backend-assessment/internal/core/domain"
	"github.com/sean-miningah/sil-backend-assessment/internal/core/ports"
	"go.opentelemetry.io/otel"
)

type CustomerService struct {
	customerRepo ports.CustomerRepository
}

func NewCustomerService(customerRepo ports.CustomerRepository) *CustomerService {
	return &CustomerService{customerRepo: customerRepo}
}

func (s *CustomerService) CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	ctx, span := otel.Tracer("").Start(ctx, "CustomerService.CreateNewCustomer")
	defer span.End()

	return s.customerRepo.CreateCustomer(ctx, customer)
}

func (s *CustomerService) GetCustomer(ctx context.Context, id uint) (*domain.Customer, error) {
	ctx, span := otel.Tracer("").Start(ctx, "CustomerService.GetCustomer")
	defer span.End()

	return s.customerRepo.GetCustomer(ctx, id)
}

func (r *CustomerService) UpsertCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	ctx, span := otel.Tracer("").Start(ctx, "CustomerService.UpsertCustomer")
	defer span.End()

	return r.customerRepo.UpsertCustomer(ctx, customer)
}
