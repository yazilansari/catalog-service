package dto

type ProductVariantResponse struct {
	ID uint64 `json:"id"`

	Name string `json:"name"`

	SKU string `json:"sku"`

	Price float64 `json:"price"`

	DiscountPrice float64 `json:"discount_price"`

	Stock int `json:"stock"`
}
