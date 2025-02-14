package repo

import (
	"context"
	"errors"

	"github.com/sean-miningah/sil-backend-assessment/internal/core/domain"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepoisotory(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
}

func (r *CustomerRepository) UpsertCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	ctx, span := otel.Tracer("").Start(ctx, "CustomerRepository.Upsert")
	defer span.End()

	var existingCustomer domain.Customer
	if err := r.db.WithContext(ctx).Where("id = ?", customer.ID).First(&existingCustomer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Customer doesn't exist, create a new one
			if err := r.db.WithContext(ctx).Create(customer).Error; err != nil {
				return nil, err
			}
			return customer, nil
		}
		// Other errors (e.g., DB error)
		return nil, err
	}

	// Update only the fields that are provided, excluding ID
	if err := r.db.WithContext(ctx).Model(&existingCustomer).Updates(map[string]interface{}{
		"name":  customer.Name,
		"email": customer.Email,
		// Add other fields to update here
	}).Error; err != nil {
		return nil, err
	}

	return &existingCustomer, nil
}

func (r *CustomerRepository) CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	ctx, span := otel.Tracer("").Start(ctx, "CustomerRepository.Create")
	defer span.End()

	// Create the customer in the database
	result := r.db.WithContext(ctx).Create(customer)
	if result.Error != nil {
		// Log the error and return it
		span.RecordError(result.Error)
		return nil, result.Error
	}

	// Retrieve the created customer from the database
	var createdCustomer domain.Customer
	if err := r.db.WithContext(ctx).First(&createdCustomer, "id = ?", customer.ID).Error; err != nil {
		// Log the error and return it
		span.RecordError(err)
		return nil, err
	}

	return &createdCustomer, nil
}

func (r *CustomerRepository) GetCustomer(ctx context.Context, id uint) (*domain.Customer, error) {
	ctx, span := otel.Tracer("").Start(ctx, "CustomerRepository.Get")
	defer span.End()

	var customer domain.Customer
	err := r.db.WithContext(ctx).First(&customer, id).Error
	return &customer, err
}
