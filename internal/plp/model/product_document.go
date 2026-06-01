package model

type ProductDocument struct {
	ID uint64 `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`

	Description string `json:"description"`

	Image string `json:"image"`

	Category string `json:"category"`

	SubCategory string `json:"subcategory"`

	Brand string `json:"brand"`

	Price float64 `json:"price"`

	DiscountPrice float64 `json:"discount_price"`

	Status string `json:"status"`

	CreatedAt string `json:"created_at"`
}
