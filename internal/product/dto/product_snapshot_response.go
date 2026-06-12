package dto

type ProductSnapshotResponse struct {
	ID uint64 `json:"id"`

	Name string `json:"name"`

	Slug string `json:"slug"`

	SKU string `json:"sku"`

	Image string `json:"image"`

	Price float64 `json:"price"`

	Stock int `json:"stock"`
}
