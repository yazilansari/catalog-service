package redis

import "time"

const (
	CategorySubCategoryTTL = 1 * time.Hour
	PageTTL                = 1 * time.Hour
	ProductTTL             = 30 * time.Minute
	PLPTTL                 = 15 * time.Minute
	SearchTTL              = 10 * time.Minute
	SEOTTL                 = 24 * time.Hour
)
