package repository

import (
	"time"

	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/seo/model"

	"go.uber.org/zap"
)

func GetSEO(
	tenantCode string,
	countryCode string,
	entityType string,
	slug string,
) (*model.SEOPage, error) {

	logger.Log.Info(
		"get seo repository called",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),

		zap.String(
			"entity_type",
			entityType,
		),

		zap.String(
			"slug",
			slug,
		),
	)

	var seo model.SEOPage

	query := database.DB.
		Where("tenant_code = ?", tenantCode).
		Where("country_code = ?", countryCode).
		Where("entity_type = ?", entityType).
		Where("slug = ?", slug)

	logger.Log.Debug(
		"executing seo database query",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),

		zap.String(
			"entity_type",
			entityType,
		),

		zap.String(
			"slug",
			slug,
		),
	)

	start := time.Now()

	err := query.Find(&seo).Error

	duration := time.Since(start)

	if err != nil {
		logger.Log.Error(
			"failed to fetch seo from database",

			zap.String(
				"tenant_code",
				tenantCode,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"entity_type",
				entityType,
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
		"seo fetched successfully",

		zap.String(
			"tenant_code",
			tenantCode,
		),

		zap.String(
			"country_code",
			countryCode,
		),

		zap.String(
			"entity_type",
			entityType,
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
				"entity_type",
				entityType,
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

	return &seo, nil
}
