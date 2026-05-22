package dto

import "catalog-service/internal/page/model"

type SubCategoryPageResponse struct {
	PageType string `json:"page_type"`

	Category model.Category `json:"category"`

	SubCategory model.Category `json:"subcategory"`

	Products []model.Product `json:"products"`
}
