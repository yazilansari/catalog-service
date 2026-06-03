package service

import (
	"catalog-service/internal/logger"
	"catalog-service/internal/product_search/dto"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

func buildCacheKey(
	query dto.ProductSearchQuery,
	tenantCode string,
	countryCode string,
) string {

	logger.Log.Info("building cache key",

		zap.String("query", query.Query),
		zap.String("tenant_code", tenantCode),
		zap.String("country_code", countryCode),
	)

	data, _ :=
		json.Marshal(query)

	hash :=
		sha256.Sum256(data)

	shortHash :=
		hex.EncodeToString(hash[:8])

	prefix := query.Query

	if prefix == "" {
		prefix = "all"
	}

	return fmt.Sprintf(
		"product-search:%s:%s:%s",
		tenantCode,
		countryCode,
		shortHash,
	)
}
