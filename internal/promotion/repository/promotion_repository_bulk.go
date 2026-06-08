package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/promotion/model"
	"time"

	"go.uber.org/zap"
)

type ProductDiscount struct {
	ProductID    uint64             `gorm:"column:product_id"`
	Promotion    model.Promotion    `gorm:"embedded"`
	DiscountRule model.DiscountRule `gorm:"embedded"`
}

func GetActiveDiscountsByProducts(
	tenantCode string,
	countryCode string,
	productIDs []uint64,
) ([]ProductDiscount, error) {

	var results []ProductDiscount

	logger.Log.Info(
		"fetching discounts by products",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Any("product_ids", productIDs),
	)

	err :=
		database.DB.
			Table("discount_products dp").
			Select(`
				dp.product_id,
				p.id,
				p.name,
				p.type,
				p.start_date,
				p.end_date,
				dr.id,
				dr.promotion_id,
				dr.apply_to,
				dr.discount_type,
				dr.discount_value,
				p.priority
			`).
			Joins(`
				INNER JOIN discount_rules dr
				ON dr.id = dp.discount_rule_id
			`).
			Joins(`
				INNER JOIN promotions p
				ON p.id = dr.promotion_id
			`).
			Where("dp.product_id IN ?", productIDs).
			Where("p.tenant_code = ?", tenantCode).
			Where("p.country_code = ?", countryCode).
			Where("p.type = ?", "discount").
			Where("p.status = ?", "active").
			Where("p.start_date <= ?", time.Now()).
			Where("p.end_date >= ?", time.Now()).
			Find(&results).
			Error

	if err != nil {
		logger.Log.Error(
			"error fetching discounts by products",
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Any("product_ids", productIDs),
			zap.Error(err),
		)
	}

	logger.Log.Info(
		"discounts by products fetched",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Any("product_ids", productIDs),
	)

	return results, err
}

type ProductCoupons struct {
	ProductID uint64 `gorm:"column:product_id"`

	CouponRuleID uint64 `gorm:"column:coupon_rule_id"`

	PromotionID uint64 `gorm:"column:promotion_id"`

	CouponCode string `gorm:"column:coupon_code"`

	CouponType string `gorm:"column:coupon_type"`

	ApplyTo string `gorm:"column:apply_to"`

	Percentage float64 `gorm:"column:percentage"`

	Amount float64 `gorm:"column:amount"`

	ProductType string `gorm:"column:product_type"`

	TotalUsed int `gorm:"column:total_used"`

	Priority int `gorm:"column:priority"`

	StartDate time.Time `gorm:"column:start_date"`

	EndDate time.Time `gorm:"column:end_date"`
}

func GetActiveCouponsByProducts(
	tenantCode string,
	countryCode string,
	productIDs []uint64,
) ([]ProductCoupons, error) {

	var results []ProductCoupons

	logger.Log.Info(
		"fetching coupons by products",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Any("product_ids", productIDs),
	)

	err :=
		database.DB.
			Table("coupon_products cp").
			Select(`
				cp.product_id,
				cr.id as coupon_rule_id,
				cr.promotion_id,
				cr.coupon_code,
				cr.coupon_type,
				cr.apply_to,
				cr.percentage,
				cr.amount,
				cr.product_type,
				cr.total_used,
				p.start_date,
				p.end_date,
				p.priority
			`).
			Joins(`
				INNER JOIN coupon_rules cr
				ON cr.id = cp.coupon_rule_id
			`).
			Joins(`
				INNER JOIN promotions p
				ON p.id = cr.promotion_id
			`).
			Where("cp.product_id IN ?", productIDs).
			Where("p.tenant_code = ?", tenantCode).
			Where("p.country_code = ?", countryCode).
			Where("p.type = ?", "coupon").
			Where("p.status = ?", "active").
			Where("p.start_date <= ?", time.Now()).
			Where("p.end_date >= ?", time.Now()).
			Find(&results).
			Error

	if err != nil {
		logger.Log.Error(
			"error fetching coupons by products",
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Any("product_ids", productIDs),
			zap.Error(err),
		)
	}

	logger.Log.Info(
		"coupons by products fetched",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Any("product_ids", productIDs),
	)

	return results, err
}

type ProductCashback struct {
	ProductID uint64
	model.CashbackRule
}

func GetCashbacksByProducts(
	tenantCode string,
	countryCode string,
	productIDs []uint64,
) ([]ProductCashback, error) {

	var results []ProductCashback

	logger.Log.Info(
		"fetching cashbacks by products",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Any("product_ids", productIDs),
	)

	err :=
		database.DB.
			Table("cashback_products cp").
			Select(`
				cp.product_id,
				cr.*
			`).
			Joins(`
				INNER JOIN cashback_rules cr
				ON cr.id = cp.cashback_rule_id
			`).
			Joins(`
				INNER JOIN promotions p
				ON p.id = cr.promotion_id
			`).
			Where("cp.product_id IN ?", productIDs).
			Where("p.tenant_code = ?", tenantCode).
			Where("p.country_code = ?", countryCode).
			Where("p.type = ?", "cashback").
			Where("p.status = ?", "active").
			Where("p.start_date <= ?", time.Now()).
			Where("p.end_date >= ?", time.Now()).
			Find(&results).
			Error

	if err != nil {
		logger.Log.Error(
			"error fetching cashbacks by products",
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Any("product_ids", productIDs),
			zap.Error(err),
		)
	}

	logger.Log.Info(
		"cashbacks by products fetched",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Any("product_ids", productIDs),
	)

	return results, err
}

type ProductFOC struct {
	ProductID uint64
	RuleID    uint64
}

func GetFOCByProducts(
	tenantCode string,
	countryCode string,
	productIDs []uint64,
) ([]ProductFOC, error) {

	var results []ProductFOC

	logger.Log.Info(
		"fetching foc by products",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Any("product_ids", productIDs),
	)

	err :=
		database.DB.
			Table("foc_products fp").
			Select(`
				fp.product_id,
				fp.foc_rule_id
			`).
			Joins(`
				INNER JOIN foc_rules fr
				ON fr.id = fp.foc_rule_id
			`).
			Joins(`
				INNER JOIN promotions p
				ON p.id = fr.promotion_id
			`).
			Where("fp.product_id IN ?", productIDs).
			Where("p.tenant_code = ?", tenantCode).
			Where("p.country_code = ?", countryCode).
			Where("p.type = ?", "foc").
			Where("p.status = ?", "active").
			Where("p.start_date <= ?", time.Now()).
			Where("p.end_date >= ?", time.Now()).
			Find(&results).
			Error

	if err != nil {
		logger.Log.Error(
			"error fetching foc by products",
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Any("product_ids", productIDs),
			zap.Error(err),
		)
	}

	logger.Log.Info(
		"foc by products fetched",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Any("product_ids", productIDs),
	)

	return results, err
}

type ProductBuyXGetY struct {
	ProductID uint64
	model.BuyXGetYRule
}

func GetBuyXGetYByProducts(
	tenantCode string,
	countryCode string,
	productIDs []uint64,
) ([]ProductBuyXGetY, error) {

	var results []ProductBuyXGetY

	logger.Log.Info(
		"fetching buy x get y by products",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Any("product_ids", productIDs),
	)

	err :=
		database.DB.
			Table("buy_x_get_y_products bp").
			Select(`
				bp.product_id,
				r.*
			`).
			Joins(`
				INNER JOIN buy_x_get_y_rules r
				ON r.id = bp.rule_id
			`).
			Joins(`
				INNER JOIN promotions p
				ON p.id = r.promotion_id
			`).
			Where("bp.product_id IN ?", productIDs).
			Where("p.tenant_code = ?", tenantCode).
			Where("p.country_code = ?", countryCode).
			Where("p.type = ?", "buy_x_get_y").
			Where("p.status = ?", "active").
			Where("p.start_date <= ?", time.Now()).
			Where("p.end_date >= ?", time.Now()).
			Find(&results).
			Error

	if err != nil {
		logger.Log.Error(
			"error fetching buy x get y by products",
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Any("product_ids", productIDs),
			zap.Error(err),
		)
	}

	logger.Log.Info(
		"buy x get y by products fetched",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Any("product_ids", productIDs),
	)

	return results, err
}
