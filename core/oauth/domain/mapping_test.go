package domain

import (
	"testing"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/suite"
)

type (
	UserMappingServiceTestSuite struct {
		suite.Suite

		mappingService *UserMappingService
	}
)

func TestUserMappingServiceTestSuite(t *testing.T) {
	suite.Run(t, &UserMappingServiceTestSuite{})
}

func (t *UserMappingServiceTestSuite) SetupTest() {
	t.mappingService = &UserMappingService{}
}

func (t *UserMappingServiceTestSuite) TestMapToUser_Default() {
	claims := map[string]interface{}{
		"sub":      "ID123456",
		"name":     "Mr. Awesome",
		"email":    "email@domain.com",
		"whatever": "whatever",
	}

	t.mappingService.idTokenMapping = config.Map{
		"sub":   "sub",
		"email": "email",
		"name":  "name",
	}

	t.Equal(&User{
		Sub:          "ID123456",
		Name:         "Mr. Awesome",
		Email:        "email@domain.com",
		CustomFields: map[string]string{},
		Type:         USER,
		Groups:       []string{""},
	}, t.mappingService.MapToUser(claims, web.EmptySession()))
}

func (t *UserMappingServiceTestSuite) TestMapToUser_AllMainFields() {
	claims := map[string]interface{}{
		"sub":         "ID123456",
		"name":        "Mr. Awesome",
		"email":       "email@domain.com",
		"salutation":  "mister",
		"firstName":   "Mr.",
		"lastName":    "Awesome",
		"street":      "some street",
		"zipCode":     "12345",
		"city":        "Whitecity",
		"dateOfBirth": "01.01.2000",
		"country":     "Mars",
		"groups":      "GROUP1,GROUP2,GROUP3",
		"whatever":    "whatever",
	}

	t.mappingService.idTokenMapping = config.Map{
		"sub":         "sub",
		"email":       "email",
		"name":        "name",
		"salutation":  "salutation",
		"firstName":   "firstName",
		"lastName":    "lastName",
		"street":      "street",
		"zipCode":     "zipCode",
		"city":        "city",
		"dateOfBirth": "dateOfBirth",
		"country":     "country",
		"groups":      "groups",
	}

	t.Equal(&User{
		Sub:          "ID123456",
		Name:         "Mr. Awesome",
		Email:        "email@domain.com",
		Salutation:   "mister",
		FirstName:    "Mr.",
		LastName:     "Awesome",
		Street:       "some street",
		ZipCode:      "12345",
		City:         "Whitecity",
		DateOfBirth:  "01.01.2000",
		Country:      "Mars",
		CustomFields: map[string]string{},
		Type:         USER,
		Groups:       []string{"GROUP1", "GROUP2", "GROUP3"},
	}, t.mappingService.MapToUser(claims, web.EmptySession()))
}

func (t *UserMappingServiceTestSuite) TestMapToUser_CustomFields() {
	claims := map[string]interface{}{
		"whatever": "value",
	}

	t.mappingService.idTokenMapping = config.Map{
		"customFields": config.Slice{"whatever"},
	}

	t.Equal(&User{
		CustomFields: map[string]string{
			"whatever": "value",
		},
		Groups: []string{""},
		Type:   USER,
	}, t.mappingService.MapToUser(claims, web.EmptySession()))
}

func (t *UserMappingServiceTestSuite) TestMapToUser_AllDifferent() {
	claims := map[string]interface{}{
		"someSub":         "ID123456",
		"someName":        "Mr. Awesome",
		"someEmail":       "email@domain.com",
		"someSalutation":  "mister",
		"someFirstName":   "Mr.",
		"someLastName":    "Awesome",
		"someStreet":      "some street",
		"someZipCode":     "12345",
		"someCity":        "Whitecity",
		"someDateOfBirth": "01.01.2000",
		"someCountry":     "Mars",
		"whatever":        "value",
		"customer_groups":        "RU",
	}

	t.mappingService.idTokenMapping = config.Map{
		"sub":          "someSub",
		"email":        "someEmail",
		"name":         "someName",
		"salutation":   "someSalutation",
		"firstName":    "someFirstName",
		"lastName":     "someLastName",
		"street":       "someStreet",
		"zipCode":      "someZipCode",
		"city":         "someCity",
		"dateOfBirth":  "someDateOfBirth",
		"country":      "someCountry",
		"customFields": config.Slice{"whatever"},
		"groups":       "userType;customer_groups",
	}

	t.Equal(&User{
		Sub:         "ID123456",
		Name:        "Mr. Awesome",
		Email:       "email@domain.com",
		Salutation:  "mister",
		FirstName:   "Mr.",
		LastName:    "Awesome",
		Street:      "some street",
		ZipCode:     "12345",
		City:        "Whitecity",
		DateOfBirth: "01.01.2000",
		Country:     "Mars",
		CustomFields: map[string]string{
			"whatever": "value",
		},
		Type:   USER,
		Groups: []string{"RU"},
	}, t.mappingService.MapToUser(claims, web.EmptySession()))
}
