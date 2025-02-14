package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sean-miningah/sil-backend-assessment/internal/core/domain"
	"github.com/sean-miningah/sil-backend-assessment/internal/core/ports"
	"go.opentelemetry.io/otel"
)

type ProductHandler struct {
	productService ports.ProductService
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=3,max=255"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	CategoryID  uint    `json:"category_id" binding:"required"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"omitempty,min=3,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	CategoryID  uint    `json:"category_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewProductHandler(ps ports.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: ps,
	}
}

// Create godoc
// @Summary Create a new product
// @Description Create a new product with the provided details
// @Tags products
// @Accept json
// @Produce json
// @Param product body CreateProductRequest true "Product details"
// @Success 201 {object} domain.Product
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "ProductHandler.Create")
	defer span.End()

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	product := &domain.Product{
		Name:       req.Name,
		Price:      req.Price,
		CategoryID: req.CategoryID,
	}

	if err := h.productService.CreateProduct(ctx, product); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// Get godoc
// @Summary Get a product by ID
// @Description Get detailed information about a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "ProductHandler.Get")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid product ID"})
		return
	}

	product, err := h.productService.GetProduct(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// List godoc
// @Summary List all products
// @Description Get a list of all products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} domain.Product
// @Failure 500 {object} ErrorResponse
// @Router /products [get]
func (h *ProductHandler) List(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "ProductHandler.List")
	defer span.End()

	products, err := h.productService.ListProducts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Update godoc
// @Summary Update a product
// @Description Update a product's details
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body UpdateProductRequest true "Product details"
// @Success 200 {object} domain.Product
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "ProductHandler.Update")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid product ID"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	product, err := h.productService.GetProduct(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Product not found"})
		return
	}

	// Update only provided fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.CategoryID > 0 {
		product.CategoryID = req.CategoryID
	}

	if err := h.productService.UpdateProduct(ctx, product); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// Delete godoc
// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	ctx, span := otel.Tracer("").Start(c.Request.Context(), "ProductHandler.Delete")
	defer span.End()

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid product ID"})
		return
	}

	if err := h.productService.DeleteProduct(ctx, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete product"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) GetAveragePriceByCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("categoryId"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid product ID"})
		return
	}

	avgPrice, err := h.productService.GetAverageCategoryPrice(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	c.JSON(http.StatusOK, avgPrice)
}
