package csp

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"go.aoe.com/flamingo/core/csrfPreventionFilter"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"
)

type (
	cspFilter struct {
		responder.ErrorAware `inject:""`
		Router               *router.Router                      `inject:""`
		NonceGenerator       csrfPreventionFilter.NonceGenerator `inject:""`
		ReportMode           bool                                `inject:"config:cspFilter.reportMode"`
	}
)

// Filter realizes the Content-Security-Policy-Report and adds nonces to the script tags
func (f *cspFilter) Filter(ctx web.Context, w http.ResponseWriter, chain *router.FilterChain) web.Response {
	response := chain.Next(ctx, w)

	if cr, ok := response.(*web.ContentResponse); ok {
		originalBody, err := ioutil.ReadAll(cr.Body)
		if err != nil {
			return f.Error(ctx, err)
		}
		nonce := f.NonceGenerator.GenerateNonce()
		newTag := []byte("<script nonce=\"" + nonce + "\"")
		cr.Body = bytes.NewBuffer(bytes.Replace(originalBody, []byte("<script"), newTag, -1))

		url := f.Router.URL("_cspreport.view", router.P{})

		if f.ReportMode {
			w.Header().Add("Content-Security-Policy-Report-Only", "default-src 'self'; script-src 'self' 'nonce-"+nonce+"'; report-uri "+url.String()+"; style-src 'self' 'unsafe-inline'")

		} else {
			w.Header().Add("Content-Security-Policy", "default-src 'self'; script-src 'self' 'nonce-"+nonce+"'; report-uri "+url.String()+"; style-src 'self' 'unsafe-inline'")
		}
	}

	return response
}
