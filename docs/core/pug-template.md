# Pug Template

## pug.js

[Pug](https://pugjs.org/api/getting-started.html) is a JavaScript template rendering engine.
 
## Flamingo Pug Template

Flamingo integrates with pug via a custom Compiler.

The very basic setup is to compile Pug's AST (Abstract Syntax Tree) to native Go templates.

Currently Flamingo uses a slightly modified version of `template/html` and `template/text` to support
certain easy development workflows.

## Dynamic JavaScript

The Pug Template engine compiles a subset of JavaScript (ES2015) to Go templates.
This allows frontend developers to use known workflows and techniques, instead of learning
a complete new template engine.

To make this possible Flamingo rewrites the JavaScript to go, on the fly.

## Supported JavaScript

### Standard Datatypes

```javascript
{"key": "value"}

"string" + "another string"

1 + 2

15 * 8

[1, 2, 3, 4, 5]
```

## Support Pug

### Interpolation

```pug
p Text #{variable} something #{1 + 2}
```

### Mixins

```pug
mixin mymixin(arg1, arg2="default)
  p.something(id=arg2)= arg1
  
+mymixing("foo")

+mymixin("foo", "bar")
```

### Loops

```pug
each value, index in  ["a", "b", "c"]
  p value #{value} at #{index}
```

## Debugging

Templates can be debugged via `/_pugtpl/debug?tpl=pages/product/view`
