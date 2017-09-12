# Flamingo

## What is Flamingo

Flamingo is a go based, opinionated framework for pluggable web projects.

## Necessary tooling

### dep

Dependency manager:

`go get -u github.com/golang/dep/cmd/dep`

Usage

`dep ensure` or `dep ensure -vendor-only`

### go-bindata

Static file compiler for fakeservices, etc

`go get -u github.com/jteeuwen/go-bindata/...`

## Docs

To read the documentation:

```
make docs
```

This will start building and previewing the mkdocs based documentation in a Docker container.

To view the docs open  [Docs](http://localhost:8000)

-----------------
