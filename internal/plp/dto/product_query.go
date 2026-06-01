package dto

type ProductQuery struct {
	Category string

	SubCategory string

	Brand string

	MinPrice float64

	MaxPrice float64

	Sort string

	Limit int

	Cursor string
}
