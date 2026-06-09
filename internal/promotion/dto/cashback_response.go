package dto

type CashbackResponse struct {
	Type       string  `json:"type"`
	Value      float64 `json:"value"`
	ExpiryDays int     `json:"expiryDays"`
}
