package model

import "time"

// =========================
// CATEGORY MODEL
// =========================

type Category struct {
	ID uint64 `gorm:"column:id"`

	TenantCode string `gorm:"column:tenant_code"`

	CountryCode string `gorm:"column:country_code"`

	Name string `gorm:"column:name"`

	Slug string `gorm:"column:slug"`

	Description string `gorm:"column:description"`

	Image string `gorm:"column:image"`

	IconImage string `gorm:"column:icon_image"`

	MenuImage string `gorm:"column:menu_image"`

	MobileImage string `gorm:"column:mobile_image"`

	ParentID *uint64 `gorm:"column:parent_id"`

	SortOrder int `gorm:"column:sort_order"`

	Status string `gorm:"column:status"`

	CreatedAt time.Time `gorm:"column:created_at"`

	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Category) TableName() string {
	return "categories"
}

// =========================
// PRODUCT MODEL
// =========================

type Product struct {
	ID uint64 `gorm:"column:id"`

	TenantCode string `gorm:"column:tenant_code"`

	CountryCode string `gorm:"column:country_code"`

	CategoryID uint64 `gorm:"column:category_id"`

	SubCategoryID uint64 `gorm:"column:subcategory_id"`

	Name string `gorm:"column:name"`

	Slug string `gorm:"column:slug"`

	SKU string `gorm:"column:sku"`

	ShortDescription string `gorm:"column:short_description"`

	Description string `gorm:"column:description"`

	Image string `gorm:"column:image"`

	Price float64 `gorm:"column:price"`

	SalePrice float64 `gorm:"column:sale_price"`

	Stock int `gorm:"column:stock"`

	IsFeatured bool `gorm:"column:is_featured"`

	IsBestSeller bool `gorm:"column:is_best_seller"`

	IsNewArrival bool `gorm:"column:is_new_arrival"`

	Status string `gorm:"column:status"`

	CreatedAt time.Time `gorm:"column:created_at"`

	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Product) TableName() string {
	return "products"
}
