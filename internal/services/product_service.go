package services

import (
	"context"

	"github.com/sean-miningah/sil-backend-assessment/internal/core/domain"
	"github.com/sean-miningah/sil-backend-assessment/internal/core/ports"
	"go.opentelemetry.io/otel"
)

type productService struct {
	productrepo ports.ProductRepository
	orderrepo   ports.OrderRepository
}

func NewProductService(product ports.ProductRepository, order ports.OrderRepository) ports.ProductService {
	return &productService{productrepo: product, orderrepo: order}
}

func (s *productService) CreateProduct(ctx context.Context, product *domain.Product) error {
	ctx, span := otel.Tracer("").Start(ctx, "ProductService.CreateProduct")
	defer span.End()

	return s.productrepo.Create(ctx, product)
}

func (s *productService) GetProduct(ctx context.Context, id uint) (*domain.Product, error) {
	ctx, span := otel.Tracer("").Start(ctx, "ProductService.GetProduct")
	defer span.End()

	return s.productrepo.Get(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context) ([]domain.Product, error) {
	ctx, span := otel.Tracer("").Start(ctx, "ProductService.ListProducts")
	defer span.End()

	return s.productrepo.List(ctx)
}

func (s *productService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	ctx, span := otel.Tracer("").Start(ctx, "ProductService.UpdateProduct")
	defer span.End()

	return s.productrepo.Update(ctx, product)
}

func (s *productService) DeleteProduct(ctx context.Context, id uint) error {
	ctx, span := otel.Tracer("").Start(ctx, "ProductService.DeleteProduct")
	defer span.End()

	return s.productrepo.Delete(ctx, id)
}

func (s *productService) GetAverageCategoryPrice(ctx context.Context, categoryID uint) (float64, error) {
	ctx, span := otel.Tracer("").Start(ctx, "ProductService.GetAverageProducgPrice")
	defer span.End()

	return s.productrepo.GetAverageCategoryPrice(ctx, categoryID)
}
