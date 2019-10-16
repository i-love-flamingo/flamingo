
# Flamingo
[![Go Report Card](https://goreportcard.com/badge/github.com/i-love-flamingo/flamingo)](https://goreportcard.com/report/github.com/i-love-flamingo/flamingo) [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo) [![Build Status](https://travis-ci.org/i-love-flamingo/flamingo.svg)](https://travis-ci.org/i-love-flamingo/flamingo)

## What is Flamingo

Flamingo is a high productivity go based framework for building fast and pluggable web projects.
It is used to build scalable and maintainable (web)applications.

Flamingo is:

* open source
* written in go
* easy to learn
* fast and flexible

Go as simple, powerful and typesafe language is great to implement and scale serverside logic.
Flamingo has a clean architecture with clear dependencies in mind and offers a typical features and support for nowadays web applications:

* Configuration with yml - supporting multiple areas and contexts
* Powerful Dependency Injection
* a Module Concept
* Authentication Concepts and Security Middleware
* Flexible Routing with support for Prefix Routes and Reverse Routing
* Controller: Request / Response / Form Handling 
* Operational Readyness: Logging, (distributed) Tracing, Metrics and Healthchecks with seperate endpoint
* Localisation support
* Commands
* Sessionmanagement
* Powerful templating support with PUG - supporting mixins and easy scripting
* GraphQL support and therefore support to build nice SPA and PWAs on top of it
* Resilience and Caching to external APIs

Flamingo itself does not contain ORM Mapper or libraries - instead it emphasizes "ports and adapters" architecture - so that you have a technology free (domain) model and any possible (and replaceable) persitence behind it.
That makes Flamingo useful to build microservices and applications - especially to build "frontends" or portals that require interaction with other (micro) services in a distributed architecture. 
When sticking to the architectural recommendation you can build modular applications with replaceable adapters that gives you independed testability.

With **"Flamingo Commerce"** there is an additional active projects that offer rich and flexible features to build modern e-commerce applications.

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
