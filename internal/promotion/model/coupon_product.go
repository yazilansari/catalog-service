package model

type CouponProduct struct {
	ID           uint64 `gorm:"column:id"`
	CouponRuleID uint64 `gorm:"column:coupon_rule_id"`
	ProductID    uint64 `gorm:"column:product_id"`
}

func (CouponProduct) TableName() string {
	return "coupon_products"
}
