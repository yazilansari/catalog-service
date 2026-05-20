package handler

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/seo/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetSEOPage(
	c *fiber.Ctx,
) error {

	tenantCode :=
		c.Locals("tenant_code").(string)

	countryCode :=
		c.Locals("country_code").(string)

	requestID, _ :=
		c.Locals("request_id").(string)

	entityType :=
		c.Params("entity")

	slug :=
		c.Params("slug")

	logger.Log.Info(
		"get seo request received",

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
			"entity_type",
			entityType,
		),

		zap.String(
			"slug",
			slug,
		),
	)

	data, err := service.GetSEOPage(
		tenantCode,
		countryCode,
		entityType,
		slug,
	)

	if err != nil {

		logger.Log.Error(
			"failed to fetch seo",

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
				"entity_type",
				entityType,
			),

			zap.String(
				"slug",
				slug,
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
		"seo fetched successfully",

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
			"entity_type",
			entityType,
		),

		zap.String(
			"slug",
			slug,
		),
	)

	return c.JSON(data)
}
