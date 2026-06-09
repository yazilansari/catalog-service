package dto

import "catalog-service/internal/product/model"

type ProductResponse struct {
	PageType string `json:"pageType"`

	Category model.Category `json:"category"`

	SubCategory model.Category `json:"subCategory"`

	Product model.Product `json:"product"`

	Images []model.ProductImage `json:"images"`

	Variants []ProductVariantResponse `json:"variants"`

	FragranceNotes []FragranceNoteResponse `json:"fragranceNotes"`

	RelatedProducts []RelatedProductResponse `json:"relatedProducts"`

	SEO model.SEO `json:"seo"`
}
