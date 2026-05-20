package routes

import (
	"catalog-service/internal/seo/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupSEORoutes(
	api fiber.Router,
) {

	seo := api.Group("/seo")

	seo.Get(
		"/:entity/:slug",
		handler.GetSEOPage,
	)
}
