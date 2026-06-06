package model

type CashbackProduct struct {
	ID                 uint64 `gorm:"column:id"`
	CashbackRuleID     uint64 `gorm:"column:cashback_rule_id"`
	ProductID          uint64 `gorm:"column:product_id"`
	CashbackCustomerID uint64 `gorm:"column:cashback_customer_id"`
}

func (CashbackProduct) TableName() string {
	return "cashback_products"
}
