package model

import "time"

type Category struct {
	ID          uint64 `gorm:"primaryKey"`
	TenantCode  string
	CountryCode string
	ParentID    *uint64
	Name        string
	Slug        string
	Image       *string
	IconImage   *string
	MenuImage   *string
	MenuImage2  *string
	MobileImage *string
	Video       *string
	Status      string
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Category) TableName() string {
	return "categories"
}
