package utils

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data" swaggertype:"object"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Total       int `json:"total"`
	Limit       int `json:"limit"`
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
}

func ResponseWithJSON(c *fiber.Ctx, code int, data interface{}, err error, pagination ...*Pagination) {
	r := Response{
		Status:  code,
		Message: strings.ToLower(http.StatusText(code)),
	}
	if len(pagination) > 0 && pagination[0] != nil {
		r.Pagination = pagination[0]
		if r.Pagination.CurrentPage <= 0 {
			r.Pagination.CurrentPage = 1
		}
		r.Pagination.LastPage = RoundUp(float64(r.Pagination.Total) / float64(r.Pagination.Limit))
		if r.Pagination.LastPage <= 0 {
			r.Pagination.LastPage = 1
		}
	}

	r.Data = data
	if err != nil {
		r.Message = err.Error()
	}

	c.Accepts("application/json")

	_ = c.JSON(r)
}
