package dto

import (
	"context"
	"errors"
	"flamingo/core/product/domain"
	"strconv"
)

// Map a product response from searchperience
func Map(ctx context.Context, productDto *Product) (domain.BasicProduct, error) {
	if productDto.ProductType == domain.TypeSimple {
		basicProduct, err := mapSimpleProduct(ctx, productDto)
		if err != nil {
			return nil, err
		}
		return basicProduct, nil
	}

	if productDto.ProductType == domain.TypeConfigurable {
		basicProduct, err := mapConfigurableProduct(ctx, productDto)
		if err != nil {
			return nil, err
		}
		return basicProduct, nil
	}

	return nil, errors.New("unknown type or product format")
}

func mapConfigurableProduct(ctx context.Context, productDto *Product) (domain.ConfigurableProduct, error) {
	configurableProduct := domain.ConfigurableProduct{}

	configurableProduct.BasicProductData = dtoConfigurableToBaseData(&productDto.ConfigurableProduct)
	addDtoProductDataToBaseData(productDto, &configurableProduct.BasicProductData)
	configurableProduct.VariantVariationAttributes = productDto.VariantVariationAttributes

	for _, variantDto := range productDto.Variants {
		variant := domain.Variant{}
		variant.Saleable = dtoVariantToSaleData(&variantDto)
		variant.BasicProductData = dtoVariantToBaseData(&variantDto)
		configurableProduct.Variants = append(configurableProduct.Variants, variant)
	}

	return configurableProduct, nil
}

func mapSimpleProduct(ctx context.Context, productDto *Product) (domain.SimpleProduct, error) {
	simpleProduct := domain.SimpleProduct{}

	if len(productDto.Variants) < 1 {
		return simpleProduct, errors.New("No Variant in simple product returned in search response")
	}
	variant1 := productDto.Variants[0]
	simpleProduct.Identifier = productDto.ForeignID

	simpleProduct.BasicProductData = dtoVariantToBaseData(&variant1)
	addDtoProductDataToBaseData(productDto, &simpleProduct.BasicProductData)
	simpleProduct.Saleable = dtoVariantToSaleData(&variant1)
	simpleProduct.Teaser = dtoTeaserToTeaser(productDto)

	return simpleProduct, nil
}

func dtoVariantToBaseData(variant1 *Variant) domain.BasicProductData {
	basicData := domain.BasicProductData{}
	basicData.Title = variant1.Title

	basicData.Attributes = domain.Attributes(variant1.Attributes)

	basicData.ShortDescription = variant1.ShortDescription
	basicData.Description = variant1.Description
	basicData.CreatedAt = variant1.CreatedAt
	basicData.MarketPlaceCode = variant1.MarketPlaceCode
	basicData.RetailerCode = variant1.RetailerCode
	basicData.RetailerSku = variant1.RetailerSku

	haslist := false
	for _, media := range variant1.Media {
		basicData.Media = append(basicData.Media, domain.Media(media))
		if media.Usage == "list" {
			haslist = true
		}
	}
	if !haslist && len(variant1.Media) > 0 {
		media := domain.Media(variant1.Media[0])
		media.Usage = "list"
		basicData.Media = append(basicData.Media, media)
	}
	return basicData
}

func dtoConfigurableToBaseData(configurable *ConfigurableProduct) domain.BasicProductData {
	basicData := domain.BasicProductData{}
	basicData.Title = configurable.Title

	basicData.Attributes = domain.Attributes(configurable.Attributes)
	// check if this fields are missing in search object
	if basicData.ShortDescription == "" && configurable.Attributes["shortDescription"] != nil {
		basicData.ShortDescription = configurable.Attributes["shortDescription"].(string)
	}
	if basicData.Description == "" && configurable.Attributes["description"] != nil {
		basicData.Description = configurable.Attributes["description"].(string)
	}
	basicData.CreatedAt = configurable.CreatedAt
	for _, media := range configurable.Media {
		basicData.Media = append(basicData.Media, domain.Media(media))
	}
	return basicData
}

func addDtoProductDataToBaseData(productDto *Product, basicData *domain.BasicProductData) {

	basicData.Keywords = productDto.Keywords

	basicData.Keywords = append(basicData.Keywords, productDto.KeywordsImportant...)

	basicData.VisibleFrom = productDto.VisibleFrom
	basicData.VisibleTo = productDto.VisibleTo
	basicData.UpdatedAt = productDto.UpdatedAt
	basicData.CategoryCodes = productDto.CategoryCodes
	basicData.CategoryPath = productDto.CategoryPath
	basicData.MarketPlaceCode = productDto.MarketPlaceCode
}

func dtoVariantToSaleData(variant1 *Variant) domain.Saleable {
	saleData := domain.Saleable{}

	// TODO - get active price from new serach response.. for now we do seperate request to priceeingie
	//saleData.ActivePrice.Default = variant1.OriginPrice
	//saleData.ActivePrice.Currency = variant1.Currency

	saleData.IsSaleable = variant1.IsSaleable
	saleData.SaleableFrom = variant1.SaleableFrom
	saleData.SaleableTo = variant1.SaleableTo

	if p, ok := variant1.Attributes["price"]; ok {
		price, _ := strconv.ParseFloat(p.(string), 64)
		saleData.ActivePrice = domain.PriceInfo{
			Default: price,
		}
	}

	return saleData

}

func dtoTeaserToTeaser(productDto *Product) domain.TeaserData {
	teaserData := domain.TeaserData{}

	teaserData.ShortDescription = productDto.TeaserData.ShortDescription
	teaserData.Teaser = productDto.TeaserData.Teaser
	for _, media := range productDto.TeaserData.Media {
		teaserData.Media = append(teaserData.Media, domain.Media(media))
	}
	teaserData.Title = productDto.TeaserData.Title
	teaserData.ShortTitle = productDto.TeaserData.ShortTitle
	return teaserData
}
