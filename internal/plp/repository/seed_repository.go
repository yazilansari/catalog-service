package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/page/model"
	"time"

	"go.uber.org/zap"
)

func GetProductsForElastic() (
	[]model.Product,
	error,
) {

	logger.Log.Info("fetching products for elasticsearch seed")

	var products []model.Product

	start := time.Now()

	query :=
		database.DB.
			Table("products").
			Where(
				"status = ?",
				"published",
			)

	err := query.Find(&products).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch products for elasticsearch seed",

			zap.Error(err),

			zap.Duration(
				"duration",
				duration,
			),
		)
		return nil, err
	}

	logger.Log.Info(
		"products fetched for elasticsearch seed",

		zap.Duration(
			"duration",
			duration,
		),
	)

	return products, nil
}
