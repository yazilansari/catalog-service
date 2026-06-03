package dto

type ProductSearchQuery struct {
	Query string `form:"q"`

	Category string `form:"category"`

	SubCategory string `form:"subcategory"`

	Brand string `form:"brand"`

	Sort string `form:"sort"`

	Cursor string `form:"cursor"`

	Limit int `form:"limit"`

	MinPrice float64 `form:"min_price"`

	MaxPrice float64 `form:"max_price"`
}
