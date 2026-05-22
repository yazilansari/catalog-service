package handler

import (
	"catalog-service/internal/page/service"

	"github.com/gofiber/fiber/v2"
)

func ResolvePage(
	c *fiber.Ctx,
) error {

	slug := c.Params("slug")

	tenantCode :=
		c.Locals("tenant_code").(string)

	countryCode :=
		c.Locals("country_code").(string)

	data, err := service.ResolvePage(
		tenantCode,
		countryCode,
		slug,
	)

	if err != nil {

		return c.Status(404).JSON(
			fiber.Map{
				"message": "page not found",
			},
		)
	}

	return c.JSON(data)
}
