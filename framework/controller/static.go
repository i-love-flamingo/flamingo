package controller

import (
	"context"
	"net/http"

	"flamingo.me/flamingo/v3/framework/web"
)

type fileResponse struct {
	r *web.Request
}

// Apply result by calling http.ServeFile
func (fr fileResponse) Apply(ctx context.Context, rw http.ResponseWriter) error {
	http.ServeFile(rw, fr.r.Request(), fr.r.Params["name"])
	return nil
}

// Static is a controller to handle file requests
type Static struct{}

// File returns a fileResponse which uses http.ServeFile to respond to the request
func (*Static) File(ctx context.Context, r *web.Request) web.Result {
	return fileResponse{r: r}
}
