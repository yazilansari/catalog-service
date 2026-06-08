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

	err := query.First(&category).Error

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
	page int,
	limit int,
	sort string,
	brand string,
	priceMin float64,
	priceMax float64,
) ([]model.Product, int64, error) {

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

		zap.Int(
			"page",
			page,
		),

		zap.Int(
			"limit",
			limit,
		),

		zap.String(
			"sort",
			sort,
		),
	)

	var products []model.Product

	var total int64

	offset := (page - 1) * limit

	start := time.Now()

	query := database.DB.
		Table("products p").
		Select("DISTINCT p.*").
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
		)

	// =========================
	// FILTERS
	// =========================

	if brand != "" {

		query = query.Where(
			"LOWER(p.brand) = LOWER(?)",
			brand,
		)
	}

	if priceMin > 0 {

		query = query.Where(
			"p.price >= ?",
			priceMin,
		)
	}

	if priceMax > 0 {

		query = query.Where(
			"p.price <= ?",
			priceMax,
		)
	}

	// =========================
	// TOTAL COUNT
	// =========================

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// =========================
	// SORTING
	// =========================

	switch sort {

	case "price_asc":
		query = query.Order("p.price ASC")

	case "price_desc":
		query = query.Order("p.price DESC")

	case "oldest":
		query = query.Order("p.id ASC")

	case "name_asc":
		query = query.Order("p.name ASC")

	case "name_desc":
		query = query.Order("p.name DESC")

	default:
		query = query.Order("p.id DESC")
	}

	// =========================
	// PAGINATION
	// =========================

	query = query.
		Offset(offset).
		Limit(limit)

	err = query.Find(&products).Error

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
		return nil, 0, err
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

	return products, total, nil
}
