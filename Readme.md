
# Flamingo Framework

<img align="right" width="159px" src="https://raw.githubusercontent.com/i-love-flamingo/flamingo/master/docs/assets/flamingo-logo-only-pink-on-white.png">


[![Go Report Card](https://goreportcard.com/badge/github.com/i-love-flamingo/flamingo)](https://goreportcard.com/report/github.com/i-love-flamingo/flamingo) [![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo) [![Build Status](https://travis-ci.org/i-love-flamingo/flamingo.svg)](https://travis-ci.org/i-love-flamingo/flamingo)


Flamingo is a web framework based on Go. It is designed to build pluggable and maintainable web projects.
It is production ready, field tested and has a growing ecosystem.


# Quick start

Initialize an empty project:

```bash
mkdir helloworld
cd helloworld
go mod init helloworld
```

Create your project main file:

```bash
cat main.go
``` 

```go
package main

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
)

func main() {
	flamingo.App([]dingo.Module{
	})
}
```

If you then start your project you will see a list of registered commands:
```bash
go run main.go
``` 

It will print something like:
```
Flamingo main

Usage:
  main [command]

Examples:
Run with -h or -help to see global debug flags

Available Commands:
  config      Config dump
  handler     Dump the Handlers and its registered methods
  help        Help about any command
  routes      Routes dump
  serve       Default serve command - starts on Port 3322

Flags:
  -h, --help   help for main

Use "main [command] --help" for more information about a command.
```

To start the server use the following sub command:

```bash
go run main.go serve
``` 

And open http://localhost:3322

## Hello World Example

Create a new module "helloworld":


```bash
mkdir helloworld
cat helloworld/module.go
``` 

With the following code for "module.go":

```go
package helloworld

import (
    	"context"
    	"net/http"
    	"strings"
    
    	"flamingo.me/dingo"
    	"flamingo.me/flamingo/v3/framework/web"
)

type Module struct{}

func (*Module) Configure(injector *dingo.Injector) {
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

This file now includes a very simple Module that can be used in the Flamingo bootstrap and binds new routes to the Flamingo router.
Now include this new module in your main.go file:

```go
package main

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3"
	"helloworld/helloworld"
)

func main() {
	flamingo.App([]dingo.Module{
        new(helloworld.Module),
	})
}
```

If you now run the server again 

```bash
go run main.go serve
``` 

And open http://localhost:3322 you will see your "hello world" string.



# Getting started

Check out the full [hello-world example](https://github.com/i-love-flamingo/example-helloworld)
and read the rendered documentation under [docs.flamingo.me](https://docs.flamingo.me/)

# Framework Details

## Feature List

* Flexible templating engines. (gotemplates and [pugtemplates](https://github.com/i-love-flamingo/pugtemplate))
* configuration concepts using [cue](https://cuelang.org/) with support for multiple config areas and additional config contexts
* dependency injection  [Dingo](https://github.com/i-love-flamingo/dingo) 
* A Module concept for building modular and pluggable applications based on Dingo
* Authentication concepts and security middleware
* Flexible routing with support for prefix routes and reverse routing
* Web Controller Support with: Request / Response / Form Handling etc
* Operational readiness: logging, (distributed) tracing, metrics and healthchecks with seperate endpoint
* Localisation support
* Commands
* Event handling
* Sessionhandling and Management (By default uses [Gorilla](https://github.com/gorilla/sessions))

## Ecosystem

* GraphQL Module (and therefore support to build SPA and PWAs on top of it)
* Caching modules providing resilience and caching for external APIs calls.
* pugtemplate template engine for server side rendering with the related frontend tooling **[Flamingo Carotene](https://github.com/i-love-flamingo/flamingo-carotene)**
* **[Flamingo Commerce](https://github.com/i-love-flamingo/flamingo-commerce)**  active projects that offer rich and flexible features to build modern e-commerce applications.

