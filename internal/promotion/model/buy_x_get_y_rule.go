package model

type BuyXGetYRule struct {
	ID          uint64 `gorm:"column:id"`
	PromotionID uint64 `gorm:"column:promotion_id"`
	BuyQuantity int    `gorm:"column:buy_quantity"`
	GetQuantity int    `gorm:"column:get_quantity"`
}

func (BuyXGetYRule) TableName() string {
	return "buy_x_get_y_rules"
}
