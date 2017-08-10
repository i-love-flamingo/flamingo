package product

import (
	"context"
	"errors"
	"flamingo/core/product/domain"
	"flamingo/om3/searchperience/infrastructure/product/dto"
)

type (
	Mapper struct{}
)

func (ps *Mapper) Map(ctx context.Context, productDto *dto.Product, priceEngine PriceEngineService) (domain.BasicProduct, error) {
	var basicProduct domain.BasicProduct

	if productDto.ProductType == "simple" {
		basicProduct, err := ps.mapSimpleProduct(ctx, productDto, priceEngine)
		if err != nil {
			return nil, err
		}
		return basicProduct, nil
	}

	if productDto.ProductType == "configurable" {
		basicProduct, err := ps.mapConfigurableProduct(ctx, productDto, priceEngine)
		if err != nil {
			return nil, err
		}
		return basicProduct, nil
	}
	return basicProduct, errors.New("Unknown type or product format")
}

func (ps *Mapper) mapConfigurableProduct(ctx context.Context, productDto *dto.Product, priceEngine PriceEngineService) (domain.ConfigurableProduct, error) {
	configurableProduct := domain.ConfigurableProduct{}

	configurableProduct.BasicProductData = ps.dtoConfigurableToBaseData(&productDto.ConfigurableProduct)
	ps.addDtoProductDataToBaseData(productDto, &configurableProduct.BasicProductData)

	for _, variantDto := range productDto.Variants {
		variant := domain.Variant{}
		variant.SaleableData = ps.dtoVariantToSaleData(&variantDto)
		variant.BasicProductData = ps.dtoVariantToBaseData(&variantDto)
		configurableProduct.Variants = append(configurableProduct.Variants, variant)
		//TODO Remove when search has it
		priceinfo, err := priceEngine.TempRequestPriceEngine(ctx, variantDto)
		if err != nil {
			return configurableProduct, err
		}
		variant.ActivePrice = priceinfo
	}

	return configurableProduct, nil
}

func (ps *Mapper) mapSimpleProduct(ctx context.Context, productDto *dto.Product, priceEngine PriceEngineService) (domain.SimpleProduct, error) {
	simpleProduct := domain.SimpleProduct{}

	if len(productDto.Variants) < 1 {
		return simpleProduct, errors.New("No Variant in simple product returned in search response")
	}
	variant1 := productDto.Variants[0]
	simpleProduct.Identifier = productDto.ForeignID

	simpleProduct.BasicProductData = ps.dtoVariantToBaseData(&variant1)
	ps.addDtoProductDataToBaseData(productDto, &simpleProduct.BasicProductData)
	simpleProduct.SaleableData = ps.dtoVariantToSaleData(&variant1)
	simpleProduct.Teaser = ps.dtoTeaserToTeaser(productDto)
	//TODO Remove when search has it
	priceinfo, err := priceEngine.TempRequestPriceEngine(ctx, variant1)
	if err != nil {
		return simpleProduct, err
	}
	simpleProduct.ActivePrice = priceinfo

	return simpleProduct, nil
}

func (ps *Mapper) dtoVariantToBaseData(variant1 *dto.Variant) domain.BasicProductData {
	basicData := domain.BasicProductData{}
	basicData.Title = variant1.Title

	basicData.Attributes = domain.Attributes(variant1.Attributes)

	basicData.ShortDescription = variant1.ShortDescription
	basicData.Description = variant1.Description
	basicData.CreatedAt = variant1.CreatedAt
	basicData.MarketPlaceCode = variant1.MarketPlaceCode
	basicData.RetailerCode = variant1.RetailerCode

	for _, media := range variant1.Media {
		basicData.Media = append(basicData.Media, domain.Media(media))
	}
	return basicData
}

func (ps *Mapper) dtoConfigurableToBaseData(configurable *dto.ConfigurableProduct) domain.BasicProductData {
	basicData := domain.BasicProductData{}
	basicData.Title = configurable.Title

	basicData.Attributes = domain.Attributes(configurable.Attributes)
	// check if this fields are missing in search object
	//basicData.ShortDescription = configurable.ShortDescription
	//basicData.Description = configurable.Description
	basicData.CreatedAt = configurable.CreatedAt
	for _, media := range configurable.Media {
		basicData.Media = append(basicData.Media, domain.Media(media))
	}
	return basicData
}

func (ps *Mapper) addDtoProductDataToBaseData(productDto *dto.Product, basicData *domain.BasicProductData) {

	basicData.Keywords = productDto.Keywords

	basicData.Keywords = append(basicData.Keywords, productDto.KeywordsImportant...)

	basicData.VisibleFrom = productDto.VisibleFrom
	basicData.VisibleTo = productDto.VisibleTo
	basicData.UpdatedAt = productDto.UpdatedAt
	basicData.CategoryCodes = productDto.CategoryCodes
	basicData.CategoryPath = productDto.CategoryPath
}

func (ps *Mapper) dtoVariantToSaleData(variant1 *dto.Variant) domain.SaleableData {
	saleData := domain.SaleableData{}

	// TODO - get active price from new serach response.. for now we do seperate request to priceeingie
	//saleData.ActivePrice.Default = variant1.OriginPrice
	//saleData.ActivePrice.Currency = variant1.Currency

	saleData.IsSaleable = variant1.IsSaleable
	saleData.SaleableFrom = variant1.SaleableFrom
	saleData.SaleableTo = variant1.SaleableTo

	saleData.RetailerSku = variant1.RetailerSku

	return saleData

}

func (ps *Mapper) dtoTeaserToTeaser(productDto *dto.Product) domain.TeaserData {
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
