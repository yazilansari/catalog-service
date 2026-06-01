package handler

import (
	"catalog-service/internal/product/service"

	"github.com/gofiber/fiber/v2"
)

func GetProductPage(
	c *fiber.Ctx,
) error {

	tenantCode :=
		c.Locals(
			"tenant_code",
		).(string)

	countryCode :=
		c.Locals(
			"country_code",
		).(string)

	slug :=
		c.Params("slug")

	response, err :=
		service.GetProductPage(
			tenantCode,
			countryCode,
			slug,
		)

	if err != nil {

		return c.Status(
			fiber.StatusNotFound,
		).JSON(
			fiber.Map{
				"message": "product not found",
			},
		)
	}

	return c.JSON(
		response,
	)
}
