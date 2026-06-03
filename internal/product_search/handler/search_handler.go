package handler

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/product_search/dto"
	"catalog-service/internal/product_search/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func SearchProducts(
	c *fiber.Ctx,
) error {

	tenantCode :=
		c.Locals("tenant_code").(string)

	countryCode :=
		c.Locals("country_code").(string)

	logger.Log.Info(
		"search products request received",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
	)

	var query dto.ProductSearchQuery

	query.Query = c.Query("q")
	query.Category = c.Query("category")
	query.SubCategory = c.Query("subcategory")
	query.Brand = c.Query("brand")
	query.Sort = c.Query("sort")
	query.Cursor = c.Query("cursor")

	if limit := c.QueryInt("limit"); limit > 0 {
		query.Limit = limit
	}

	if minPrice, err := strconv.ParseFloat(
		c.Query("min_price"),
		64,
	); err == nil {
		query.MinPrice = minPrice
	}

	if maxPrice, err := strconv.ParseFloat(
		c.Query("max_price"),
		64,
	); err == nil {
		query.MaxPrice = maxPrice
	}

	logger.Log.Info(
		"parsed search query",
		zap.Any("query", query),
	)

	// if err != nil {

	// 	logger.Log.Error(
	// 		"invalid query params",
	// 		zap.Error(err),
	// 	)

	// 	return c.Status(400).JSON(
	// 		fiber.Map{
	// 			"message": "invalid query params",
	// 		},
	// 	)
	// }

	response, err :=
		service.SearchProducts(
			query,
			tenantCode,
			countryCode,
		)

	if err != nil {

		logger.Log.Error(
			"error searching products",
			zap.Error(err),
		)

		return c.Status(500).JSON(
			fiber.Map{
				"message": err.Error(),
			},
		)
	}

	logger.Log.Info(
		"search products request completed",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Any("cache_key", response),
	)

	return c.JSON(response)
}
