package repository

import (
	"time"

	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/page/model"

	"go.uber.org/zap"
)

func FindProductBySlug(
	tenantCode string,
	countryCode string,
	slug string,
) (*model.Product, error) {

	logger.Log.Info(
		"find product by slug",

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

	var product model.Product

	start := time.Now()

	query := database.DB.
		// Debug().
		Table("products").
		Where("tenant_code = ?", tenantCode).
		Where("country_code = ?", countryCode).
		Where("slug = ?", slug).
		Where("status = ?", "published")

	err := query.First(&product).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch product repository",

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
		"product fetched successfully",

		zap.Uint64(
			"product_id",
			product.ID,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return &product, nil
}

func GetCategoryByProductID(
	productID uint64,
) (*model.Category, error) {

	logger.Log.Info(
		"fetching category by product id",

		zap.Uint64(
			"product_id",
			productID,
		),
	)

	var category model.Category

	query := database.DB.
		Table("categories c").
		Select("c.*").
		Joins(
			"INNER JOIN product_categories pc ON pc.category_id = c.id",
		).
		Where(
			"pc.product_id = ?",
			productID,
		).
		Where(
			"c.parent_id = 0",
		).
		Limit(1)

	start := time.Now()

	err := query.First(&category).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch category",

			zap.Uint64(
				"product_id",
				productID,
			),

			zap.Error(err),
		)

		return nil, err
	}

	logger.Log.Info(
		"category fetched successfully",

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

func GetSubCategoryByID(
	productID uint64,
) (*model.Category, error) {

	logger.Log.Info(
		"fetching subcategory",

		zap.Uint64(
			"product_id",
			productID,
		),
	)

	var category model.Category

	start := time.Now()

	query := database.DB.
		Table("categories c").
		Select("c.*").
		Joins(
			"INNER JOIN product_categories pc ON pc.category_id = c.id",
		).
		Where(
			"pc.product_id = ?",
			productID,
		).
		Where(
			"c.parent_id != 0",
		).
		Limit(1)

	err := query.First(&category).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch subcategory",

			zap.Uint64(
				"product_id",
				productID,
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
			category.ID,
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return &category, nil
}

func GetRelatedProducts(
	subCategoryID uint64,
	productID uint64,
	page int,
	limit int,
	sort string,
	brand string,
	priceMin float64,
	priceMax float64,
) ([]model.Product, int64, error) {

	logger.Log.Info(
		"fetching related products",

		zap.Uint64(
			"category_id",
			subCategoryID,
		),

		zap.Uint64(
			"product_id",
			productID,
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
			subCategoryID,
		).
		Where(
			"p.id != ?",
			productID,
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

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch related products",

			zap.Uint64(
				"product_id",
				productID,
			),

			zap.Duration(
				"duration",
				duration,
			),

			zap.Error(err),
		)

		return nil, 0, err
	}

	logger.Log.Info(
		"related products fetched successfully",

		zap.Int(
			"count",
			len(products),
		),

		zap.Duration(
			"duration",
			duration,
		),
	)

	return products, total, nil
}
