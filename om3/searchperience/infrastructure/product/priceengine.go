package product

import (
	"bytes"
	"context"
	"encoding/json"
	"flamingo/core/product/domain"
	"flamingo/om3/searchperience/infrastructure/product/dto"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
)

// TODO: Complete file to be replaced when Searchperience does the callto priceegine instead

type (
	PriceEngineService struct {
		TempPriceEngineBaseURL string `inject:"config:searchperience.priceengine.baseurl"`
	}

	TempPriceEngineResponseDto []struct {
		MarketplaceCode string `json:"marketplaceCode"`
		ActivePrice     struct {
			Default          float64 `json:"default"`
			Discounted       float64 `json:"discounted"`
			DiscountText     string  `json:"discountText"`
			ActiveBase       float64 `json:"activeBase"`
			ActiveBaseAmount int     `json:"activeBaseAmount"`
			ActiveBaseUnit   string  `json:"activeBaseUnit"`
			Context          struct {
				CustomerGroup interface{} `json:"customerGroup"`
				ChannelCode   string      `json:"channelCode"`
				Locale        string      `json:"locale"`
			} `json:"context"`
		} `json:"activePrice"`
		AvailablePrices []struct {
			Default          int     `json:"default"`
			Discounted       float64 `json:"discounted"`
			DiscountText     string  `json:"discountText"`
			ActiveBase       int     `json:"activeBase"`
			ActiveBaseAmount int     `json:"activeBaseAmount"`
			ActiveBaseUnit   string  `json:"activeBaseUnit"`
			Context          struct {
				CustomerGroup interface{} `json:"customerGroup"`
				ChannelCode   string      `json:"channelCode"`
				Locale        string      `json:"locale"`
			} `json:"context"`
		} `json:"availablePrices"`
	}
)

func (ps *PriceEngineService) TempRequestPriceEngine(ctx context.Context, variant dto.Variant) (domain.PriceInfo, error) {

	var priceinfo domain.PriceInfo
	//Set default baseurl
	if ps.TempPriceEngineBaseURL == "" {
		ps.TempPriceEngineBaseURL = "http://priceengine"
	}
	u, _ := url.Parse(ps.TempPriceEngineBaseURL)
	u.Path += "/prices"
	var brand string
	if val, ok := variant.Attributes["brand"]; ok {
		brand = val.(string)
	} else {
		brand = ""
	}

	var campaignCodesJson string
	campaignCodesJson = "[]"
	if val, ok := variant.Attributes["campaignCodes"]; ok {
		b, err := json.Marshal(val)
		if err == nil {
			campaignCodesJson = string(b)
		}
	}

	//log.Printf("Call to %v", u)
	body := []byte(`{
  "context": {
    "customerGroup": null,
    "channelCode": "mainstore",
    "locale": "de_DE"
  },
  "products": [
    {
      "marketplaceCode": "` + variant.MarketPlaceCode + `",
      "retailerCode": "` + variant.RetailerCode + `",
      "brandCode": "` + brand + `",
      "campaignCodes": ` + campaignCodesJson + `,
      "categoryCodes": [
        "winter"
      ],
      "pricedata": {
        "default": ` + strconv.FormatFloat(variant.OriginPrice, 'f', 1, 64) + `,
        "currency": "EUR",
        "taxClass": "test",
        "base": ` + strconv.FormatFloat(variant.OriginBasePrice, 'f', 1, 64) + `,
        "baseUnit": "` + variant.OriginBasePriceUnit + `",
        "baseAmount": ` + strconv.FormatFloat(variant.OriginBasePriceAmount, 'f', 1, 64) + `,
        "special": ` + strconv.FormatFloat(variant.SpecialPrice, 'f', 1, 64) + `,
        "specialTo": null,
        "specialFrom": null,
        "groupPrices": [
        ]
      }
    }]}`)

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	log.Printf("Priceengine Request Url %s \n", u.String())
	//log.Printf("Priceengine Request Body %v \n", string(body))
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return priceinfo, err
	}
	var tempPriceEngineResponseDto TempPriceEngineResponseDto
	err = json.NewDecoder(resp.Body).Decode(&tempPriceEngineResponseDto)
	//log.Printf("Resp %=v", resp)
	if err != nil {
		return priceinfo, errors.WithStack(err)
	}
	//log.Printf("Priceengine Response %=v", tempPriceEngineResponseDto)
	if len(tempPriceEngineResponseDto) != 1 {
		return priceinfo, errors.New("Priceengine response has not exactly one array entry")
	}
	priceinfo.Default = tempPriceEngineResponseDto[0].ActivePrice.Default
	return priceinfo, nil
}
