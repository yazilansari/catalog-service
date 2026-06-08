package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/logger"
	"catalog-service/internal/promotion/model"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetActiveDiscountByProduct(
	tenantCode string,
	countryCode string,
	productID uint64,
) (*model.Promotion, *model.DiscountRule, error) {

	var promotion model.Promotion
	var rule model.DiscountRule

	logger.Log.Info(
		"fetching discount by product",

		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Uint64("product_id", productID),
	)

	err :=
		database.DB.
			Table("promotions p").
			Select("p.*, dr.*").
			Joins(`
				INNER JOIN discount_rules dr
				ON dr.promotion_id = p.id
			`).
			Joins(`
				INNER JOIN discount_products dp
				ON dp.discount_rule_id = dr.id
			`).
			Where("dp.product_id = ?", productID).
			Where("p.tenant_code = ?", tenantCode).
			Where("p.country_code = ?", countryCode).
			Where("p.type = ?", "discount").
			Where("p.status = ?", "active").
			Where("p.start_date <= ?", time.Now()).
			Where("p.end_date >= ?", time.Now()).
			First(&promotion).
			Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Log.Info(
			"no active discount found",
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Uint64("product_id", productID),
		)

		return nil, nil, nil
	}

	if err != nil {

		logger.Log.Error(
			"error fetching discount by product",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Uint64("product_id", productID),
			zap.Error(err),
		)

		return nil, nil, err
	}

	err =
		database.DB.
			Where("promotion_id = ?", promotion.ID).
			First(&rule).
			Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Log.Info(
			"no discount rule found",
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Uint64("product_id", productID),
		)

		return &promotion, nil, nil
	}

	if err != nil {

		logger.Log.Error(
			"error fetching discount rule by product",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Uint64("product_id", productID),
			zap.Error(err),
		)

		return nil, nil, err
	}

	logger.Log.Info(
		"discount by product fetched",

		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Uint64("product_id", productID),
	)

	return &promotion, &rule, nil
}

type ProductCoupon struct {
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

func GetActiveCouponsByProduct(
	tenantCode string,
	countryCode string,
	productID uint64,
) ([]ProductCoupon, error) {

	var coupons []ProductCoupon

	logger.Log.Info(
		"fetching coupons by product",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Uint64("product_id", productID),
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
			Where("cp.product_id = ?", productID).
			Where("p.tenant_code = ?", tenantCode).
			Where("p.country_code = ?", countryCode).
			Where("p.type = ?", "coupon").
			Where("p.status = ?", "active").
			Where("p.start_date <= ?", time.Now()).
			Where("p.end_date >= ?", time.Now()).
			Find(&coupons).
			Error

	if err != nil {

		logger.Log.Error(
			"error fetching coupons by product",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Uint64("product_id", productID),
			zap.Error(err),
		)

		return nil, err
	}

	logger.Log.Info(
		"coupons by product fetched",

		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Uint64("product_id", productID),
	)

	return coupons, err
}

func GetCashbackByProduct(
	tenantCode string,
	countryCode string,
	productID uint64,
) (*model.CashbackRule, error) {

	var cashback model.CashbackRule

	logger.Log.Info(
		"fetching cashback by product",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Uint64("product_id", productID),
	)

	err :=
		database.DB.
			Table("cashback_rules cr").
			Select("cr.*").
			Joins(`
				INNER JOIN cashback_products cp
				ON cp.cashback_rule_id = cr.id
			`).
			Joins(`
				INNER JOIN promotions p
				ON p.id = cr.promotion_id
			`).
			Where("cp.product_id = ?", productID).
			Where("p.tenant_code = ?", tenantCode).
			Where("p.country_code = ?", countryCode).
			Where("p.type = ?", "cashback").
			Where("p.status = ?", "active").
			Where("p.start_date <= ?", time.Now()).
			Where("p.end_date >= ?", time.Now()).
			First(&cashback).
			Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Log.Info(
			"no cashback found",
			zap.Uint64("product_id", productID),
		)
		return nil, nil
	}

	if err != nil {

		logger.Log.Error(
			"error fetching cashback by product",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Uint64("product_id", productID),
			zap.Error(err),
		)

		return nil, err
	}

	logger.Log.Info(
		"cashback by product fetched",

		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Uint64("product_id", productID),
	)

	return &cashback, nil
}

func GetFOCByProduct(
	tenantCode string,
	countryCode string,
	productID uint64,
) (*model.FOCProduct, error) {

	var foc model.FOCProduct

	logger.Log.Info(
		"fetching foc by product",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Uint64("product_id", productID),
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
			Where("fp.product_id = ?", productID).
			Where("p.tenant_code = ?", tenantCode).
			Where("p.country_code = ?", countryCode).
			Where("p.type = ?", "foc").
			Where("p.status = ?", "active").
			Where("p.start_date <= ?", time.Now()).
			Where("p.end_date >= ?", time.Now()).
			First(&foc).
			Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Log.Info(
			"no foc found",
			zap.Uint64("product_id", productID),
		)
		return nil, nil
	}

	if err != nil {

		logger.Log.Error(
			"error fetching foc by product",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Uint64("product_id", productID),
			zap.Error(err),
		)

		return nil, err
	}

	logger.Log.Info(
		"foc by product fetched",

		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Uint64("product_id", productID),
	)

	return &foc, nil
}

func GetBuyXGetYByProduct(
	tenantCode string,
	countryCode string,
	productID uint64,
) (*model.BuyXGetYRule, error) {

	var rule model.BuyXGetYRule

	logger.Log.Info(
		"fetching buy x get y by product",
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Uint64("product_id", productID),
	)

	err :=
		database.DB.
			Table("buy_x_get_y_rules r").
			Select("r.*").
			Joins(`
				INNER JOIN buy_x_get_y_products p
				ON p.rule_id = r.id
			`).
			Joins(`
				INNER JOIN promotions promo
				ON promo.id = r.promotion_id
			`).
			Where("p.product_id = ?", productID).
			Where("promo.tenant_code = ?", tenantCode).
			Where("promo.country_code = ?", countryCode).
			Where("promo.type = ?", "buy_x_get_y").
			Where("promo.status = ?", "active").
			Where("promo.start_date <= ?", time.Now()).
			Where("promo.end_date >= ?", time.Now()).
			First(&rule).
			Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Log.Info(
			"no buy x get y found",
			zap.Uint64("product_id", productID),
		)
		return nil, nil
	}

	if err != nil {

		logger.Log.Error(
			"error fetching buy x get y by product",

			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Uint64("product_id", productID),
			zap.Error(err),
		)

		return nil, err
	}

	logger.Log.Info(
		"buy x get y by product fetched",

		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Uint64("product_id", productID),
	)

	return &rule, nil
}
