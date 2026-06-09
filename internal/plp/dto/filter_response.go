package dto

type FilterResponse struct {
	Brands []string `json:"brands"`

	PriceRange PriceRange `json:"priceRange"`
}

type PriceRange struct {
	Min float64 `json:"min"`

	Max float64 `json:"max"`
}
