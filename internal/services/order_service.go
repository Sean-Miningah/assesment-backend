package services

import (
	"context"
	"fmt"

	"github.com/sean-miningah/sil-backend-assessment/internal/core/domain"
	"github.com/sean-miningah/sil-backend-assessment/internal/core/ports"
	"go.opentelemetry.io/otel"
)

type orderService struct {
	orderRepo        ports.OrderRepository
	productRepo      ports.ProductRepository
	notificationRepo ports.NotificationRepository
}

func NewOrderService(orderRepo ports.OrderRepository, productRepo ports.ProductRepository, notificationRepo ports.NotificationRepository) ports.OrderService {
	return &orderService{
		orderRepo:        orderRepo,
		productRepo:      productRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, order *domain.Order) error {
	ctx, span := otel.Tracer("").Start(ctx, "OrderService.CreateOrder")
	defer span.End()

	email, ok := ctx.Value("email").(string)
	if !ok || email == "" {
		return fmt.Errorf("email not found in context")
	}

	order, err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return err
	}
	// Send message to admin of order creation and send message to customer
	message := fmt.Sprintf("New Order Created!\n\nOrder ID: %d\nTotal Price: %.2f", order.ID, order.TotalPrice)

	err = s.notificationRepo.SendEmail(ctx, message, "New Order Purchased", email)
	if err != nil {
		return err
	}

	// Send customer message notifying of orde creation
	// err = s.notificationRepo.SendSms(ctx, []string{"+254700000000"}, "Sil", "Order has been created")
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (s *orderService) ListOrders(ctx context.Context) ([]domain.Order, error) {
	ctx, span := otel.Tracer("").Start(ctx, "OrderService.ListOrders")
	defer span.End()

	return s.orderRepo.List(ctx)
}

func (s *orderService) GetOrder(ctx context.Context, id uint) (*domain.Order, error) {
	ctx, span := otel.Tracer("").Start(ctx, "OrderService.GetOrder")
	defer span.End()

	return s.orderRepo.Get(ctx, id)
}

func (s *orderService) UpdateOrder(ctx context.Context, order *domain.Order) error {
	ctx, span := otel.Tracer("").Start(ctx, "OrderService.UpdateOrder")
	defer span.End()

	return s.orderRepo.Update(ctx, order)
}

// create a delete order route that delete orders items forthat order then it deletes the order after
func (s *orderService) DeleteOrder(ctx context.Context, id uint) error {
	ctx, span := otel.Tracer("").Start(ctx, "OrderService.DeleteOrder")
	defer span.End()

	return s.orderRepo.Delete(ctx, id)
}

func (s *orderService) GetOrderProduct(ctx context.Context, productID uint) (*domain.Product, error) {
	ctx, span := otel.Tracer("").Start(ctx, "OrderService.GetOrderProduct")
	defer span.End()

	return s.orderRepo.GetOrderProduct(ctx, productID)
}
