package model

type CashbackRule struct {
	ID                 uint64  `gorm:"column:id"`
	PromotionID        uint64  `gorm:"column:promotion_id"`
	CustomerType       string  `gorm:"column:customer_type"`
	ProductType        string  `gorm:"column:product_type"`
	CashbackPercentage float64 `gorm:"column:cashback_percentage"`
	CashbackAmount     float64 `gorm:"column:cashback_amount"`
	ExpiryDays         int     `gorm:"column:expiry_days"`
}

func (CashbackRule) TableName() string {
	return "cashback_rules"
}
