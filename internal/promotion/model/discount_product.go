package model

type DiscountProduct struct {
	ID             uint64 `gorm:"column:id"`
	DiscountRuleID uint64 `gorm:"column:discount_rule_id"`
	ProductID      uint64 `gorm:"column:product_id"`
}

func (DiscountProduct) TableName() string {
	return "discount_products"
}
