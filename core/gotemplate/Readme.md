# Gotemplate

Flamingo comes with a wrapped `html/template` as simple default template engine.

Refer to [golang.org/pkg/html/template/](https://golang.org/pkg/html/template/) for the basic documentation

## Structured templating

### Template directory

This module allows to set up a deeply nested directory structure with template (html) files. 
These files can be referenced from a controller by just using the path without `.html`.

For example, to render `deep/nested/index.html` in your controller, just call

```go
return controller.responder.Render("deep/nested/index")
```

### Layout templates

In addition, a set of base layout templates can be defined in a separate directory. These layout templates can be included
into all rendered templates.

If you want to define a site template, just call different sub templates inside like

```gotemplate
{{template "content" .}}
```

In your rendered template, you can call the layout template and define all needed blocks:

```gotemplate
{{template "pages/site.html" .}}

{{define "content"}}
  <h1>The site content</h1>
{{end}}
```

The layout templates can also be used to define common "snippets" which can be used in every rendered template, for example:

```gotemplate
{{range $i, $p := .Products}}
  {{if gt $i 0}}
    <hr/>
  {{end}}
  <div class="row">
    {{template "blocks/product.html" $p}}
  </div>
{{end}}

```

## Configuration

```yaml
gotemplates:
  engine:
    templates:
      basepath: "templates", # template directory
    layout:
      dir: "layouts", #  layout directory within the template directory
```

# Static assets
You can use Flamingo’s built-in static file handler to serve static assets that you need in your template.

Set it up by adding a route with a param called “name” (that will get the name of the asset), such as 

```
polls/urls.go:
```

```go
func (u *urls) Routes(registry *web.RouterRegistry) {
    // ...
    registry.MustRoute("/asset/*name", `flamingo.static.file(name, dir?="asset")`)
}
```

Or via the routes.yml feature:

Edit your file `config/routes.yml`

```yaml
- controller: flamingo.static.file(name, dir?="asset")
  path: /asset/*name
```

Then, in your template, create the url by doing:
````html
<script src="{{ url "flamingo.static.file" "name" "polls.js"}}"></script>
````
(essentially calling ‘flamingo.static.file(name=”polls.js”)’. The “dir” variable has been defined to default to “asset” in the registration.

This way flamingo will automatically serve the asset form your assets folder.