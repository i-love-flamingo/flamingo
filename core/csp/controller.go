package csp

import (
	"encoding/json"
	"io/ioutil"

	"flamingo.me/flamingo/core/csrfPreventionFilter"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

type (
	// cspReportController shows information about a csp report
	cspReportController struct {
		Logger flamingo.Logger `inject:""`
	}
	report struct {
		CspReport struct {
			BlockedURI        string `json:"blocked-uri"`
			DocumentURI       string `json:"document-uri"`
			OriginalPolicy    string `json:"original-policy"`
			Referrer          string `json:"referrer"`
			SourceFile        string `json:"source-file"`
			ViolatedDirective string `json:"violated-directive"`
			ScriptSample      string `json:"script-sample"`
			LineNumber        int    `json:"line-number"`
		} `json:"csp-report"`
	}
)

// Post logs the csp report
func (dc *cspReportController) Post(ctx web.Context) web.Response {

	if ctx.Request().Header.Get("Content-Type") == "application/csp-report" {

		b, _ := ioutil.ReadAll(ctx.Request().Body)
		var data report
		json.Unmarshal(b, &data)
		l := dc.Logger.WithField("BlockedURI", data.CspReport.BlockedURI)
		l = l.WithField("DocumentURI", data.CspReport.DocumentURI)
		l = l.WithField("OriginalPolicy", data.CspReport.OriginalPolicy)
		l = l.WithField("Referrer", data.CspReport.Referrer)
		l = l.WithField("ScriptSample", data.CspReport.ScriptSample)

		l.Warn("csp policy report")

	}
	return &web.JSONResponse{}

}

// CheckOption takes care that the csrfPreventionFilter will be ignored
func (dc *cspReportController) CheckOption(option router.ControllerOption) bool {
	return option == csrfPreventionFilter.Ignore
}
