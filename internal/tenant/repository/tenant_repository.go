package repository

import (
	"time"

	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/tenant/model"

	"go.uber.org/zap"
)

func FindTenantByDomain(domain string) (*model.Tenant, error) {

	logger.Log.Info(
		"find tenant by domain called",

		zap.String(
			"domain",
			domain,
		),
	)

	var tenant model.Tenant

	start := time.Now()

	err := database.DB.
		Where("domain = ?", domain).
		Where("active = ?", true).
		First(&tenant).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to find tenant by domain",

			zap.String(
				"domain",
				domain,
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
		"tenant fetched successfully",

		zap.String(
			"domain",
			domain,
		),

		zap.String(
			"tenant_code",
			tenant.TenantCode,
		),

		zap.String(
			"country_code",
			tenant.CountryCode,
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
			"slow tenant query detected",

			zap.String(
				"domain",
				domain,
			),

			zap.Duration(
				"duration",
				duration,
			),
		)
	}

	return &tenant, nil
}
