package model

import promotionDto "catalog-service/internal/promotion/dto"

// =========================
// CATEGORY MODEL
// =========================

type Category struct {
	ID uint64 `gorm:"column:id" json:"id"`

	// TenantCode string `gorm:"column:tenant_code"`

	// CountryCode string `gorm:"column:country_code"`

	Name string `gorm:"column:name" json:"name"`

	Slug string `gorm:"column:slug" json:"slug"`

	Description string `gorm:"column:description" json:"description"`

	Image string `gorm:"column:image" json:"image"`

	IconImage string `gorm:"column:icon_image" json:"iconImage"`

	MenuImage string `gorm:"column:menu_image" json:"menuImage"`

	MobileImage string `gorm:"column:mobile_image" json:"mobileImage"`

	ParentID *uint64 `gorm:"column:parent_id" json:"-"`

	// SortOrder int `gorm:"column:sort_order"`

	// Status string `gorm:"column:status"`

	// CreatedAt time.Time `gorm:"column:created_at"`

	// UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Category) TableName() string {
	return "categories"
}

// =========================
// PRODUCT MODEL
// =========================

type Product struct {
	ID uint64 `gorm:"column:id" json:"id"`

	// TenantCode string `gorm:"column:tenant_code"`

	// CountryCode string `gorm:"column:country_code"`

	// CategoryID uint64 `gorm:"column:category_id"`

	// SubCategoryID uint64 `gorm:"column:subcategory_id"`

	Name string `gorm:"column:name" json:"name"`

	Slug string `gorm:"column:slug" json:"slug"`

	SKU string `gorm:"column:sku" json:"sku"`

	ShortDescription string `gorm:"column:description" json:"shortDescription"`

	// Description string `gorm:"column:content"`

	Image string `gorm:"column:image" json:"image"`

	Price float64 `gorm:"column:price" json:"price"`

	// SalePrice float64 `gorm:"column:sale_price"`

	Stock int `gorm:"column:stock" json:"stock"`

	IsFeatured bool `gorm:"column:is_featured" json:"isFeatured"`

	IsBestSeller bool `gorm:"column:is_best_seller" json:"isBestSeller"`

	IsNewArrival bool `gorm:"column:is_new_arrival" json:"isNewArrival"`

	// Status string `gorm:"column:status"`

	// CreatedAt time.Time `gorm:"column:created_at"`

	// UpdatedAt time.Time `gorm:"column:updated_at"`

	Promotion *promotionDto.PromotionResponse `json:"promotion,omitempty" gorm:"-"`
}

func (Product) TableName() string {
	return "products"
}

// type ProductImage struct {
// 	ID uint64 `gorm:"column:id"`

// 	ProductID uint64 `gorm:"column:product_id"`

// 	Image string `gorm:"column:image"`

// 	SortOrder int `gorm:"column:sort_order"`
// }

// func (ProductImage) TableName() string {
// 	return "product_images"
// }
