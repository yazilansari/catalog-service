package dto

import promotionDto "catalog-service/internal/promotion/dto"

type RelatedProductResponse struct {
	ID uint64 `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`

	Price float64 `json:"price"`

	// DiscountPrice float64 `json:"discount_price"`

	Promotion *promotionDto.PromotionResponse `json:"promotion,omitempty" gorm:"-"`
}
