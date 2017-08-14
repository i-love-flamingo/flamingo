package cmspage

import (
	"context"
	"encoding/json"
	"flamingo/core/cms/domain"
	"fmt"
	"io/ioutil"
)

// FakePageService for CMS Pages
type FakePageService struct{}

// Get returns a CMS Page struct
func (ps *FakePageService) Get(ctx context.Context, name string) (*domain.Page, error) {
	var page domain.Page

	fmt.Println("Fake Page Service Call")
	if name == "bluefoot" {
		b, _ := ioutil.ReadFile("../om3/fakeservices/cmspage/service.cms.page.bluefoot.mock.json")
		json.Unmarshal(b, &page)
	} else {
		b, _ := ioutil.ReadFile("../om3/fakeservices/cmspage/service.cms.page.mock.json")
		json.Unmarshal(b, &page)
	}

	page.Identifier = name
	fmt.Println(page)
	return &page, nil
}
