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
			"categories",
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
		Where("parent_id = ?", 0)

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
			"categories",
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
				"categories",
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
			"categories",
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
			"slow category page query detected",

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
				"categories",
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

	logger.Log.Info(
		"get subcategory page repository called",

		zap.Uint64(
			"parent_ID",
			parentID,
		),

		zap.String(
			"type",
			"subcategories",
		),
	)

	var categories []model.Category

	query := database.DB.
		Table("categories").
		Where("parent_id = ?", parentID).
		Where("parent_id != ?", 0).
		Order("sort_order asc")

	logger.Log.Debug(
		"executing subcategory page database query",

		zap.String(
			"type",
			"subcategories",
		),

		zap.Uint64(
			"parent_ID",
			parentID,
		),
	)

	start := time.Now()

	err := query.Find(&categories).Error

	duration := time.Since(start)

	if err != nil {
		logger.Log.Error(
			"failed to fetch subcategory page from database",

			zap.String(
				"type",
				"subcategories",
			),

			zap.Uint64(
				"parent_ID",
				parentID,
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
		"subcategory page fetched successfully",

		zap.String(
			"type",
			"subcategories",
		),

		zap.Uint64(
			"parent_ID",
			parentID,
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
			"slow subcategory page query detected",

			zap.String(
				"type",
				"subcategories",
			),

			zap.Uint64(
				"parent_ID",
				parentID,
			),

			zap.Duration(
				"duration",
				duration,
			),
		)
	}

	return categories, nil
}

func GetProductsByCategory(
	categoryID uint64,
) ([]model.Product, error) {

	logger.Log.Info(
		"get product page repository called",

		zap.Uint64(
			"category_ID",
			categoryID,
		),

		zap.String(
			"type",
			"products",
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
			categoryID,
		).
		Where(
			"p.status = ?",
			"published",
		).
		Order("p.id DESC").
		Limit(20)

	logger.Log.Debug(
		"executing product page database query",

		zap.String(
			"type",
			"products",
		),

		zap.Uint64(
			"category_ID",
			categoryID,
		),
	)

	err := query.Find(&products).Error

	duration := time.Since(start)

	if err != nil {
		logger.Log.Error(
			"failed to fetch product page from database",

			zap.String(
				"type",
				"products",
			),

			zap.Uint64(
				"category_ID",
				categoryID,
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
		"product page fetched successfully",

		zap.String(
			"type",
			"products",
		),

		zap.Uint64(
			"category_ID",
			categoryID,
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
			"slow product page query detected",

			zap.String(
				"type",
				"products",
			),

			zap.Uint64(
				"category_ID",
				categoryID,
			),

			zap.Duration(
				"duration",
				duration,
			),
		)
	}

	return products, nil
}
