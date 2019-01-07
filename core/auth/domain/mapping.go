package domain

import (
	"encoding/gob"

	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/web"
	"github.com/coreos/go-oidc"
)

type (
	UserMapping struct {
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
	}

	UserMappingService struct {
		idTokenMapping config.Map
	}
)

func (ums *UserMappingService) Inject(config *struct {
	IdTokenMapping config.Map `inject:"config:auth.mapping.idToken"`
}) {
	ums.idTokenMapping = config.IdTokenMapping
}

func (ums *UserMappingService) UserFromIDToken(idToken *oidc.IDToken, session *web.Session) (*User, error) {
	var claims map[string]interface{}
	err := idToken.Claims(&claims)
	if err != nil {
		return nil, err
	}

	return ums.MapToUser(claims, session), nil
}

func (ums *UserMappingService) GetMapping(config.Map) (UserMapping, error) {
	var mapping UserMapping
	err := ums.idTokenMapping.MapInto(&mapping)

	return mapping, err
}

func (ums *UserMappingService) MapToUser(claims map[string]interface{}, session *web.Session) *User {
	mapping, err := ums.GetMapping(ums.idTokenMapping)
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
