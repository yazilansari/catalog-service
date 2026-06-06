package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/promotion/dto"
	"catalog-service/internal/promotion/mapper"
	"catalog-service/internal/promotion/repository"
	"time"

	"go.uber.org/zap"
)

func GetProductPromotions(
	tenantCode string,
	countryCode string,
	productID uint64,
	productPrice float64,
) (*dto.PromotionResponse, error) {

	logger.Log.Info(
		"fetching product promotions",

		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Uint64("product_id", productID),
	)

	response :=
		&dto.PromotionResponse{}

	// =================================
	// DISCOUNT
	// =================================

	discountStart := time.Now()

	promotion,
		discountRule,
		err :=
		repository.GetActiveDiscountByProduct(
			tenantCode,
			countryCode,
			productID,
		)

	discountDuration := time.Since(discountStart)

	// =========================
	// SLOW DISCOUNT QUERY DETECTION
	// =========================

	if discountDuration > time.Second {

		logger.Log.Warn(
			"slow discount query detected",

			zap.Uint64("product_id", productID),

			zap.Duration(
				"duration",
				discountDuration,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"error fetching discount",
			zap.Error(err),
			zap.Duration("duration", discountDuration),
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Uint64("product_id", productID),
		)
		return nil, err
	}

	if promotion != nil &&
		discountRule != nil {

		logger.Log.Info(
			"started mapping for discount",
		)

		response.Discount =
			mapper.BuildDiscountResponse(
				productPrice,
				promotion,
				discountRule,
			)

		logger.Log.Info(
			"done mapping for discount",
			zap.Duration("duration", discountDuration),
		)
	}

	// =================================
	// COUPONS
	// =================================

	couponStart := time.Now()

	coupons,
		err :=
		repository.GetActiveCouponsByProduct(
			tenantCode,
			countryCode,
			productID,
		)

	couponDuration := time.Since(couponStart)

	// =========================
	// SLOW COUPON QUERY DETECTION
	// =========================

	if couponDuration > time.Second {

		logger.Log.Warn(
			"slow coupon query detected",

			zap.Uint64("product_id", productID),

			zap.Duration(
				"duration",
				couponDuration,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"error fetching coupons",
			zap.Error(err),
			zap.Duration("duration", couponDuration),
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Uint64("product_id", productID),
		)
		return nil, err
	}

	logger.Log.Info(
		"started mapping for coupons",
	)

	for _, coupon := range coupons {

		response.Coupons =
			append(
				response.Coupons,

				dto.CouponResponse{
					Code: coupon.CouponCode,
					Type: coupon.CouponType,
					Value: func() float64 {

						if coupon.Percentage > 0 {
							return coupon.Percentage
						}

						return coupon.Amount

					}(),
				},
			)
	}

	logger.Log.Info(
		"done mapping for coupons",
		zap.Duration("duration", couponDuration),
	)

	// =================================
	// CASHBACK
	// =================================

	cashbackStart := time.Now()

	cashback,
		err :=
		repository.GetCashbackByProduct(
			tenantCode,
			countryCode,
			productID,
		)

	cashbackDuration := time.Since(cashbackStart)

	// =========================
	// SLOW CASHBACK QUERY DETECTION
	// =========================

	if cashbackDuration > time.Second {

		logger.Log.Warn(
			"slow cashback query detected",

			zap.Uint64("product_id", productID),

			zap.Duration(
				"duration",
				cashbackDuration,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"error fetching cashback",
			zap.Error(err),
			zap.Duration("duration", cashbackDuration),
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Uint64("product_id", productID),
		)
		return nil, err
	}

	if cashback != nil {

		logger.Log.Info(
			"started mapping for cashback",
		)

		response.Cashback =
			mapper.BuildCashbackResponse(
				cashback,
			)
	}

	logger.Log.Info(
		"done mapping for cashback",
		zap.Duration("duration", cashbackDuration),
	)

	// =================================
	// FOC
	// =================================

	focStart := time.Now()

	foc,
		err :=
		repository.GetFOCByProduct(
			tenantCode,
			countryCode,
			productID,
		)

	focDuration := time.Since(focStart)

	// =========================
	// SLOW FOC QUERY DETECTION
	// =========================

	if focDuration > time.Second {

		logger.Log.Warn(
			"slow FOC query detected",

			zap.Uint64("product_id", productID),

			zap.Duration(
				"duration",
				focDuration,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"error fetching FOC",
			zap.Error(err),
			zap.Duration("duration", focDuration),
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Uint64("product_id", productID),
		)
		return nil, err
	}

	if foc != nil {

		logger.Log.Info(
			"started mapping for FOC",
		)

		response.FOC =
			mapper.BuildFOCResponse(
				foc.ProductID,
			)
	}

	logger.Log.Info(
		"done mapping for FOC",
		zap.Duration("duration", focDuration),
	)

	// =================================
	// BUY X GET Y
	// =================================

	buyxgetyStart := time.Now()

	buyXGetY,
		err :=
		repository.GetBuyXGetYByProduct(
			tenantCode,
			countryCode,
			productID,
		)

	buyXGetYDuration := time.Since(buyxgetyStart)

	// =========================
	// SLOW BUY X GET Y QUERY DETECTION
	// =========================

	if buyXGetYDuration > time.Second {

		logger.Log.Warn(
			"slow buy x get y query detected",

			zap.Uint64("product_id", productID),

			zap.Duration(
				"duration",
				buyXGetYDuration,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"error fetching buy x get y",
			zap.Error(err),
			zap.Duration("duration", buyXGetYDuration),
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Uint64("product_id", productID),
		)
		return nil, err
	}

	if buyXGetY != nil {

		logger.Log.Info(
			"started mapping for buy x get y",
		)

		response.BuyXGetY =
			mapper.BuildBuyXGetYResponse(
				buyXGetY,
			)
	}

	logger.Log.Info(
		"done mapping for buy x get y",
		zap.Duration("duration", buyXGetYDuration),
	)

	return response, nil
}

func GetProductsPromotions(
	tenantCode string,
	countryCode string,
	productIDs []uint64,
	productPriceMap map[uint64]float64,
) (map[uint64]*dto.PromotionResponse, error) {

	logger.Log.Info(
		"fetching products promotions",

		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
		zap.Any("product_ids", productIDs),
	)

	result :=
		make(
			map[uint64]*dto.PromotionResponse,
		)

	for _, productID := range productIDs {

		result[productID] =
			&dto.PromotionResponse{}
	}

	// =========================
	// DISCOUNTS
	// =========================

	discountStart := time.Now()

	discounts, err :=
		repository.GetActiveDiscountsByProducts(
			tenantCode,
			countryCode,
			productIDs,
		)

	discountDuration := time.Since(discountStart)

	// =========================
	// SLOW DISCOUNT QUERY DETECTION
	// =========================

	if discountDuration > time.Second {

		logger.Log.Warn(
			"slow discount query detected",

			zap.Any("product_ids", productIDs),

			zap.Duration(
				"duration",
				discountDuration,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"error fetching discounts",
			zap.Error(err),
			zap.Duration("duration", discountDuration),
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Any("product_ids", productIDs),
		)
		return nil, err
	}

	logger.Log.Info(
		"started fetching discounts",
	)

	for _, item := range discounts {

		result[item.ProductID].
			Discount =
			mapper.BuildDiscountResponse(
				productPriceMap[item.ProductID],
				&item.Promotion,
				&item.DiscountRule,
			)
	}

	logger.Log.Info(
		"done fetching discounts",
		zap.Duration("duration", discountDuration),
	)

	// =========================
	// COUPONS
	// =========================

	couponStart := time.Now()

	coupons, err :=
		repository.GetActiveCouponsByProducts(
			tenantCode,
			countryCode,
			productIDs,
		)

	couponDuration := time.Since(couponStart)

	// =========================
	// SLOW COUPON QUERY DETECTION
	// =========================

	if couponDuration > time.Second {

		logger.Log.Warn(
			"slow coupon query detected",

			zap.Any("product_ids", productIDs),

			zap.Duration(
				"duration",
				couponDuration,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"error fetching coupons",
			zap.Error(err),
			zap.Duration("duration", couponDuration),
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Any("product_ids", productIDs),
		)
		return nil, err
	}

	logger.Log.Info(
		"started fetching coupons",
	)

	for _, coupon := range coupons {

		productID :=
			coupon.ProductID

		result[productID].
			Coupons =
			append(
				result[productID].Coupons,

				dto.CouponResponse{
					Code: coupon.CouponCode,
					Type: coupon.CouponType,
					Value: func() float64 {

						if coupon.Percentage > 0 {
							return coupon.Percentage
						}

						return coupon.Amount

					}(),
				},
			)
	}

	logger.Log.Info(
		"done fetching coupons",
		zap.Duration("duration", couponDuration),
	)

	// =========================
	// CASHBACK
	// =========================

	cashbackStart := time.Now()

	cashbacks, err :=
		repository.GetCashbacksByProducts(
			tenantCode,
			countryCode,
			productIDs,
		)

	cashbackDuration := time.Since(cashbackStart)

	// =========================
	// SLOW CASHBACK QUERY DETECTION
	// =========================

	if cashbackDuration > time.Second {

		logger.Log.Warn(
			"slow cashback query detected",

			zap.Any("product_ids", productIDs),

			zap.Duration(
				"duration",
				cashbackDuration,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"error fetching cashbacks",
			zap.Error(err),
			zap.Duration("duration", cashbackDuration),
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Any("product_ids", productIDs),
		)
		return nil, err
	}

	logger.Log.Info(
		"started fetching cashbacks",
		zap.Duration("duration", cashbackDuration),
	)

	for _, cashback := range cashbacks {

		result[cashback.ProductID].
			Cashback =
			mapper.BuildCashbackResponse(
				&cashback.CashbackRule,
			)
	}

	logger.Log.Info(
		"done fetching cashbacks",
		zap.Duration("duration", cashbackDuration),
	)

	// =========================
	// FOC
	// =========================

	focStart := time.Now()

	focs, err :=
		repository.GetFOCByProducts(
			tenantCode,
			countryCode,
			productIDs,
		)

	focDuration := time.Since(focStart)

	// =========================
	// SLOW FOC QUERY DETECTION
	// =========================

	if focDuration > time.Second {

		logger.Log.Warn(
			"slow foc query detected",

			zap.Any("product_ids", productIDs),

			zap.Duration(
				"duration",
				focDuration,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"error fetching focs",
			zap.Error(err),
			zap.Duration("duration", focDuration),
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Any("product_ids", productIDs),
		)
		return nil, err
	}

	logger.Log.Info(
		"started fetching focs",
		zap.Duration("duration", focDuration),
	)

	for _, foc := range focs {

		result[foc.ProductID].
			FOC =
			mapper.BuildFOCResponse(
				foc.ProductID,
			)
	}

	logger.Log.Info(
		"done fetching focs",
		zap.Duration("duration", focDuration),
	)

	// =========================
	// BUY X GET Y
	// =========================

	buyXGetYStart := time.Now()

	buyXGetYRules, err :=
		repository.GetBuyXGetYByProducts(
			tenantCode,
			countryCode,
			productIDs,
		)

	buyXGetYDuration := time.Since(buyXGetYStart)

	// =========================
	// SLOW BUY X GET Y QUERY DETECTION
	// =========================

	if buyXGetYDuration > time.Second {

		logger.Log.Warn(
			"slow buy x get y query detected",

			zap.Any("product_ids", productIDs),

			zap.Duration(
				"duration",
				buyXGetYDuration,
			),

			zap.String(
				"country_code",
				countryCode,
			),

			zap.String(
				"tenant_code",
				tenantCode,
			),
		)
	}

	if err != nil {

		logger.Log.Error(
			"error fetching buy x get y rules",
			zap.Error(err),
			zap.Duration("duration", buyXGetYDuration),
			zap.String("tenant_code", tenantCode),
			zap.String("country_code", countryCode),
			zap.Any("product_ids", productIDs),
		)
		return nil, err
	}

	logger.Log.Info(
		"started fetching buy x get y rules",
		zap.Duration("duration", buyXGetYDuration),
	)

	for _, rule := range buyXGetYRules {

		result[rule.ProductID].
			BuyXGetY =
			mapper.BuildBuyXGetYResponse(
				&rule.BuyXGetYRule,
			)
	}

	logger.Log.Info(
		"done fetching buy x get y rules",
		zap.Duration("duration", buyXGetYDuration),
	)

	return result, nil
}
