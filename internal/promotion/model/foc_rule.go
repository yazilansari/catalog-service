package model

type FOCRule struct {
	ID           uint64  `gorm:"column:id"`
	PromotionID  uint64  `gorm:"column:promotion_id"`
	MinThreshold float64 `gorm:"column:min_threshold"`
	MaxThreshold float64 `gorm:"column:max_threshold"`
}

func (FOCRule) TableName() string {
	return "foc_rules"
}
