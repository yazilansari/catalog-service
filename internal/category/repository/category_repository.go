package repository

import (
	"time"

	"catalog-service/internal/category/model"
	"catalog-service/internal/database"
	"catalog-service/internal/logger"

	"go.uber.org/zap"
)

func GetCategories(
	tenantCode string,
	countryCode string,
) ([]model.Category, error) {

	logger.Log.Info(
		"get categories repository called",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),
	)

	var categories []model.Category

	query := database.DB.
		Where("tenant_code = ?", tenantCode).
		Where("country_code = ?", countryCode).
		Where("status = ?", "published").
		Order("sort_order asc")

	logger.Log.Debug(
		"executing category database query",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),
	)

	start := time.Now()

	err := query.Find(&categories).Error

	duration := time.Since(start)

	if err != nil {
		logger.Log.Error(
			"failed to fetch categories from database",

			zap.String(
				"tenant_code",
				tenantCode,
			),

			zap.String(
				"country_code",
				countryCode,
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
		"categories fetched successfully",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),

		zap.Int(
			"count",
			len(categories),
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
			"slow category query detected",

			zap.String(
				"tenant_code",
				tenantCode,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.Duration(
				"duration",
				duration,
			),
		)
	}

	return categories, nil
}
