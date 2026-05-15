package dto

type CategoryResponse struct {
	ID                   uint64              `json:"id"`
	Name                 string              `json:"name"`
	Slug                 string              `json:"slug"`
	Image                *string             `json:"image,omitempty"`
	IconImage            *string             `json:"icon_image,omitempty"`
	MenuImage            *string             `json:"menu_image,omitempty"`
	MenuImage2           *string             `json:"menu_image2,omitempty"`
	MobileImage          *string             `json:"mobile_image,omitempty"`
	Video                *string             `json:"video"`
	ProductSubCategories []*CategoryResponse `json:"productSubCategories,omitempty"`
}
