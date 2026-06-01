package routes

import (
	"catalog-service/internal/product/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutes(
	api fiber.Router,
) {

	product :=
		api.Group(
			"/products",
		)

	product.Get(
		"/:slug",
		handler.GetProductPage,
	)
}
