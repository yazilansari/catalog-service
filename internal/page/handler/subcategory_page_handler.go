package handler

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/page/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetSubCategoryPage(
	c *fiber.Ctx,
) error {

	slug := c.Params("slug")

	tenantCode :=
		c.Locals("tenant_code").(string)

	countryCode :=
		c.Locals("country_code").(string)

	requestID, _ :=
		c.Locals("request_id").(string)

	logger.Log.Info(
		"get subcategory page request received",

		zap.String(
			"request_id",
			requestID,
		),

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),

		zap.String(
			"path",
			c.OriginalURL(),
		),

		zap.String(
			"method",
			c.Method(),
		),

		zap.String(
			"slug",
			slug,
		),
	)

	data, err := service.GetSubCategoryPage(
		tenantCode,
		countryCode,
		slug,
	)

	if err != nil {

		logger.Log.Error(
			"failed to fetch subcategory page",

			zap.String(
				"request_id",
				requestID,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"slug",
				slug,
			),

			zap.Error(err),
		)

		return c.Status(404).JSON(
			fiber.Map{
				"message": "subcategory page not found",
			},
		)
	}

	logger.Log.Info(
		"subcategory page fetched successfully",

		zap.String(
			"request_id",
			requestID,
		),

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),

		zap.String(
			"slug",
			slug,
		),
	)

	return c.JSON(data)
}
