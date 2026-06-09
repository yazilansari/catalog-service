package dto

type DiscountResponse struct {
	Type           string  `json:"type"`
	Value          float64 `json:"value"`
	DiscountAmount float64 `json:"discountAmount"`
	FinalPrice     float64 `json:"finalPrice"`
	StartDate      string  `json:"startDate"`
	EndDate        string  `json:"endDate"`
	Priority       int     `json:"priority"`
}
