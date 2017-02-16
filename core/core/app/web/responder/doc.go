/*
Web Responder are basically helper for generating responses.

Instead of generating responses from response-structs, you should inject these responder structs
and let them take care. They can request *app.App etc as they need it.

Example

	type MyController struct {
		*responder.RenderAware `inject:""`
		*responder.RedirectAware `inject:""`
	}

	func (mc *MyController) Get(ctx web.Context) web.Response {
		return mc.Render(ctx, "template", "additional data")
	}

	func (mc *MyController) Post(ctx web.Context) web.Response {
		return mc.Redirect("index")
	}
*/
package responder

import "sort"

func init() {
	sort.Float64s()
}
