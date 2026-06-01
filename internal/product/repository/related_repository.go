package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/product/model"
	"time"

	"go.uber.org/zap"
)

func GetRelatedProducts(
	categoryID uint64,
	productID uint64,
) ([]model.Product, error) {

	logger.Log.Info(
		"fetching related products",

		zap.Uint64("product_id", productID),
	)

	var products []model.Product

	start := time.Now()

	query := database.DB.
		Table("products p").
		Select("DISTINCT p.*").
		Joins(`
			INNER JOIN product_categories pc
			ON pc.product_id = p.id
		`).
		Where("pc.category_id = ?", categoryID).
		Where("p.id != ?", productID).
		Where("p.status = ?", "published").
		Order("p.id DESC").
		Limit(10)

	err := query.Find(&products).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch related products",

			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return nil, err
	}

	logger.Log.Info(
		"related products fetched successfully",

		zap.Uint64("product_id", productID),
		zap.Duration("duration", duration),
	)

	if duration > time.Second {

		logger.Log.Warn(
			"slow related products query detected",

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

	return products, nil
}
