package cmspage

//go:generate go-bindata -pkg cmspage -prefix mocks/ mocks/

import (
	"context"
	"encoding/json"
	"fmt"

	"go.aoe.com/flamingo/core/cms/domain"
)

// FakePageService for CMS Pages
type FakePageService struct{}

// Get returns a CMS Page struct
func (ps *FakePageService) Get(ctx context.Context, name string) (*domain.Page, error) {
	var page domain.Page

	fmt.Println("Fake Page Service Call")
	if name == "bluefoot" {
		b, _ := Asset("service.cms.page.bluefoot.mock.json")
		json.Unmarshal(b, &page)
	} else {
		b, _ := Asset("service.cms.page.mock.json")
		json.Unmarshal(b, &page)
	}

	page.Identifier = name
	fmt.Println(page)

	return &page, nil
}
