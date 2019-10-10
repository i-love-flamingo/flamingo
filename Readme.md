
# Flamingo
[![Go Report Card](https://goreportcard.com/badge/github.com/i-love-flamingo/flamingo)](https://goreportcard.com/report/github.com/i-love-flamingo/flamingo) [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo) [![Build Status](https://travis-ci.org/i-love-flamingo/flamingo.svg)](https://travis-ci.org/i-love-flamingo/flamingo)

## What is Flamingo

Flamingo is a go based, framework for pluggable web projects.
It is used to build scalable and maintainable web-applications. 
It's architecture is especially useful to build "frontends" for your headless microservice architecture.


* open source 
* written in go
* easy to learn
* fast and flexible

Go as simple, powerful and typesafe language is great to implement and scale serverside logic.
Flamingo has a clean architecture and uses "Domain Driven Design" and "Ports and Adapters" Layering - with clean and clear dependencies in mind.

# Getting Started / Hello World

Initialize your Project

```bash
mkdir helloworld
cd helloworld
go mod init helloworld
```

Create `main.go`:
```go
package main

import (
	"context"
	"net/http"
	"strings"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"flamingo.me/flamingo/v3/core/requestlogger"
	"flamingo.me/flamingo/v3/framework/web"
)

func main() {
	flamingo.App([]dingo.Module{
		new(requestlogger.Module),
		new(module),
	})
}

type module struct{}

func (*module) Configure(injector *dingo.Injector) {
	web.BindRoutes(injector, new(routes))
}

type routes struct{}

func (*routes) Routes(registry *web.RouterRegistry) {
	registry.Route("/", "home")
	registry.HandleAny("home", indexHandler)
}

func indexHandler(ctx context.Context, req *web.Request) web.Result {
	return &web.Response{
		Status: http.StatusOK,
		Body:   strings.NewReader("Hello World!"),
	}
}
```

Start server
```bash
go run main.go serve
``` 

Open http://localhost:3322

# The Flamingo Ecosystem 
With "Flamingo Commerce" and "Flamingo Carotene" you get your toolkit for building **Blazing fast commerce experience layers**

## Flamingo Commerce

* Contains flamingo modules that provide „domain“, „application“ and „interface“ logic around commerce features
* According to „ports and adapters“ these modules can be used with your own „Adapters“ to interact with any API or microservice you want.

## Flamingo Carotene
Is the frontend build pipeline featuring pug and atomic design pattern

It can be used to implement modern and blazing fast commerce web applications.

# Getting started

Check out the hello-world example
and read the rendered documentation under http://docs.flamingo.me

# External Links:
* http://www.flamingo.me
* http://docs.flamingo.me
