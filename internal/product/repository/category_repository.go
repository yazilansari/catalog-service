package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/product/model"
	"time"

	"go.uber.org/zap"
)

func GetCategoryAndSubCategory(
	productID uint64,
) (*model.Category, *model.Category, error) {

	logger.Log.Info(
		"fetching category and subcategory",

		zap.Uint64("product_id", productID),
	)

	var subCategory model.Category

	start := time.Now()

	query := database.DB.
		Table("categories c").
		Select("c.*").
		Joins(`
			INNER JOIN product_categories pc
			ON pc.category_id = c.id
		`).
		Where("pc.product_id = ?", productID).
		Where("c.parent_id != ?", 0).
		Limit(1)

	err := query.First(&subCategory).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch subcategory",

			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return nil, nil, err
	}

	var category model.Category

	query = database.DB.
		Table("categories").
		Where("id = ?", subCategory.ParentID)

	err = query.First(&category).Error

	if err != nil {

		logger.Log.Error(
			"failed to fetch category",

			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return nil, nil, err
	}

	logger.Log.Info(
		"product fetched successfully",

		zap.Uint64("product_id", productID),
		zap.Duration("duration", duration),
	)

	// =========================
	// SLOW QUERY DETECTION
	// =========================

	if duration > time.Second {

		logger.Log.Warn(
			"slow category and subcategory query detected",

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

	return &category, &subCategory, nil
}
