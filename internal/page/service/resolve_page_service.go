package service

import (
	"catalog-service/internal/page/dto"
	"catalog-service/internal/page/repository"
)

func ResolvePage(
	tenantCode string,
	countryCode string,
	slug string,
) (*dto.ResolvePageResponse, error) {

	// =========================
	// CATEGORY
	// =========================

	_, err :=
		repository.FindCategoryBySlug(
			tenantCode,
			countryCode,
			slug,
		)

	if err == nil {

		return &dto.ResolvePageResponse{
			PageType: "category",

			Slug: slug,

			RedirectURL: "/category/" + slug,
		}, nil
	}

	// =========================
	// SUBCATEGORY
	// =========================

	subCategory, err :=
		repository.FindSubCategoryBySlug(
			tenantCode,
			countryCode,
			slug,
		)

	if err == nil {

		return &dto.ResolvePageResponse{
			PageType: "subcategory",

			Slug: subCategory.Slug,

			RedirectURL: "/subcategory/" + slug,
		}, nil
	}

	// =========================
	// PRODUCT
	// =========================

	product, err :=
		repository.FindProductBySlug(
			tenantCode,
			countryCode,
			slug,
		)

	if err == nil {

		return &dto.ResolvePageResponse{
			PageType: "product",

			Slug: product.Slug,

			RedirectURL: "/product/" + slug,
		}, nil
	}

	return nil, err
}
