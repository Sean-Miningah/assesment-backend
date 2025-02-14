package ports

import (
	"context"

	"github.com/sean-miningah/sil-backend-assessment/internal/core/domain"
)

type ProductService interface {
	CreateProduct(ctx context.Context, product *domain.Product) error
	GetProduct(ctx context.Context, id uint) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product) error
	DeleteProduct(ctx context.Context, id uint) error
	GetAverageCategoryPrice(ctx context.Context, id uint) (float64, error)
}

type CategoryService interface {
	CreateCategory(ctx context.Context, category *domain.Category) error
	GetCategory(ctx context.Context, id uint) (*domain.Category, error)
	ListCategories(ctx context.Context) ([]domain.Category, error)
	UpdateCategory(ctx context.Context, category *domain.Category) error
	DeleteCategory(ctx context.Context, id uint) error
}

type OrderService interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
	ListOrders(ctx context.Context) ([]domain.Order, error)
	GetOrder(ctx context.Context, id uint) (*domain.Order, error)
	UpdateOrder(ctx context.Context, order *domain.Order) error
	DeleteOrder(ctx context.Context, id uint) error
	GetOrderProduct(ctx context.Context, productID uint) (*domain.Product, error)
}

type CustomerService interface {
	CreateCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
	GetCustomer(ctx context.Context, id uint) (*domain.Customer, error)
	UpsertCustomer(ctx context.Context, customer *domain.Customer) (*domain.Customer, error)
}
