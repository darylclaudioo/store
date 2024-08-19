package service

import (
	"context"
	"hexagon-architecture/internal/domain/products/entity"
	"hexagon-architecture/internal/infrastructure"
	"hexagon-architecture/internal/utils"
	"net/http"
)

type Products struct {
	ID     string `json:"id"`
	Name   string `json:"name_product"`
	Stock  int `json:"stock"`
}

type GetProductsRequest struct {
	Page  int
	Limit int
	Name  string
}

func (s *service) GetProduct(ctx context.Context, id string) (*Products, int, error) {
	_, span := infrastructure.Tracer().Start(ctx, "service:GetProduct")
	defer span.End()

	product, code, err := s.products.GetProductByID(ctx, id)
	if err != nil {

		return nil, code, err
	}

	return &Products{
		ID:     product.ID,
		Name:   product.Name,
		Stock:   product.Stock,
	}, code, nil
}

func (s *service) GetProducts(ctx context.Context, req GetProductsRequest) ([]*Products, *utils.Pagination, int, error) {
	_, span := infrastructure.Tracer().Start(ctx, "service:GetCompanies")
	defer span.End()

	products, pagination, code, err := s.products.GetCompanies(ctx, entity.GetProductsRequest{
		Page:  req.Page,
		Limit: req.Limit,
		Name:  req.Name,
	})
	if err != nil {
		return nil, nil, code, err
	}

	productsDTO := make([]*Products, len(products))
	for i, product := range products {
		productDTO := &Products{
			ID:         product.ID,
			Name:       product.Name,
			Stock: product.Stock,
		}
		productsDTO[i] = productDTO
	}

	return productsDTO, pagination, code, nil
}

type CreateProductRequest struct {
	Name string `json:"name_product" validate:"required"`
	Stock int `json:"stock" validate:"required"`
}

func (s *service) CreateProduct(ctx context.Context, data CreateProductRequest) (*Products, int, error) {
	_, span := infrastructure.Tracer().Start(ctx, "service:CreateProduct")
	defer span.End()

	if err := utils.Validate(&data); err != nil {
		return nil, http.StatusBadRequest, err
	}

	product, code, err := s.products.CreateProduct(ctx, entity.Products{
		Name: data.Name,
		Stock : data.Stock,
	})
	if err != nil {
		return nil, code, err
	}

	return &Products{
		ID:   product.ID,
		Name: product.Name,
		Stock: product.Stock,
	}, code, nil
}

type UpdateProductRequest struct {
	Name  string `json:"name_product"`
	Stock *int   `json:"stock"`
}

func (s *service) UpdateProduct(ctx context.Context, id string, updateData UpdateProductRequest) (*Products, int, error) {
	_, span := infrastructure.Tracer().Start(ctx, "service:UpdateProduct")
	defer span.End()

	product, code, err := s.products.UpdateProduct(ctx, id, entity.UpdateProductsRequest{
		Name:  updateData.Name,
		Stock: updateData.Stock,
	})
	if err != nil {
		return nil, code, err
	}

	return &Products{
		ID:    product.ID,
		Name:  product.Name,
		Stock: product.Stock,
	}, code, nil
}

func (s *service) DeleteProduct(ctx context.Context, id string) (int, error) {
	_, span := infrastructure.Tracer().Start(ctx, "service:DeleteProduct")
	defer span.End()

	code, err := s.products.DeleteProduct(ctx, id)
	if err != nil {
		return code, err
	}

	return http.StatusOK, nil
}