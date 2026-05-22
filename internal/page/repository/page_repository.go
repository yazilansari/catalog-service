package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/page/model"
	"time"

	"go.uber.org/zap"
)

func FindCategoryBySlug(
	tenantCode string,
	countryCode string,
	slug string,
) (*model.Category, error) {

	logger.Log.Info(
		"get category page repository called",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),

		zap.String(
			"type",
			"category",
		),

		zap.String(
			"slug",
			slug,
		),
	)

	var category model.Category

	query := database.DB.
		Table("categories").
		Where("tenant_code = ?", tenantCode).
		Where("country_code = ?", countryCode).
		Where("slug = ?", slug).
		Where("parent_id IS NULL")

	logger.Log.Debug(
		"executing category page database query",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),

		zap.String(
			"type",
			"category",
		),

		zap.String(
			"slug",
			slug,
		),
	)

	start := time.Now()

	err := query.Find(&category).Error

	duration := time.Since(start)

	if err != nil {
		logger.Log.Error(
			"failed to fetch category page from database",

			zap.String(
				"tenant_code",
				tenantCode,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"type",
				"category",
			),

			zap.String(
				"slug",
				slug,
			),

			zap.Duration(
				"query_duration",
				duration,
			),

			zap.Error(err),
		)
		return nil, err
	}

	logger.Log.Info(
		"category page fetched successfully",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),

		zap.String(
			"type",
			"category",
		),

		zap.String(
			"slug",
			slug,
		),

		zap.Duration(
			"query_duration",
			duration,
		),
	)

	// =========================
	// SLOW QUERY DETECTION
	// =========================

	if duration > time.Second {

		logger.Log.Warn(
			"slow seo query detected",

			zap.String(
				"tenant_code",
				tenantCode,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"type",
				"category",
			),

			zap.String(
				"slug",
				slug,
			),

			zap.Duration(
				"duration",
				duration,
			),
		)
	}

	return &category, nil
}

func GetSubCategories(
	parentID uint64,
) ([]model.Category, error) {

	var categories []model.Category

	err := database.DB.
		Table("categories").
		Where("parent_id = ?", parentID).
		Order("sort_order asc").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func GetProductsByCategory(
	categoryID uint64,
) ([]model.Product, error) {

	var products []model.Product

	err := database.DB.
		Table("products").
		Where("category_id = ?", categoryID).
		Where("status = ?", "published").
		Order("id desc").
		Limit(20).
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	return products, nil
}
