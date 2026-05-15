package routes

import (
	"catalog-service/internal/category/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupCategoryRoutes(
	app fiber.Router,
) {

	group := app.Group("/categories")

	group.Get(
		"/",
		handler.GetCategories,
	)
}
