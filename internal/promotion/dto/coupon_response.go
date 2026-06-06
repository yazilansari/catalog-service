package dto

type CouponResponse struct {
	Code      string  `json:"code"`
	Type      string  `json:"type"`
	Value     float64 `json:"value"`
	StartDate string  `json:"start_date"`
	EndDate   string  `json:"end_date"`
}
