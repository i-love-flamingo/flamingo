# Contributing to Flamingo

## CLA

Flamingo is Open Source and therefore we need a signed CLA (Contributor License Agreements) before we can include your contributions.
See [AOE Individual Contributor License Agreement](CLA.md).

## Workflow

### Committing

Please contribute by opening merge requests in github.
This allows us to review, and optionally squash changes to an atomic change.

You can always create branches in this repository, not necessary to do so in a custom fork.

### Commit messages

A commit message should always start with the affected modules, such as:

`framework/dingo: add command line flags`
or
`core/pugtemplate/pugjs: update parser`

A commit message is supposed to tell what has changed.


### Releases

Flamingo follows Semantic Versioning. Minor versions are created any now and then. 

## Coding guidelines

There are a few things to note when coding, so please try to follow these guidelines.

### Generals

- go fmt applies
- go vet applies
- [Effective Go](https://golang.org/doc/effective_go.html) applies
- Always handle errors!
- Always verify type assertions with the second `_, ok := v.(T)` bool flag.

### Usage of context.Context

The Context is always the first argument to a method/function.
Do not check, and do not pass, nil contexts, as they are forbidden!
Do not use the context for general passing of scoped data, unless it's explicit necessary.

Use the opencensus/tags package for tracing/stats tags.

If necessary use the session.FromContext and web.FromContext to retrieve current session or request.
The `web.Request` has a `Values` map which allows you to pass request-scoped data. Please note that this works like a 
request-global variable, thus you might create implicit dependencies. Use a private type for the map-key.

### Interfaces

Limit interface usage to the necessary. There is no need to define a Service and a DefaultServiceImplementation, as long
as you do not plan to provide more than one implementation.

### Public interfaces and methods

Limit the public surface to the bare minimum. Use dingo-compatible `Inject(...)` methods rather than `inject:""` tags
on public members.

### Dingo

Dingo is mighty, always think about if/where/when to use it.
`Inject` methods are always to be used in favor of struct tags.

## Testing

Every public exposed Method/Function/Type needs appropriate unittests.

If possible please provide examples as well.

## Documentation / godoc

Every public exposed Method/Function/Type needs appropriate godoc comments.

If necessary you can write documentation in the docs folder.

As a rule of thumb:
- Code documentation belongs in the godoc.
- Functionality/Feature/Usage documentation belongs in the docs folder.
