package handler

import (
	"catalog-service/internal/category/service"
	"catalog-service/internal/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetCategories(
	c *fiber.Ctx,
) error {

	tenantCode :=
		c.Locals("tenant_code").(string)

	countryCode :=
		c.Locals("country_code").(string)

	requestID, _ :=
		c.Locals("request_id").(string)

	logger.Log.Info(
		"get categories request received",

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
	)

	data, err := service.GetCategoryTree(
		tenantCode,
		countryCode,
	)

	if err != nil {

		logger.Log.Error(
			"failed to fetch category tree",

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

			zap.Error(err),
		)

		return c.Status(500).JSON(
			fiber.Map{
				"message": "failed",
			},
		)
	}

	logger.Log.Info(
		"categories fetched successfully",

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

		zap.Int(
			"root_categories",
			len(data),
		),
	)

	return c.JSON(data)
}
