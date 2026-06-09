package model

import promotionDto "catalog-service/internal/promotion/dto"

type Product struct {
	ID uint64 `gorm:"column:id" json:"id"`

	Name string `gorm:"column:name" json:"name"`
	Slug string `gorm:"column:slug" json:"slug"`

	ShortDescription string `gorm:"column:description" json:"shortDescription"`

	Description string `gorm:"column:content" json:"description"`

	SKU string `gorm:"column:sku" json:"sku"`

	Image string `gorm:"column:image" json:"image"`

	Price float64 `gorm:"column:price" json:"price"`

	// DiscountPrice float64 `gorm:"column:discount_price"`

	// Brand string `gorm:"column:brand"`

	Stock int `gorm:"column:stock" json:"stock"`

	// Status string `gorm:"column:status"`

	// TenantCode  string `gorm:"column:tenant_code"`
	// CountryCode string `gorm:"column:country_code"`

	// CreatedAt time.Time `gorm:"column:created_at"`
	// UpdatedAt time.Time `gorm:"column:updated_at"`

	IsFeatured bool `gorm:"column:is_featured" json:"isFeatured"`

	IsBestSeller bool `gorm:"column:is_best_seller" json:"isBestSeller"`

	IsNewArrival bool `gorm:"column:is_new_arrival" json:"isNewArrival"`

	Promotion *promotionDto.PromotionResponse `json:"promotion,omitempty" gorm:"-"`
}

func (Product) TableName() string {
	return "products"
}
