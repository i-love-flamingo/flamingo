# OM3 Pug Templating

## Flamingo

OM3 uses the [core/pug_template](../Framework/pug-template.go) package for complex templates.

The frontend pipeline is in `akl/frontend`, and uses `yarn` for dependency management,
build tasks, etc.

Flamingo integrates with webpack's dev-server, so the easiest thing
for live-working with flamingo and pug is to run `yarn run dev`.
This way you get live template reloading.

## Static Rendering

OM3 templates are setup to be compiled to static HTML files.

If you run `yarn run dev` it is possible to access these master template via [WebPack](http://localhost:1337).
