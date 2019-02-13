s/\"flamingo.me\/flamingo\/framework\/dingo/\"flamingo.me\/dingo/
s/\"flamingo.me\/flamingo\/core\/canonicalUrl/\"flamingo.me\/v3\/core\/canonicalurl/
s/\"flamingo.me\/flamingo\/core\/requestTask/\"flamingo.me\/v3\/core\/requesttask/

s/\"flamingo.me\/flamingo\/framework\/event/\"flamingo.me\/v3\/framework\/flamingo/
s/ event.Event/ flamingo.Event/



s/\"flamingo.me\/framework\/web\/responder\"/\"flamingo.me\/v3\/framework\/web\"/
s/ responder.JSONAware/ Responder *web.Responder `inject:\"\"`/
s/ responder.RenderAware/ Responder *web.Responder `inject:\"\"`/
s/ responder.ErrorAware/ Responder *web.Responder `inject:\"\"`/
s/\.JSON(/\.Responder\.Data(/
s/\.Render(/\.Responder\.Render(/
s/\.Redirect(/\.Responder\.Redirect(/


s/\"flamingo.me\/flamingo\/framework\/router\"/\"flamingo.me\/v3\/framework\/web\"/
s/router\.Bind(/web\.Bind(/

s/\"flamingo.me\/flamingo\/framework\/template\"/\"flamingo.me\/v3\/framework\/flamingo\"/
s/template\.BindFunc(/flamingo\.BindTemplateFunc(/
s/template\.BindCtxFunc(/flamingo\.BindTemplateFunc(/


s/\"flamingo.me\/flamingo\/framework\/session\"/\"flamingo.me\/v3\/framework\/flamingo\"/
s/session\.Module(/flamingo\.SessionModule(/



s/\"flamingo.me\/flamingo\/core\/redirects/\"flamingo.me\/redirects\//
s/\"flamingo.me\/flamingo\/core\/pugtemplate/\"flamingo.me\/pugtemplate\//
s/\"flamingo.me\/flamingo\/core\/form2/\"flamingo.me\/form/
s/\"flamingo.me\/flamingo\/core\/form/\"flamingo.me\/form/

s/\"flamingo.me\/flamingo\/core\/csrf/\"flamingo.me\/csrf\//
s/\"flamingo.me\/flamingo\/core\/csp/\"flamingo.me\/csp\//
s/\"flamingo.me\/flamingo\/core\/captcha/\"flamingo.me\/captcha\//

s/\"flamingo.me\/flamingo\//"flamingo.me\/v3\//

s/\"flamingo.me\/flamingo-commerce\//"flamingo.me\/flamingo-commerce\/v3\//
