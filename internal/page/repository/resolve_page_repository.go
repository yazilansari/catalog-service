package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/page/dto"
	"time"

	"go.uber.org/zap"
)

func ResolvePage(
	tenantCode string,
	countryCode string,
	slug string,
) (*dto.ResolvePageResponse, error) {

	logger.Log.Info(
		"resolve page repository called",

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

	start := time.Now()

	var response dto.ResolvePageResponse

	// =========================
	// CATEGORY
	// =========================

	err := database.DB.
		Table("categories").
		Select(
			"'category' AS page_type",
			"slug",
			"CONCAT('/category/', slug) AS redirect_url",
		).
		Where("tenant_code = ?", tenantCode).
		Where("country_code = ?", countryCode).
		Where("slug = ?", slug).
		Where("parent_id IS NULL").
		Limit(1).
		Scan(&response).Error

	if err == nil && response.PageType != "" {

		logger.Log.Info(
			"category resolved",

			zap.String(
				"slug",
				slug,
			),
		)

		return &response, nil
	}

	// =========================
	// SUBCATEGORY
	// =========================

	err = database.DB.
		Table("categories").
		Select(
			"'subcategory' AS page_type",
			"slug",
			"CONCAT('/subcategory/', slug) AS redirect_url",
		).
		Where("tenant_code = ?", tenantCode).
		Where("country_code = ?", countryCode).
		Where("slug = ?", slug).
		Where("parent_id IS NOT NULL").
		Limit(1).
		Scan(&response).Error

	if err == nil && response.PageType != "" {

		logger.Log.Info(
			"subcategory resolved",

			zap.String(
				"slug",
				slug,
			),
		)

		return &response, nil
	}

	// =========================
	// PRODUCT
	// =========================

	err = database.DB.
		Table("products").
		Select(
			"'product' AS page_type",
			"slug",
			"CONCAT('/product/', slug) AS redirect_url",
		).
		Where("tenant_code = ?", tenantCode).
		Where("country_code = ?", countryCode).
		Where("slug = ?", slug).
		Limit(1).
		Scan(&response).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to resolve page",

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

			zap.Duration(
				"query_duration",
				duration,
			),

			zap.Error(err),
		)

		return nil, err
	}

	if response.PageType == "" {

		logger.Log.Warn(
			"page not found",

			zap.String(
				"slug",
				slug,
			),
		)

		return nil, nil
	}

	logger.Log.Info(
		"page resolved successfully",

		zap.String(
			"page_type",
			response.PageType,
		),

		zap.String(
			"slug",
			response.Slug,
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
			"slow resolve query detected",

			zap.Duration(
				"duration",
				duration,
			),

			zap.String(
				"slug",
				slug,
			),
		)
	}

	return &response, nil
}
