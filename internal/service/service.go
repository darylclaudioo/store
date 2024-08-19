package service

import (
	"context"
	productsRepo "hexagon-architecture/internal/domain/products/repository"
	"hexagon-architecture/internal/utils"
)

type Service interface {
	GetProduct(ctx context.Context, id string) (*Products, int, error)
	GetProducts(ctx context.Context, req GetProductsRequest) ([]*Products, *utils.Pagination, int, error)
	CreateProduct(ctx context.Context, data CreateProductRequest) (*Products, int, error)
	UpdateProduct(ctx context.Context, id string, updateData UpdateProductRequest) (*Products, int, error)
	DeleteProduct(ctx context.Context, id string) (int, error)
}

type service struct {
	products productsRepo.Repository
}

func New(
	products productsRepo.Repository,
) Service {
	return &service{
		products: products,
	}
}
