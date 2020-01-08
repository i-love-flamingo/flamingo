package oauth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClaimSet_AddVoluntaryClaim(t *testing.T) {
	claimSet := &ClaimSet{}
	claimSet.AddVoluntaryClaim(TopLevelClaimUserInfo, "name")

	opt, err := claimSet.AuthCodeOption()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, `{k:claims v:{"userinfo":{"name":null}}}`, fmt.Sprintf("%+v", opt))
}

func TestClaimSet_AddClaimWithValue_Essential(t *testing.T) {
	claimSet := &ClaimSet{}
	claimSet.AddClaimWithValue(TopLevelClaimUserInfo, "name", true, "nameValue")

	opt, err := claimSet.AuthCodeOption()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, `{k:claims v:{"userinfo":{"name":{"essential":true,"value":"nameValue"}}}}`, fmt.Sprintf("%+v", opt))
}

func TestClaimSet_AddClaimWithValue_Voluntary(t *testing.T) {
	claimSet := &ClaimSet{}
	claimSet.AddClaimWithValue(TopLevelClaimUserInfo, "name", false, "nameValue")

	opt, err := claimSet.AuthCodeOption()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, `{k:claims v:{"userinfo":{"name":{"value":"nameValue"}}}}`, fmt.Sprintf("%+v", opt))
}

func TestClaimSet_AddClaimWithValues_Essential(t *testing.T) {
	claimSet := &ClaimSet{}
	claimSet.AddClaimWithValues(TopLevelClaimIDToken, "email", true, "emailValue", "mailValue")

	opt, err := claimSet.AuthCodeOption()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, `{k:claims v:{"id_token":{"email":{"essential":true,"values":["emailValue","mailValue"]}}}}`, fmt.Sprintf("%+v", opt))
}

func TestClaimSet_AddClaimWithValues_Voluntary(t *testing.T) {
	claimSet := &ClaimSet{}
	claimSet.AddClaimWithValues(TopLevelClaimIDToken, "email", false, "emailValue", "mailValue")

	opt, err := claimSet.AuthCodeOption()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, `{k:claims v:{"id_token":{"email":{"values":["emailValue","mailValue"]}}}}`, fmt.Sprintf("%+v", opt))
}
