package routes

import (
	"catalog-service/internal/page/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupPageRoutes(
	api fiber.Router,
) {

	page := api.Group("/pages")

	page.Get(
		"/category/:slug",
		handler.GetCategoryPage,
	)

	page.Get(
		"/subcategory/:slug",
		handler.GetSubCategoryPage,
	)

	page.Get(
		"/product/:slug",
		handler.GetProductPage,
	)

	page.Get(
		"/resolve/:slug",
		handler.ResolvePage,
	)
}
