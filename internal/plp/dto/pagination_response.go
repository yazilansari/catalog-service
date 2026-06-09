package dto

type PaginationResponse struct {
	NextCursor string `json:"nextCursor"`

	Limit int `json:"limit"`

	HasMore bool `json:"hasMore"`
}
