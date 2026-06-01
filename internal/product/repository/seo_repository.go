package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/product/model"
	"time"

	"go.uber.org/zap"
)

func GetProductSEO(
	slug string,
) (*model.SEO, error) {

	logger.Log.Info(
		"fetching seo",

		zap.String("slug", slug),
	)

	var seo model.SEO

	start := time.Now()

	query := database.DB.
		Table("seo_pages").
		Where("slug = ?", slug).
		Where("page_type = ?", "product")

	err := query.First(&seo).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch seo",

			zap.Error(err),
			zap.Duration("duration", duration),
		)
		return nil, err
	}

	logger.Log.Info(
		"seo fetched successfully",

		zap.String("slug", slug),
		zap.Duration("duration", duration),
	)

	if duration > time.Second {

		logger.Log.Warn(
			"slow product query detected",

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
