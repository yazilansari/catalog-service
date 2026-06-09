package dto

import "catalog-service/internal/page/model"

type SubCategoryPageResponse struct {
	PageType string `json:"pageType"`

	Category model.Category `json:"category"`

	SubCategory model.Category `json:"subCategory"`

	Products []model.Product `json:"products"`
}
