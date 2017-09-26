package category

import (
	"context"
	"encoding/json"
	"flamingo/core/category/domain"
	"flamingo/om3/searchperience/infrastructure"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

type (
	// Client for searchperience category requests
	Client struct {
		SearchperienceClient infrastructure.SearchperienceClient `inject:""`
	}

	// Service for flamingo categories
	Service struct {
		Client  *Client `inject:""`
		Locale  string  `inject:"config:locale"`
		Channel string  `inject:"config:searchperience.frontend.channel"`
	}

	searchCategory struct {
		Locale      string            `json:"locale"`
		Channel     string            `json:"channel"`
		Path        string            `json:"path"`
		CCode       string            `json:"code"`
		Media       []interface{}     `json:"media"`
		CName       string            `json:"name"`
		Content     string            `json:"content"`
		CCategories []*searchCategory `json:"categories"`
		active      bool
	}
)

var (
	_ domain.CategoryService = new(Service)
	_ domain.Category        = new(searchCategory)
)

// Code getter
func (sc *searchCategory) Code() string {
	return sc.CCode
}

// Name getter
func (sc *searchCategory) Name() string {
	return sc.CName
}

// Categories getter
func (sc *searchCategory) Categories() []domain.Category {
	categories := make([]domain.Category, len(sc.CCategories))
	for i, c := range sc.CCategories {
		categories[i] = c
	}
	return categories
}

// Active indicator
func (sc *searchCategory) Active() bool {
	return sc.active
}

// Get a category request
func (cc *Client) Get(ctx context.Context, category string, query url.Values) (*http.Response, error) {
	return cc.SearchperienceClient.Request(ctx, "categories/tree", query)
}

// Get a category object
func (cs *Service) Get(ctx context.Context, categoryCode string) (domain.Category, error) {
	query := url.Values{
		"locale":  {cs.Locale},
		"channel": {cs.Channel},
	}

	resp, err := cs.Client.Get(ctx, categoryCode, query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.NotFound
	}

	res := new(searchCategory)
	err = json.NewDecoder(resp.Body).Decode(res)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	markActive(res, categoryCode)

	return res, nil
}

func markActive(sc *searchCategory, categoryCode string) (marked bool) {
	for _, sub := range sc.CCategories {
		if markActive(sub, categoryCode) {
			sc.active = true
			return true
		}
	}
	if sc.CCode == categoryCode {
		sc.active = true
		return true
	}
	return
}
