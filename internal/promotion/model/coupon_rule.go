package model

type CouponRule struct {
	ID          uint64  `gorm:"column:id"`
	PromotionID uint64  `gorm:"column:promotion_id"`
	CouponCode  string  `gorm:"column:coupon_code"`
	CouponType  string  `gorm:"column:coupon_type"`
	ApplyTo     string  `gorm:"column:apply_to"`
	Percentage  float64 `gorm:"column:percentage"`
	Amount      float64 `gorm:"column:amount"`
	ProductType string  `gorm:"column:product_type"`
	TotalUsed   int     `gorm:"column:total_used"`
}
