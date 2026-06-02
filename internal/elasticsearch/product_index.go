package elasticsearch

import "strings"

func GetProductIndex(
	tenantCode string,
	countryCode string,
) string {

	return "products_" +
		strings.ToLower(tenantCode) +
		"_" +
		strings.ToLower(countryCode)
}
