package dto

type PaginationResponse struct {
	NextCursor string `json:"next_cursor"`

	Limit int `json:"limit"`

	HasMore bool `json:"has_more"`
}
