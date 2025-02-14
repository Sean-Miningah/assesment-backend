package services

import (
	"context"

	"github.com/sean-miningah/sil-backend-assessment/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	args := m.Called(ctx, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) List(ctx context.Context) ([]domain.Order, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Order), args.Error(1)
}

func (m *MockOrderRepository) Get(ctx context.Context, id uint) (*domain.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) Update(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockOrderRepository) DeleteOrderItems(ctx context.Context, orderID uint) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrderProduct(ctx context.Context, productID uint) (*domain.Product, error) {
	args := m.Called(ctx, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}
