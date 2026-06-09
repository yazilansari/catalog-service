package dto

type CouponResponse struct {
	Code      string  `json:"code"`
	Type      string  `json:"type"`
	Value     float64 `json:"value"`
	StartDate string  `json:"startDate"`
	EndDate   string  `json:"endDate"`
	Priority  int     `json:"priority"`
}
