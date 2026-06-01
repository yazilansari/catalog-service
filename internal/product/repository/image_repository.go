package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/product/model"
	"time"

	"go.uber.org/zap"
)

func GetProductImages(
	productID uint64,
) ([]model.ProductImage, error) {

	logger.Log.Info(
		"fetching images",

		zap.Uint64("product_id", productID),
	)

	var images []model.ProductImage

	start := time.Now()

	query := database.DB.
		Table("product_images").
		Where("product_id = ?", productID).
		Order("sort_order asc")

	err := query.Find(&images).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch images",

			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return nil, err
	}

	logger.Log.Info(
		"images fetched successfully",

		zap.Uint64("product_id", productID),
		zap.Duration("duration", duration),
	)

	if duration > time.Second {

		logger.Log.Warn(
			"slow images query detected",

			zap.Uint64(
				"product_id",
				productID,
			),

			zap.Duration(
				"duration",
				duration,
			),
		)
	}

	return images, nil
}
