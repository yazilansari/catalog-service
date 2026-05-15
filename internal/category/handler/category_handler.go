package handler

import (
	"catalog-service/internal/category/service"

	"github.com/gofiber/fiber/v2"
)

func GetCategories(
	c *fiber.Ctx,
) error {

	tenantCode :=
		c.Locals("tenant_code").(string)

	countryCode :=
		c.Locals("country_code").(string)

	data, err := service.GetCategoryTree(
		tenantCode,
		countryCode,
	)

	if err != nil {

		return c.Status(500).JSON(
			fiber.Map{
				"message": "failed",
			},
		)
	}

	return c.JSON(data)
}
