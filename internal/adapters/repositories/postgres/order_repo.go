package repo

import (
	"context"

	"github.com/sean-miningah/sil-backend-assessment/internal/core/domain"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	ctx, span := otel.Tracer("").Start(ctx, "OrderRepository.Create")
	defer span.End()

	err := r.db.WithContext(ctx).Create(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (r *OrderRepository) List(ctx context.Context) ([]domain.Order, error) {
	ctx, span := otel.Tracer("").Start(ctx, "OrderRepository.List")
	defer span.End()

	var orders []domain.Order
	err := r.db.WithContext(ctx).Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) Get(ctx context.Context, id uint) (*domain.Order, error) {
	ctx, span := otel.Tracer("").Start(ctx, "OrderRepository.Get")
	defer span.End()

	var order domain.Order
	err := r.db.WithContext(ctx).First(&order, id).Error
	return &order, err
}

func (r *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	ctx, span := otel.Tracer("").Start(ctx, "OrderRepository.Update")
	defer span.End()

	return r.db.WithContext(ctx).Save(order).Error
}

func (r *OrderRepository) Delete(ctx context.Context, id uint) error {
	ctx, span := otel.Tracer("").Start(ctx, "OrderRepository.Delete")
	defer span.End()

	return r.db.WithContext(ctx).Delete(&domain.Order{}, id).Error
}

func (r *OrderRepository) DeleteOrderItems(ctx context.Context, orderID uint) error {
	ctx, span := otel.Tracer("").Start(ctx, "OrderRepository.DeleteOrderItems")
	defer span.End()

	return r.db.WithContext(ctx).Where("order_id = ?", orderID).Delete(&domain.OrderItem{}).Error
}

func (r *OrderRepository) GetOrderProduct(ctx context.Context, productID uint) (*domain.Product, error) {
	ctx, span := otel.Tracer("").Start(ctx, "OrderRepository.GetOrderProduct")
	defer span.End()

	var product domain.Product
	err := r.db.WithContext(ctx).
		First(&product, productID).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}
