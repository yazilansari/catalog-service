package dto

type PromotionResponse struct {
	Discount *DiscountResponse `json:"discount,omitempty"`
	Coupons  []CouponResponse  `json:"coupons,omitempty"`
	Cashback *CashbackResponse `json:"cashback,omitempty"`
	FOC      *FOCResponse      `json:"foc,omitempty"`
	BuyXGetY *BuyXGetYResponse `json:"buy_x_get_y,omitempty"`
}
