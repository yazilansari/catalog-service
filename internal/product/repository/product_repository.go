package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/product/dto"
	"catalog-service/internal/product/model"
	"time"

	"go.uber.org/zap"
)

func GetProductBySlug(
	tenantCode string,
	countryCode string,
	slug string,
) (*model.Product, error) {

	logger.Log.Info(
		"fetching product",

		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.String("slug", slug),
	)

	var product model.Product

	start := time.Now()

	query := database.DB.
		Table("products").
		Where("tenant_code = ?", tenantCode).
		Where("country_code = ?", countryCode).
		Where("slug = ?", slug).
		Where("status = ?", "published")

	err := query.First(&product).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch product",

			zap.Error(err),
			zap.Duration("duration", duration),
		)

		return nil, err
	}

	logger.Log.Info(
		"product fetched successfully",

		zap.Uint64("product_id", product.ID),
		zap.Duration("duration", duration),
	)

	// =========================
	// SLOW QUERY DETECTION
	// =========================

	if duration > time.Second {

		logger.Log.Warn(
			"slow product query detected",

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
				"duration",
				duration,
			),
		)
	}

	return &product, nil
}

func GetProductByID(
	tenantCode string,
	countryCode string,
	productID uint64,
) (
	*dto.ProductSnapshotResponse,
	error,
) {

	logger.Log.Info(
		"fetching product by id",
		zap.Uint64("product_id", productID),
	)

	var product dto.ProductSnapshotResponse

	err :=
		database.DB.
			Table(
				"products p",
			).
			Select(`
			p.id,
			p.name,
			p.slug,
			p.sku,
			p.image,
			p.price,
			p.quantity AS stock
		`).
			Where(
				`
			p.id = ?
			AND p.tenant_code = ?
			AND p.country_code = ?
			AND p.status = ?
			`,
				productID,
				tenantCode,
				countryCode,
				"published",
			).
			Take(
				&product,
			).Error

	if err != nil {

		logger.Log.Error(
			"failed to fetch product by id",
			zap.Error(err),
		)

		return nil,
			err
	}

	return &product,
		nil

}
