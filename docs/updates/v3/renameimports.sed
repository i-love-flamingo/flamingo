s/\"flamingo.me\/flamingo\/framework\/dingo/\"flamingo.me\/dingo/
s/\"flamingo.me\/flamingo\/core\/canonicalUrl/\"flamingo.me\/flamingo\/v3\/core\/canonicalurl/
s/\"flamingo.me\/flamingo\/core\/requestTask/\"flamingo.me\/flamingo\/v3\/core\/requesttask/

s/\"flamingo.me\/flamingo\/framework\/event/\"flamingo.me\/flamingo\/v3\/framework\/flamingo/
s/ event.Event/ flamingo.Event/


s/\"flamingo.me\/flamingo\/core\/redirects/\"flamingo.me\/redirects/
s/\"flamingo.me\/flamingo\/core\/pugtemplate/\"flamingo.me\/pugtemplate/
s/\"flamingo.me\/flamingo\/core\/form2/\"flamingo.me\/form/
s/\"flamingo.me\/flamingo\/core\/form/\"go.aoe.com\/flamingo\/form/

s/\"flamingo.me\/flamingo\/core\/csrf/\"flamingo.me\/csrf\//
s/\"flamingo.me\/flamingo\/core\/csp/\"flamingo.me\/csp\//
s/\"flamingo.me\/flamingo\/core\/captcha/\"flamingo.me\/captcha\//


s/\"flamingo.me\/flamingo\/framework\/web\/responder/\"flamingo.me\/flamingo\/v3\/framework\/web/
s/responder.JSONAware/Responder *web.Responder/
s/responder.RenderAware/Responder *web.Responder/
s/responder.ErrorAware/Responder *web.Responder/
s/\.JSON([a-z]*,/\.Responder\.Data(/
s/\.Render([a-z]*,/\.Responder\.Render(/
s/\.Redirect([a-z]*,/\.Responder\.Redirect(/
s/ web\.Response/ web\.Result/
s/\.Param1(\(.*\))/\.Params\[\1]/


s/\"flamingo.me\/flamingo\/framework\/router\"/\"flamingo.me\/flamingo\/v3\/framework\/web\"/
s/router\.Bind(/web\.BindRoutes(/
s/router\.Router/web\.Router/
s/router\.Registry/web\.RouterRegistry/



s/\"flamingo.me\/flamingo\/framework\/template\"/\"flamingo.me\/flamingo\/v3\/framework\/flamingo\"/
s/template\.BindFunc(/flamingo\.BindTemplateFunc(/
s/template\.BindCtxFunc(/flamingo\.BindTemplateFunc(/


s/\"flamingo.me\/flamingo\/framework\/session\"/\"flamingo.me\/flamingo\/v3\/framework\/flamingo\"/
s/session\.Module(/flamingo\.SessionModule(/

s/\"flamingo.me\/flamingo\/core\/cmd/"flamingo.me\/flamingo\/v3\/framework\/cmd/

s/\"flamingo.me\/flamingo\//"flamingo.me\/flamingo\/v3\//
s/\"flamingo.me\/flamingo\/v3\/v3\//"flamingo.me\/flamingo\/v3\//


s/\"flamingo.me\/flamingo-commerce\//"flamingo.me\/flamingo-commerce\/v3\//
