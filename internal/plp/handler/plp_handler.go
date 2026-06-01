package handler

import (
	"catalog-service/internal/plp/dto"
	"catalog-service/internal/plp/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetProducts(
	c *fiber.Ctx,
) error {

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
