package rest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sean-miningah/sil-backend-assessment/internal/core/domain"
	"github.com/sean-miningah/sil-backend-assessment/internal/core/ports"
	"go.opentelemetry.io/otel"
)

type OrderHandler struct {
	orderService ports.OrderService
}

type OrderItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
type OrderRequest struct {
	Items []OrderItem `json:"items"`
}

type UpdateOrderRequest struct {
	Items []OrderItem `json:"items"`
}

func NewOrderHandler(os ports.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: os,
	}
}

func (h *OrderHandler) Create(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "OrderHandler.Create")
	defer span.End()

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	ctx = context.WithValue(ctx, "email", email)
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var orderItems []domain.OrderItem
	var totalOrderPrice float64
	for _, item := range req.Items {
		// Fetch product to get its price
		product, err := h.orderService.GetOrderProduct(ctx, uint(item.ProductID))
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: fmt.Sprintf("Invalid product ID: %d", item.ProductID),
			})
			return
		}

		// Calculate item total price based on quantity
		itemPrice := product.Price * float64(item.Quantity)

		orderItems = append(orderItems, domain.OrderItem{
			ProductID: uint(item.ProductID),
			Quantity:  item.Quantity,
			Price:     itemPrice, // Store the total price for this item
		})

		totalOrderPrice += itemPrice
	}

	order := &domain.Order{
		CustomerID: 1,
		Items:      orderItems,
		CreatedAt:  time.Now(),
		TotalPrice: totalOrderPrice,
	}

	err := h.orderService.CreateOrder(ctx, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Order created successfully",
		"data":    order,
	})
}

func (h *OrderHandler) List(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "OrderHandler.List")
	defer span.End()

	orders, err := h.orderService.ListOrders(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) Get(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "OrderHandler.Get")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid order ID"})
		return
	}

	order, err := h.orderService.GetOrder(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) Update(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "OrderHandler.Update")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid order ID"})
		return
	}

	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	order, err := h.orderService.GetOrder(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Order not found"})
		return
	}

	if err := h.orderService.UpdateOrder(ctx, order); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) Delete(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "OrderHandler.Delete")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid order ID"})
		return
	}

	if err := h.orderService.DeleteOrder(ctx, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete order"})
		return
	}

	if err := h.orderService.DeleteOrder(ctx, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete order"})
		return
	}

	c.Status(http.StatusNoContent)
}
