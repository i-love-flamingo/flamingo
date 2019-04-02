package domain

import (
	"encoding/gob"
	"strings"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/coreos/go-oidc"
)

type (
	userMapping struct {
		Sub          string
		Name         string
		Email        string
		Salutation   string
		FirstName    string
		LastName     string
		Street       string
		ZipCode      string
		City         string
		DateOfBirth  string
		Country      string
		CustomFields []string
		Groups       string
	}

	// UserMappingService maps a user based on data available via the idTokenMapping setting
	UserMappingService struct {
		idTokenMapping config.Map
	}
)

// Inject dependencies
func (ums *UserMappingService) Inject(config *struct {
	IDTokenMapping config.Map `inject:"config:oauth.mapping.idToken"`
}) {
	ums.idTokenMapping = config.IDTokenMapping
}

// UserFromIDToken returns a user mapped with data from a provided OpenID connect token
func (ums *UserMappingService) UserFromIDToken(idToken *oidc.IDToken, session *web.Session) (*User, error) {
	var claims map[string]interface{}
	err := idToken.Claims(&claims)
	if err != nil {
		return nil, err
	}

	return ums.MapToUser(claims, session), nil
}

func (ums *UserMappingService) getMapping(config.Map) (userMapping, error) {
	var mapping userMapping
	err := ums.idTokenMapping.MapInto(&mapping)

	return mapping, err
}

// MapToUser returns the user mapped from the claims
func (ums *UserMappingService) MapToUser(claims map[string]interface{}, session *web.Session) *User {
	mapping, err := ums.getMapping(ums.idTokenMapping)
	if err != nil {
		panic(err)
	}

	claims = ensureClaims(claims, session)

	return &User{
		Sub:          ums.mapField(mapping.Sub, claims),
		Name:         ums.mapField(mapping.Name, claims),
		Email:        ums.mapField(mapping.Email, claims),
		Salutation:   ums.mapField(mapping.Salutation, claims),
		FirstName:    ums.mapField(mapping.FirstName, claims),
		LastName:     ums.mapField(mapping.LastName, claims),
		Street:       ums.mapField(mapping.Street, claims),
		ZipCode:      ums.mapField(mapping.ZipCode, claims),
		City:         ums.mapField(mapping.City, claims),
		DateOfBirth:  ums.mapField(mapping.DateOfBirth, claims),
		Country:      ums.mapField(mapping.Country, claims),
		CustomFields: ums.mapCustomFields(mapping.CustomFields, claims),
		Type:         USER,
		Groups:       ums.mapSliceField(mapping.Groups, claims),
	}
}

type cachedClaims string

const sessionkey cachedClaims = "cachedClaims"

func init() {
	gob.Register(sessionkey)
	gob.Register(map[string]interface{}{})
}

func ensureClaims(claims map[string]interface{}, session *web.Session) map[string]interface{} {
	var cached map[string]interface{}
	if raw, ok := session.Load(sessionkey); ok {
		cached, _ = raw.(map[string]interface{})
	}

	for k, v := range cached {
		if _, known := claims[k]; !known {
			claims[k] = v
		}
	}

	session.Store(sessionkey, claims)

	return claims
}

func (ums *UserMappingService) mapCustomFields(mapping []string, claims map[string]interface{}) map[string]string {
	custom := map[string]string{}

	for _, name := range mapping {
		value, ok := claims[name].(string)
		if !ok {
			continue
		}
		custom[name] = value
	}

	return custom
}

func (ums *UserMappingService) mapField(mappedFieldName string, claims map[string]interface{}) string {
	value, ok := claims[mappedFieldName].(string)
	if !ok {
		return ""
	}
	return value
}

func (ums *UserMappingService) mapSliceField(mappedFieldName string, claims map[string]interface{}) []string {
	return strings.Split(ums.mapField(mappedFieldName, claims), ",")
}
