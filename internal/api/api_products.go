package api

import (
	"encoding/json"
	"hexagon-architecture/internal/infrastructure"
	"hexagon-architecture/internal/service"
	"hexagon-architecture/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (api *API) handleGetProduct(c *fiber.Ctx) error {
	ctx, span := infrastructure.Tracer().Start(c.UserContext(), "ping:handleGetProduct")
	defer span.End()

	id := c.Params("id")

	result, code, err := api.service.GetProduct(ctx, id)

	utils.ResponseWithJSON(c, code, result, err, nil)
	return nil
}

func (api *API) handleGetProducts(c *fiber.Ctx) error {
	ctx, span := infrastructure.Tracer().Start(c.UserContext(), "ping:handleGetProducts")
	defer span.End()

	name := c.Query("name")

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	request := service.GetProductsRequest{
		Page:  page,
		Limit: limit,
		Name:  name,
	}

	result, pagination, code, err := api.service.GetProducts(ctx, request)

	utils.ResponseWithJSON(c, code, result, err, pagination)
	return nil
}

func (api *API) handleCreateProduct(c *fiber.Ctx) error {
	ctx, span := infrastructure.Tracer().Start(c.UserContext(), "ping:handleCreateProduct")
	defer span.End()

	var request service.CreateProductRequest
	_ = json.Unmarshal(c.Body(), &request)

	result, code, err := api.service.CreateProduct(ctx, request)

	utils.ResponseWithJSON(c, code, result, err, nil)
	return nil
}

func (api *API) handleUpdateProduct(c *fiber.Ctx) error {
	ctx, span := infrastructure.Tracer().Start(c.UserContext(), "ping:handleUpdateProduct")
	defer span.End()

	id := c.Params("id")

	var updateRequest service.UpdateProductRequest
	_ = json.Unmarshal(c.Body(), &updateRequest)

	result, code, err := api.service.UpdateProduct(ctx, id, updateRequest)

	utils.ResponseWithJSON(c, code, result, err, nil)
	return nil
}

func (api *API) handleDeleteProduct(c *fiber.Ctx) error {
	ctx, span := infrastructure.Tracer().Start(c.UserContext(), "ping:handleDeleteProduct")
	defer span.End()

	id := c.Params("id")

	code, err := api.service.DeleteProduct(ctx, id)
	if err != nil {
		return c.Status(code).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(code).JSON(fiber.Map{"message": "Product deleted successfully"})
}
