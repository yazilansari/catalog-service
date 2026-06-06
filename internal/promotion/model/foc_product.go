package model

type FOCProduct struct {
	ID        uint64 `gorm:"column:id"`
	FOCRuleID uint64 `gorm:"column:foc_rule_id"`
	ProductID uint64 `gorm:"column:product_id"`
}

func (FOCProduct) TableName() string {
	return "foc_products"
}
