# Gotemplate

Flamingo comes with a wrapped `html/template` as simple default template engine.

## Structured templating

### Template directory

This module allows you to set up a deeply nested directory structure containing template files with the `.html` type ending. 
These files can be referenced from a controller by just using the path without `.html`.

An example for such a directory structure could be:

```text
/project/
  /templates/
    /deep/
      /nested/
        /index.html
  main.go
  go.mod
```

Here `index.html` resembles an example template. For our purposes, it contains the following content:

```html
<!-- /templates/deep/nested/index.html -->
<!doctype html>
<html>
  <head>
      <meta charset="utf-8">
      <title>Hello World</title>
  </head>
  <body>
    <main>
      <h1>Huzzah! It works!</h1>
      <p>
        This is an example text.
      </p>
    </main>
  </body>
</html>
```

You can refer to the [html/template documentation](golang.org/pkg/html/template/) for further information on how to fill your template.

To render `index.html` in your controller, just call

```go
return controller.responder.Render("templates/deep/nested/index")
```

You can also [register the `templates` directory as the directory containing all static assets](#static-assets), therefore reducing the render call to:

```go
return controller.responder.Render("deep/nested/index")
```

### Layout templates

Layouts can be used to reduce boilerplate html when creating templates by encapsulating your templates.

To begin, let's start by creating a new `layouts` folder and a `base.html` layout file. You can configure the location of your layouts folder in your project [configuration](#configuration).

```text
/project/
  /templates/
    /deep/
      /nested/
        /index.html
    /layouts/
      base.html
  main.go
  go.mod
```

A layout contains all the html that you want to reuse. Therefore, we first want to extract all the boilerplate html out of our templates and place it into our `base.html` and update our `index.html` accordingly.

Here is what our `index.html` looks like, after refactoring:

```html
<!-- /templates/deep/nested/index.html -->
{{template "layouts/base.html" .}}

{{define "title"}}
Hello World
{{end}}

{{define "content"}}
<main>
  <h1>Huzzah! It works!</h1>
  <p>
    This is an example text.
  </p>
</main>
{{end}}
 
```

Ok, let's look at what we did step by step:

1. We moved all the html we want to reuse into the `base.html` layout
2. We defined into which layout our template will be inserted into via the `{{template "layouts/base.html" .}}` action. (The dot after the path hands the data to the specified layout when everything is being rendered)
3. We defined our template blocks via the `{{define "<block-name>"}}` action and closed said definitoin with the `{{end}}` action.

Next, let's look at our layout file:

```html
<!-- /templates/layouts/base.html -->
<!doctype html>
<html>
  <head>
      <meta charset="utf-8">
      <title>{{template "title" .}}</title>
  </head>
  <body>
    {{template "content" .}}
  </body>
</html>
```

As you can see, this is where most of the html from `index.html` has ended up. You may also have notice that the previously defined template blocks have been called upon at there corresponding positions.

If you were to now render `index.html` you would recieve an html page like the one we started out with.

Congratulations! You have understood the basic concept of layouts, but that's not all!

You can make use of layouts and templates to create fragments which can then be called upon dynamically, like in this example:

```html
{{range $i, $p := .Products}}
  {{if le $i 0}}
    {{template "deep/otherNest/noProducts.html" $p}}
  {{end}}
  <div class="row">
    {{template "deep/otherNest/product.html" $p}}
  </div>
{{end}}
```

## Configuration

Within your `config.yml` you can define the paths for your template and layout directory.

```yml
# /config/config.yml
gotemplates:
  engine:
    templates:
      basepath: "templates", # template directory
    layout:
      dir: "layouts", #  layout directory within the template directory
```

# Static assets
You can use Flamingo’s built-in static file handler to automatically serve necessary static assets from your asset folder.

## Code-based configuration

You can set it up by adding a route and setting the `name` param to the name of your asset folder.

In the following example, our assets lie in the `asset` folder:

```go
// /polls/urls.go
func (u *urls) Routes(registry *web.RouterRegistry) {
    // ...
    registry.MustRoute("/asset/*name", `flamingo.static.file(name, dir?="asset")`)
}
```

Or via your routes.yml configuration file:

## routes.yml configuration

```yaml
# /config/routes.yml
- controller: flamingo.static.file(name, dir?="asset")
  path: /asset/*name
```

Then, set up a reference to the the url by adding the following script to your template:

````html
<!-- /templates/deep/nested/index.html -->
<script src="{{ url "flamingo.static.file" "name" "polls.js"}}"></script>
````

*This essentially calls the ‘flamingo.static.file(name=”polls.js”)’ command with the `dir` param set to its default, which you defined in your `routes.yml`.*
