package model

import "time"

type Product struct {
	ID uint64 `gorm:"column:id"`

	Name string `gorm:"column:name"`
	Slug string `gorm:"column:slug"`

	Description string `gorm:"column:description"`

	SKU string `gorm:"column:sku"`

	Price float64 `gorm:"column:price"`

	DiscountPrice float64 `gorm:"column:discount_price"`

	Brand string `gorm:"column:brand"`

	Stock int `gorm:"column:stock"`

	Status string `gorm:"column:status"`

	TenantCode  string `gorm:"column:tenant_code"`
	CountryCode string `gorm:"column:country_code"`

	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Product) TableName() string {
	return "products"
}
