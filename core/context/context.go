package context

type Context struct {
	Name    string
	BaseUrl string

	Parent *Context

	Configuration map[string]string

	Routes map[string]string

	Handler map[string]interface{}
}
