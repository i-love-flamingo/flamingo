# Flamingo Dev Hints

## Necessary tooling

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
