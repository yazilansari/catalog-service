package handler

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/plp/dto"
	"catalog-service/internal/plp/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func GetProducts(
	c *fiber.Ctx,
) error {

	tenantCode :=
		c.Locals("tenant_code").(string)

	countryCode :=
		c.Locals("country_code").(string)

	requestID, _ :=
		c.Locals("request_id").(string)

	logger.Log.Info(
		"get products request received",

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
	)

	limit, _ :=
		strconv.Atoi(
			c.Query(
				"limit",
				"20",
			),
		)

	minPrice, _ :=
		strconv.ParseFloat(
			c.Query(
				"min_price",
				"0",
			),
			64,
		)

	maxPrice, _ :=
		strconv.ParseFloat(
			c.Query(
				"max_price",
				"0",
			),
			64,
		)

	query :=
		dto.ProductQuery{
			Category: c.Query(
				"category",
			),

			SubCategory: c.Query(
				"subcategory",
			),

			Brand: c.Query(
				"brand",
			),

			Sort: c.Query(
				"sort",
				"latest",
			),

			Cursor: c.Query(
				"cursor",
			),

			Limit: limit,

			MinPrice: minPrice,

			MaxPrice: maxPrice,
		}

	response, err :=
		service.GetProducts(
			query,
			tenantCode,
			countryCode,
		)

	if err != nil {

		return c.Status(
			500,
		).JSON(
			fiber.Map{
				"message": err.Error(),
			},
		)
	}

	return c.JSON(
		response,
	)
}
