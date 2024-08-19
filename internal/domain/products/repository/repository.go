package repository

import (
	"context"
	"hexagon-architecture/internal/domain/products/entity"
	"hexagon-architecture/internal/utils"
)

type Repository interface {
	GetProductByID(ctx context.Context, id string) (*entity.Products, int, error)
	GetCompanies(ctx context.Context, data entity.GetProductsRequest) ([]*entity.Products, *utils.Pagination, int, error)
	CreateProduct(ctx context.Context, product entity.Products) (*entity.Products, int, error)
	UpdateProduct(ctx context.Context, id string, updateData entity.UpdateProductsRequest) (*entity.Products, int, error)
	DeleteProduct(ctx context.Context, id string) (int, error)
}