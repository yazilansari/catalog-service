package repository

import (
	"catalog-service/internal/database"
	"catalog-service/internal/tenant/model"
)

func FindTenantByDomain(domain string) (*model.Tenant, error) {

	var tenant model.Tenant

	err := database.DB.
		Where("domain = ?", domain).
		Where("active = ?", true).
		First(&tenant).Error

	if err != nil {
		return nil, err
	}

	return &tenant, nil
}
