package api

import (
	"hexagon-architecture/internal/infrastructure"
	"hexagon-architecture/internal/service"
	"hexagon-architecture/internal/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type API struct {
	service service.Service
	env     string
}

func New(service service.Service, env string) *API {
	return &API{
		service: service,
		env:     env,
	}
}

func (api *API) Register(r *fiber.App) {
	r.Route("/", func(router fiber.Router) {
		router.Get("/", api.handleRoot)
		router.Get("/ping", api.handlePing)
		router.Get("/favicon.ico", api.handleFavIcon)

		router.Get("/product/:id", api.handleGetProduct)
		router.Get("/products", api.handleGetProducts)
		router.Post("/product", api.handleCreateProduct)
		router.Put("/product/:id", api.handleUpdateProduct)
		router.Delete("/product/:id", api.handleDeleteProduct)
	})
}

func (api *API) handleRoot(c *fiber.Ctx) error {
	utils.ResponseWithJSON(c, http.StatusOK, "ok", nil)
	return nil
}

func (api *API) handlePing(c *fiber.Ctx) error {
	_, span := infrastructure.Tracer().Start(c.UserContext(), "ping:handlePing")
	defer span.End()

	utils.ResponseWithJSON(c, http.StatusOK, "pong", nil)
	return nil
}

func (api *API) handleNotFound(c *fiber.Ctx) error {
	utils.ResponseWithJSON(c, http.StatusNotFound, nil, nil)
	return nil
}

func (api *API) handleFavIcon(c *fiber.Ctx) error {
	utils.ResponseWithJSON(c, http.StatusOK, "ok", nil)
	return nil
}
