package repository

import (
	"time"

	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/page/model"

	"go.uber.org/zap"
)

func FindSubCategoryBySlug(
	tenantCode string,
	countryCode string,
	slug string,
) (*model.Category, error) {

	logger.Log.Info(
		"find subcategory by slug",

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

	var subCategory model.Category

	start := time.Now()

	query := database.DB.
		Table("categories").
		Where("tenant_code = ?", tenantCode).
		Where("country_code = ?", countryCode).
		Where("slug = ?", slug).
		Where("parent_id != ?", 0)

	err := query.First(&subCategory).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch subcategory",

			zap.String(
				"slug",
				slug,
			),

			zap.Duration(
				"duration",
				duration,
			),

			zap.Error(err),
		)

		return nil, err
	}

	logger.Log.Info(
		"subcategory fetched successfully",

		zap.Uint64(
			"subcategory_id",
			subCategory.ID,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return &subCategory, nil
}

func GetCategoryByID(
	categoryID uint64,
) (*model.Category, error) {

	logger.Log.Info(
		"fetching parent category",

		zap.Uint64(
			"category_id",
			categoryID,
		),
	)

	var category model.Category

	start := time.Now()

	query := database.DB.
		Table("categories").
		Where("id = ?", categoryID).
		Where("parent_id = ?", 0)

	err := query.First(&category).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch parent category",

			zap.Uint64(
				"category_id",
				categoryID,
			),

			zap.Duration(
				"duration",
				duration,
			),

			zap.Error(err),
		)

		return nil, err
	}

	logger.Log.Info(
		"parent category fetched successfully",

		zap.Uint64(
			"category_id",
			category.ID,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return &category, nil
}

func GetProductsBySubCategory(
	subCategoryID uint64,
) ([]model.Product, error) {

	logger.Log.Info(
		"fetching subcategory products",

		zap.Uint64(
			"subcategory_id",
			subCategoryID,
		),
	)

	var products []model.Product

	start := time.Now()

	query := database.DB.
		Table("products p").
		Select("p.*").
		Joins(`
			INNER JOIN product_categories pc
			ON pc.product_id = p.id
		`).
		Where(
			"pc.category_id = ?",
			subCategoryID,
		).
		Where(
			"p.status = ?",
			"published",
		).
		Order("p.id DESC").
		Limit(20)

	err := query.Find(&products).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch subcategory products",

			zap.Uint64(
				"subcategory_id",
				subCategoryID,
			),

			zap.Duration(
				"duration",
				duration,
			),

			zap.Error(err),
		)

		return nil, err
	}

	logger.Log.Info(
		"subcategory products fetched successfully",

		zap.Uint64(
			"subcategory_id",
			subCategoryID,
		),

		zap.Int(
			"count",
			len(products),
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return products, nil
}
