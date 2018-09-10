# Flamingo Dev Hints

## Setup

In order to work properly, Flamingo needs to be checked out at the following location: `$GOPATH/src/flamingo.me/flamingo`
```sh
git clone git@gitlab.aoe.com:shared/i-love-flamingo/flamingo.git $GOPATH/src/flamingo.me/flamingo
```

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

# Mockery (to create Mocks)

https://github.com/vektra/mockery

Usage:
mockery -name <Name of Interface>
