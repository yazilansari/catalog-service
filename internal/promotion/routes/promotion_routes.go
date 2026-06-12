package routes

import (
	"catalog-service/internal/promotion/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupPromotionRoutes(
	api fiber.Router,
) {

	product :=
		api.Group(
			"/promotions",
		)

	product.Get(
		"/:id",
		handler.GetPromotionSnapshot,
	)
}
