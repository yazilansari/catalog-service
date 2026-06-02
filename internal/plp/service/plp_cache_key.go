package service

import (
	"catalog-service/internal/plp/dto"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type cacheKeyPayload struct {
	TenantCode  string `json:"tenant_code"`
	CountryCode string `json:"country_code"`

	Category    string `json:"category,omitempty"`
	SubCategory string `json:"subcategory,omitempty"`
	Brand       string `json:"brand,omitempty"`

	Sort   string `json:"sort,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Cursor string `json:"cursor,omitempty"`

	MinPrice *float64 `json:"min_price,omitempty"`
	MaxPrice *float64 `json:"max_price,omitempty"`
}

func buildCacheKey(
	query dto.ProductQuery,
	tenantCode string,
	countryCode string,
) string {

	payload := cacheKeyPayload{
		TenantCode:  tenantCode,
		CountryCode: countryCode,

		Category:    query.Category,
		SubCategory: query.SubCategory,
		Brand:       query.Brand,

		Sort:   query.Sort,
		Limit:  query.Limit,
		Cursor: query.Cursor,

		// Future-proof
		// MinPrice: query.MinPrice,
		// MaxPrice: query.MaxPrice,
	}

	data, err := json.Marshal(payload)
	if err != nil {

		// Fallback key
		return fmt.Sprintf(
			"plp:%s:%s",
			tenantCode,
			countryCode,
		)
	}

	hash := sha256.Sum256(data)

	shortHash := hex.EncodeToString(hash[:8])

	prefix := "all"

	if query.Category != "" {
		prefix = query.Category
	}

	if query.SubCategory != "" {
		prefix = query.SubCategory
	}

	if query.Brand != "" {
		prefix = query.Brand
	}

	return fmt.Sprintf(
		"plp:%s:%s",
		prefix,
		shortHash,
	)
}
