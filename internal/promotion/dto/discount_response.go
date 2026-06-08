package dto

type DiscountResponse struct {
	Type           string  `json:"type"`
	Value          float64 `json:"value"`
	DiscountAmount float64 `json:"discount_amount"`
	FinalPrice     float64 `json:"final_price"`
	StartDate      string  `json:"start_date"`
	EndDate        string  `json:"end_date"`
	Priority       int     `json:"priority"`
}
