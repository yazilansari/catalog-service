package model

import "time"

type Promotion struct {
	ID          uint64    `gorm:"column:id"`
	Name        string    `gorm:"column:name"`
	Type        string    `gorm:"column:type"`
	Description string    `gorm:"column:description"`
	StartDate   time.Time `gorm:"column:start_date"`
	EndDate     time.Time `gorm:"column:end_date"`
}

func (Promotion) TableName() string {
	return "promotions"
}
