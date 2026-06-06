package mapper

import (
	"catalog-service/internal/promotion/dto"
	"catalog-service/internal/promotion/model"
)

func BuildDiscountResponse(
	productPrice float64,
	promotion *model.Promotion,
	rule *model.DiscountRule,
) *dto.DiscountResponse {

	if promotion == nil ||
		rule == nil {

		return nil
	}

	var discountAmount float64
	var finalPrice float64

	switch rule.DiscountType {

	case "percent":

		discountAmount =
			productPrice *
				rule.DiscountValue /
				100

		finalPrice =
			productPrice -
				discountAmount

	case "fixed":

		discountAmount =
			rule.DiscountValue

		finalPrice =
			productPrice -
				discountAmount

		if finalPrice < 0 {
			finalPrice = 0
		}
	}

	return &dto.DiscountResponse{
		Type:           rule.DiscountType,
		Value:          rule.DiscountValue,
		DiscountAmount: discountAmount,
		FinalPrice:     finalPrice,
		StartDate: promotion.StartDate.Format(
			"2006-01-02 15:04:05",
		),
		EndDate: promotion.EndDate.Format(
			"2006-01-02 15:04:05",
		),
	}
}

func BuildCouponResponse(
	promotion *model.Promotion,
	rule *model.CouponRule,
) dto.CouponResponse {

	var value float64

	if rule.Percentage > 0 {
		value = rule.Percentage
	} else {
		value = rule.Amount
	}

	return dto.CouponResponse{
		Code:  rule.CouponCode,
		Type:  rule.CouponType,
		Value: value,
		StartDate: promotion.StartDate.Format(
			"2006-01-02 15:04:05",
		),
		EndDate: promotion.EndDate.Format(
			"2006-01-02 15:04:05",
		),
	}
}

func BuildCashbackResponse(
	rule *model.CashbackRule,
) *dto.CashbackResponse {

	if rule == nil {
		return nil
	}

	var value float64

	if rule.CashbackPercentage > 0 {
		value = rule.CashbackPercentage
	} else {
		value = rule.CashbackAmount
	}

	return &dto.CashbackResponse{
		Type:       rule.ProductType,
		Value:      value,
		ExpiryDays: rule.ExpiryDays,
	}
}

func BuildFOCResponse(
	productID uint64,
) *dto.FOCResponse {

	return &dto.FOCResponse{
		ProductID: productID,
	}
}

func BuildBuyXGetYResponse(
	rule *model.BuyXGetYRule,
) *dto.BuyXGetYResponse {

	if rule == nil {
		return nil
	}

	return &dto.BuyXGetYResponse{
		BuyQuantity: rule.BuyQuantity,
		GetQuantity: rule.GetQuantity,
	}
}
