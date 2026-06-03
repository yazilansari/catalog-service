package route

import (
	"catalog-service/internal/product_search/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupProductSearchRoutes(
	api fiber.Router,
) {

	api.Get("/products/search", handler.SearchProducts)
}
