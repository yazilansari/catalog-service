package routes

import (
	"catalog-service/internal/plp/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupPLPRoutes(
	api fiber.Router,
) {

	products :=
		api.Group(
			"/products",
		)

	products.Get(
		"/",
		handler.GetProducts,
	)
}
