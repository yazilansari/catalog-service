package dto

import "catalog-service/internal/product/model"

type ProductResponse struct {
	PageType string `json:"page_type"`

	Category model.Category `json:"category"`

	SubCategory model.Category `json:"subcategory"`

	Product model.Product `json:"product"`

	Images []model.ProductImage `json:"images"`

	Variants []ProductVariantResponse `json:"variants"`

	FragranceNotes []FragranceNoteResponse `json:"fragrance_notes"`

	RelatedProducts []RelatedProductResponse `json:"related_products"`

	SEO model.SEO `json:"seo"`
}
