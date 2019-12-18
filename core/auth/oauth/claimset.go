package oauth

import (
	"encoding/json"
	"golang.org/x/oauth2"
)

// JSON keys for top-level member of the Claims request JSON.
const (
	TopLevelClaimUserInfo = "userinfo"
	TopLevelClaimIdToken  = "id_token"
)

// claim describes JSON object used to specify additional information
// about the Claim being requested.
type claim struct {
	Essential bool     `json:"essential,omitempty"`
	Value     string   `json:"value,omitempty"`
	Values    []string `json:"values,omitempty"`
}

// ClaimSet contains map with members of the Claims request.
// It provides methods to add specific Claims.
type ClaimSet struct {
	claims map[string]map[string]*claim
}

// AuthCodeOption returns oauth2.AuthCodeOption with json encoded body containing all defined claims.
// It returns error in case when json marshalling can not be performed.
func (c *ClaimSet) AuthCodeOption() (oauth2.AuthCodeOption, error) {
	body, err := json.Marshal(c.claims)
	if err != nil {
		return nil, err
	}

	return oauth2.SetAuthURLParam("claims", string(body)), nil
}

// AddVoluntaryClaim adds the Claim being requested in default manner,
// as a Voluntary Claim.
func (c *ClaimSet) AddVoluntaryClaim(topLevelName string, claimName string) {
	c.initializeTopLevelMember(topLevelName)
	c.claims[topLevelName][claimName] = nil
}

// AddClaimWithValue adds the Claim being requested to return a particular value.
// The Claim can be defined as an Essential Claim.
func (c *ClaimSet) AddClaimWithValue(topLevelName string, claimName string, essential bool, value string) {
	c.initializeTopLevelMember(topLevelName)
	c.claims[topLevelName][claimName] = &claim{
		Essential: essential,
		Value:     value,
	}
}

// AddClaimWithValues adds the Claim being requested to return
// one of a set of values, with the values appearing in order of preference.
// The Claim can be defined as an Essential Claim.
func (c *ClaimSet) AddClaimWithValues(topLevelName string, claimName string, essential bool, values ...string) {
	c.initializeTopLevelMember(topLevelName)
	c.claims[topLevelName][claimName] = &claim{
		Essential: essential,
		Values:    values,
	}
}

// initializeTopLevelMember checks if top-level member is initialized
// as a map of JSON Claims objects and initializes it, if it's needed.
func (c *ClaimSet) initializeTopLevelMember(topLevelName string) {
	if c.claims == nil {
		c.claims = map[string]map[string]*claim{}
	}
	if c.claims[topLevelName] == nil {
		c.claims[topLevelName] = map[string]*claim{}
	}
}
