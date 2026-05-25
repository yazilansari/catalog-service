package handler

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/page/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetProductPage(
	c *fiber.Ctx,
) error {

	slug := c.Params("slug")

	tenantCode :=
		c.Locals("tenant_code").(string)

	countryCode :=
		c.Locals("country_code").(string)

	requestID, _ :=
		c.Locals("request_id").(string)

	page, _ := strconv.Atoi(
		c.Query("page", "1"),
	)

	limit, _ := strconv.Atoi(
		c.Query("limit", "20"),
	)

	sort := c.Query(
		"sort",
		"newest",
	)

	brand := c.Query("brand")

	priceMin, _ := strconv.ParseFloat(
		c.Query("price_min", "0"),
		64,
	)

	priceMax, _ := strconv.ParseFloat(
		c.Query("price_max", "0"),
		64,
	)

	logger.Log.Info(
		"get product page request received",

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

		zap.Int(
			"page",
			page,
		),

		zap.Int(
			"limit",
			limit,
		),

		zap.String(
			"sort",
			sort,
		),

		zap.String(
			"brand",
			brand,
		),

		zap.Float64(
			"price_min",
			priceMin,
		),

		zap.Float64(
			"price_max",
			priceMax,
		),
	)

	data, err := service.GetProductPage(
		tenantCode,
		countryCode,
		slug,
		page,
		limit,
		sort,
		brand,
		priceMin,
		priceMax,
	)

	if err != nil {

		logger.Log.Error(
			"failed to fetch product page",

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
				"message": "product page not found",
			},
		)
	}

	logger.Log.Info(
		"product page fetched successfully",

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
