package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/product/model"
	"time"

	"go.uber.org/zap"
)

func GetVariantsByProductID(
	productID uint64,
) ([]model.ProductVariant, error) {

	logger.Log.Info(
		"fetching variants",

		zap.Uint64("product_id", productID),
	)

	var variants []model.ProductVariant

	start := time.Now()

	query := database.DB.
		Table("product_variants").
		Where("product_id = ?", productID).
		Where("status = ?", "published").
		Order("id asc")

	err := query.Find(&variants).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch variants",

			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return nil, err
	}

	logger.Log.Info(
		"variants fetched successfully",

		zap.Uint64("product_id", productID),
		zap.Duration("duration", duration),
	)

	if duration > time.Second {

		logger.Log.Warn(
			"slow variants query detected",

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

	return variants, nil
}
