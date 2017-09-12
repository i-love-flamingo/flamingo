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

	"math"

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
			Default           float64  `json:"default"`
			Discounted        float64  `json:"discounted"`
			IsDiscounted      bool     `json:"isDiscounted"`
			DiscountText      string   `json:"discountText"`
			ActiveBase        float64  `json:"activeBase"`
			ActiveBaseAmount  float64  `json:"activeBaseAmount"`
			ActiveBaseUnit    string   `json:"activeBaseUnit"`
			CampaignRules     []string `json:"campaignRules"`
			DenyMoreDiscounts bool     `json:"denyMoreDiscounts"`
			Currency          string   `json:"currency"`
			Context           struct {
				CustomerGroup interface{} `json:"customerGroup"`
				ChannelCode   string      `json:"channelCode"`
				Locale        string      `json:"locale"`
			} `json:"context"`
		} `json:"activePrice"`
		AvailablePrices []struct {
			Default           float64  `json:"default"`
			IsDiscounted      bool     `json:"isDiscounted"`
			Discounted        float64  `json:"discounted"`
			DiscountText      string   `json:"discountText"`
			ActiveBase        float64  `json:"activeBase"`
			ActiveBaseAmount  float64  `json:"activeBaseAmount"`
			ActiveBaseUnit    string   `json:"activeBaseUnit"`
			CampaignRules     []string `json:"campaignRules"`
			DenyMoreDiscounts bool     `json:"denyMoreDiscounts"`
			Currency          string   `json:"currency"`
			Context           struct {
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

	price := 0.0
	if _price, ok := variant.Attributes["price"].(string); ok {
		price, _ = strconv.ParseFloat(_price, 64)
	}
	basePrice := 0.0
	if _price, ok := variant.Attributes["basePrice"].(string); ok {
		basePrice, _ = strconv.ParseFloat(_price, 64)
	}

	//log.Printf("Call to %v", u)
	priceEngineRequest := []byte(`{
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
        "default": ` + strconv.FormatFloat(price, 'f', 1, 64) + `,
        "currency": "EUR",
        "taxClass": "test",
        "base": ` + strconv.FormatFloat(basePrice, 'f', 1, 64) + `,
        "baseUnit": "` + variant.OriginBasePriceUnit + `",
        "baseAmount": ` + strconv.FormatFloat(variant.OriginBasePriceAmount, 'f', 1, 64) + `,
        "special": ` + strconv.FormatFloat(variant.SpecialPrice, 'f', 1, 64) + `,
        "specialTo": null,
        "specialFrom": null,
        "groupPrices": [
        ]
      }
    }]}`)

	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(priceEngineRequest))
	if err != nil {
		panic(err)
	}
	log.Printf("Priceengine Request Url %s \n", u.String())
	//log.Printf("Priceengine Request Body %v \n", string(priceEngineRequest))
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Priceengine Error  %s \n", err)
		return priceinfo, err
	}
	var tempPriceEngineResponseDto TempPriceEngineResponseDto
	err = json.NewDecoder(resp.Body).Decode(&tempPriceEngineResponseDto)
	//log.Printf("Resp %=v", resp)
	if err != nil {
		log.Printf("Priceengine Error  %s \n", err)
		return priceinfo, errors.WithStack(err)
	}
	//log.Printf("Priceengine Response %=v", tempPriceEngineResponseDto)
	if len(tempPriceEngineResponseDto) != 1 {
		log.Printf("Priceengine Error  - len mismatch")
		return priceinfo, errors.New("Priceengine response has not exactly one array entry")
	}

	priceinfo.Default = tempPriceEngineResponseDto[0].ActivePrice.Default + 0.42
	//priceinfo.Discounted = tempPriceEngineResponseDto[0].ActivePrice.Discounted
	priceinfo.Discounted = math.Ceil((priceinfo.Default/100)*80) + 0.41
	priceinfo.DiscountText = tempPriceEngineResponseDto[0].ActivePrice.DiscountText
	priceinfo.CampaignRules = tempPriceEngineResponseDto[0].ActivePrice.CampaignRules
	priceinfo.Currency = tempPriceEngineResponseDto[0].ActivePrice.Currency
	priceinfo.ActiveBase = tempPriceEngineResponseDto[0].ActivePrice.ActiveBase
	priceinfo.ActiveBaseAmount = tempPriceEngineResponseDto[0].ActivePrice.ActiveBaseAmount
	priceinfo.IsDiscounted = tempPriceEngineResponseDto[0].ActivePrice.IsDiscounted
	priceinfo.IsDiscounted = true

	return priceinfo, nil
}
