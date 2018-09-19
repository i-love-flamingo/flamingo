package redirects

import (
	"context"
	"errors"
	"net/http"

	"flamingo.me/flamingo/core/redirects/infrastructure"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	redirector struct {
		responder.RedirectAware
		responder.ErrorAware
		logger          flamingo.Logger
		redirectDataMap map[string]infrastructure.CsvContent
	}
)

func newRedirector(
	redirectAware responder.RedirectAware,
	errorAware responder.ErrorAware,
	logger flamingo.Logger,
	redirectData *infrastructure.RedirectData,
) *redirector {
	r := &redirector{
		RedirectAware:   redirectAware,
		ErrorAware:      errorAware,
		logger:          logger,
		redirectDataMap: make(map[string]infrastructure.CsvContent),
	}

	rd := redirectData.Get()

	for i := range rd {
		r.redirectDataMap[rd[i].OriginalPath] = rd[i]
	}

	for _, entry := range r.redirectDataMap {
		foundEntry, err := r.findEntryInRedirectMap(entry.RedirectTarget)
		if err == nil {
			logger.Error("ERROR: found a chained redirectData for ", foundEntry, " to ", foundEntry, " Rejecting redirects because of loop risk")
			r.redirectDataMap = nil

			break
		}
	}

	return r
}

// TryServeHTTP - implementation of OptionalHandler (from prefixrouter package)
func (r *redirector) TryServeHTTP(rw http.ResponseWriter, req *http.Request) (bool, error) {
	contextPath := req.RequestURI
	// r.Logger.Debug("TryServeHTTP called with %v", contextPath)
	status, location, err := r.processRedirects(contextPath)
	if err != nil {
		return true, errors.New("no redirect found")
	}
	if location != "" {
		rw.Header().Set("Location", location)
	}
	rw.WriteHeader(status)
	return false, nil
}

var _ router.Filter = (*redirector)(nil)

// Filter - implementation of Filter interface (from router package)
func (r *redirector) Filter(ctx context.Context, req *web.Request, w http.ResponseWriter, chain *router.FilterChain) web.Response {

	contextPath := req.Request().RequestURI

	status, location, err := r.processRedirects(contextPath)
	if err != nil {
		return chain.Next(ctx, req, w)
	}

	switch code := status; code {
	case http.StatusMovedPermanently:
		return r.RedirectPermanentURL(location)
	case http.StatusFound:
		return r.RedirectURL(location)
	case http.StatusGone:
		return r.ErrorWithCode(ctx, errors.New("page is gone"), http.StatusGone)
	case http.StatusNotFound:
		return r.ErrorNotFound(ctx, errors.New("page not found"))
	}

	return chain.Next(ctx, req, w) // never reached
}

// processRedirects - if a redirect is existing for given contextPath it returns the desired HTTP Response status and location (if relevant for the status) - if nothing is found it return 0
func (r *redirector) processRedirects(contextPath string) (int, string, error) {

	entry, err := r.findEntryInRedirectMap(contextPath)
	if err != nil {
		// nothing found for contextPath
		return 0, "", errors.New("contextPath not found")
	}

	r.logger.Debug("Redirecting from %s to %s by %d", entry.OriginalPath, entry.RedirectTarget, entry.HTTPStatusCode)

	switch code := entry.HTTPStatusCode; code {
	case http.StatusMovedPermanently, http.StatusFound:
		return code, entry.RedirectTarget, nil
	case http.StatusGone:
		return http.StatusGone, "", nil
	}

	// unsupported status - return 404 status
	return http.StatusNotFound, "", nil
}

func (r *redirector) findEntryInRedirectMap(contextPath string) (*infrastructure.CsvContent, error) {
	if r.redirectDataMap == nil {
		return nil, errors.New("no redirects loaded")
	}
	if currentRedirect, ok := r.redirectDataMap[contextPath]; ok {
		return &currentRedirect, nil
	}
	return nil, errors.New("not found")
}
