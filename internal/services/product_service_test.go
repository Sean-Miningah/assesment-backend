package services

import (
	"context"
	"testing"

	"github.com/sean-miningah/sil-backend-assessment/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Get(ctx context.Context, id uint) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) List(ctx context.Context) ([]domain.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) GetAverageCategoryPrice(ctx context.Context, productID uint) (float64, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).(float64), args.Error(1)
}

func TestProductService_CreateProduct(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockOrderRepo := new(MockOrderRepository)
	service := NewProductService(mockProductRepo, mockOrderRepo)

	product := &domain.Product{
		Name:       "Test Product",
		Price:      99.99,
		CategoryID: 1,
	}

	mockProductRepo.On("Create", mock.Anything, product).Return(nil)

	err := service.CreateProduct(context.Background(), product)
	assert.NoError(t, err)
	mockProductRepo.AssertExpectations(t)
}

func TestProductService_GetProduct(t *testing.T) {
	mockProductRepo := new(MockProductRepository)
	mockOrderRepo := new(MockOrderRepository)
	service := NewProductService(mockProductRepo, mockOrderRepo)

	expectedProduct := &domain.Product{
		ID:         1,
		Name:       "Test Product",
		Price:      99.99,
		CategoryID: 1,
	}

	mockProductRepo.On("Get", mock.Anything, uint(1)).Return(expectedProduct, nil)

	product, err := service.GetProduct(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)
	mockProductRepo.AssertExpectations(t)
}
