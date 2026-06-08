package model

import promotionDto "catalog-service/internal/promotion/dto"

type Product struct {
	ID uint64 `gorm:"column:id"`

	Name string `gorm:"column:name"`
	Slug string `gorm:"column:slug"`

	ShortDescription string `gorm:"column:description"`

	Description string `gorm:"column:content"`

	SKU string `gorm:"column:sku"`

	Image string `gorm:"column:image"`

	Price float64 `gorm:"column:price"`

	// DiscountPrice float64 `gorm:"column:discount_price"`

	// Brand string `gorm:"column:brand"`

	Stock int `gorm:"column:stock"`

	// Status string `gorm:"column:status"`

	// TenantCode  string `gorm:"column:tenant_code"`
	// CountryCode string `gorm:"column:country_code"`

	// CreatedAt time.Time `gorm:"column:created_at"`
	// UpdatedAt time.Time `gorm:"column:updated_at"`

	IsFeatured bool `gorm:"column:is_featured"`

	IsBestSeller bool `gorm:"column:is_best_seller"`

	IsNewArrival bool `gorm:"column:is_new_arrival"`

	Promotion *promotionDto.PromotionResponse `json:"promotion,omitempty" gorm:"-"`
}

func (Product) TableName() string {
	return "products"
}
