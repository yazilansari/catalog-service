package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"time"

	"go.uber.org/zap"
)

type ProductElasticDocument struct {
	ID uint64 `gorm:"column:id" json:"id"`

	Name string `gorm:"column:name" json:"name"`

	Slug string `gorm:"column:slug" json:"slug"`

	Category string `gorm:"column:category" json:"category"`

	SubCategory string `gorm:"column:subcategory" json:"subcategory"`

	Brand string `gorm:"column:brand" json:"brand"`

	Price float64 `gorm:"column:price" json:"price"`

	// DiscountPrice float64 `gorm:"column:discount_price" json:"discount_price"`

	Status string `gorm:"column:status" json:"status"`

	CreatedAt string `gorm:"column:created_at" json:"created_at"`
}

func GetProductsForElastic(
	tenantCode string,
	countryCode string,
) (
	[]ProductElasticDocument,
	error,
) {

	logger.Log.Info("fetching products for elasticsearch seed")

	var products []ProductElasticDocument

	start := time.Now()

	query :=
		database.DB.
			Table("products p").
			Select(`
			p.id,
			p.name,
			p.slug,
			p.price,
			p.sale_price AS discount_price,
			p.status,
			p.created_at,

			MAX(
				CASE
					WHEN c.parent_id = 0
					THEN c.slug
				END
			) AS category,

			MAX(
				CASE
					WHEN c.parent_id <> 0
					THEN c.slug
				END
			) AS subcategory,

			b.slug AS brand
		`).
			Joins(`
			LEFT JOIN product_categories pc
			ON pc.product_id = p.id
		`).
			Joins(`
			LEFT JOIN categories c
			ON c.id = pc.category_id
		`).
			Joins(`
			LEFT JOIN brands b
			ON b.id = p.brand_id
		`).
			Where(
				"p.status = ?",
				"published",
			).
			Where(
				"p.tenant_code = ?",
				tenantCode,
			).
			Where(
				"p.country_code = ?",
				countryCode,
			).
			Group(`
			p.id,
			p.name,
			p.slug,
			p.price,
			p.sale_price,
			p.status,
			p.created_at,
			b.slug
		`)

	err := query.Find(&products).Error

	duration := time.Since(start)

	if err != nil {

		logger.Log.Error(
			"failed to fetch products for elasticsearch seed",

			zap.Error(err),

			zap.Duration(
				"duration",
				duration,
			),
		)
		return nil, err
	}

	logger.Log.Info(
		"products fetched for elasticsearch seed",

		zap.Duration(
			"duration",
			duration,
		),
	)

	return products, nil
}
