package model

type BuyXGetYProduct struct {
	ID        uint64 `gorm:"column:id"`
	RuleID    uint64 `gorm:"column:rule_id"`
	ProductID uint64 `gorm:"column:product_id"`
	Type      string `gorm:"column:type"`
}

func (BuyXGetYProduct) TableName() string {
	return "buy_x_get_y_products"
}
