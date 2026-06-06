package model

type DiscountRule struct {
	ID            uint64  `gorm:"column:id"`
	PromotionID   uint64  `gorm:"column:promotion_id"`
	ApplyTo       string  `gorm:"column:apply_to"`
	DiscountType  string  `gorm:"column:discount_type"`
	DiscountValue float64 `gorm:"column:discount_value"`
}

func (DiscountRule) TableName() string {
	return "discount_rules"
}
