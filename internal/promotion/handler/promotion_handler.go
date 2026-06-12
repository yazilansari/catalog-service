package handler

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/promotion/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetPromotionSnapshot(
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

	id,
		err :=
		c.ParamsInt(
			"id",
		)

	price, _ := strconv.ParseFloat(
		c.Query("price", "0"),
		64,
	)

	logger.Log.Info("getting product page",
		zap.Int("id", id),
		zap.String("tenant", tenantCode),
		zap.String("country", countryCode),
		zap.Float64("price", price),
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

	response, err :=
		service.GetProductPromotions(
			tenantCode,
			countryCode,
			uint64(id),
			price,
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
