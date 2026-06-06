package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type cacheKeyPayload struct {
	// TenantCode  string `json:"tenant_code"`
	// CountryCode string `json:"country_code"`

	// Slug string `json:"slug"`
	Page int `json:"page,omitempty"`

	Sort  string `json:"sort,omitempty"`
	Limit int    `json:"limit,omitempty"`
	// Cursor string `json:"cursor,omitempty"`

	MinPrice *float64 `json:"min_price,omitempty"`
	MaxPrice *float64 `json:"max_price,omitempty"`
}

func buildCacheKey(
	tenantCode string,
	countryCode string,
	slug string,
	page int,
	sort string,
	limit int,
	minPrice float64,
	maxPrice float64,
) string {

	payload := cacheKeyPayload{
		// TenantCode:  tenantCode,
		// CountryCode: countryCode,

		// Slug: slug,
		Page: page,

		Sort:  sort,
		Limit: limit,
		// Cursor: query.Cursor,

		// Future-proof
		MinPrice: &minPrice,
		MaxPrice: &maxPrice,
	}

	data, err := json.Marshal(payload)
	if err != nil {

		// Fallback key
		return fmt.Sprintf(
			"page:%s:%s:%s:%d:%s:%d::%.2f:%.2f",
			slug,
			tenantCode,
			countryCode,
			page,
			sort,
			limit,
			minPrice,
			maxPrice,
		)
	}

	hash := sha256.Sum256(data)

	shortHash := hex.EncodeToString(hash[:8])

	return fmt.Sprintf(
		"page:%s:%s:%s:%s",
		slug,
		tenantCode,
		countryCode,
		shortHash,
	)
}
