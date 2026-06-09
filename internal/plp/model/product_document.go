package model

import promotionDto "catalog-service/internal/promotion/dto"

type ProductDocument struct {
	ID uint64 `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`

	Description string `json:"shortDescription"`

	Image string `json:"image"`

	Category string `json:"category"`

	SubCategory string `json:"subCategory"`

	Brand string `json:"brand"`

	Price float64 `json:"price"`

	// DiscountPrice float64 `json:"discountPrice"`

	// Status string `json:"status"`

	// CreatedAt string `json:"created_at"`

	Promotion *promotionDto.PromotionResponse `json:"promotion,omitempty"`
}
