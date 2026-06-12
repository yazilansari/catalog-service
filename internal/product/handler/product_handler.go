package handler

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/product/service"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetProductPage(
	c *fiber.Ctx,
) error {

	logger.Log.Info("getting product page")

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

	logger.Log.Info("getting product page",
		zap.String("slug", slug),
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
	)

	response, err :=
		service.GetProductPage(
			tenantCode,
			countryCode,
			slug,
		)

	if err != nil {

		logger.Log.Error("product not found", zap.Error(err))

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

func GetProductSnapshot(
	c *fiber.Ctx,
) error {

	logger.Log.Info("getting product snapshot")

	tenant :=
		c.Locals(
			"tenant_code",
		).(string)

	country :=
		c.Locals(
			"country_code",
		).(string)

	id,
		err :=
		c.ParamsInt(
			"id",
		)

	logger.Log.Info("getting product snapshot",
		zap.Int("id", id),
		zap.String("tenant", tenant),
		zap.String("country", country),
	)

	if err != nil {

		logger.Log.Error("invalid product id", zap.Error(err))

		return c.Status(
			400,
		).JSON(
			fiber.Map{
				"message": "invalid product id",
			},
		)
	}

	response,
		err :=
		service.GetProductSnapshot(
			tenant,
			country,
			uint64(id),
		)

	if err != nil {

		logger.Log.Error("product not found", zap.Error(err))

		return c.Status(
			404,
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
