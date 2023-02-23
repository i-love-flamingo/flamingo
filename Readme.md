# Flamingo Framework

<img align="right" width="159px" src="https://raw.githubusercontent.com/i-love-flamingo/flamingo/master/docs/assets/flamingo-logo-only-pink-on-white.png">


[![Go Report Card](https://goreportcard.com/badge/github.com/i-love-flamingo/flamingo)](https://goreportcard.com/report/github.com/i-love-flamingo/flamingo) 
[![GoDoc](https://godoc.org/github.com/i-love-flamingo/flamingo?status.svg)](https://godoc.org/github.com/i-love-flamingo/flamingo)
[![Tests](https://github.com/i-love-flamingo/flamingo/workflows/Tests/badge.svg?branch=master)](https://github.com/i-love-flamingo/flamingo/actions?query=branch%3Amaster+workflow%3ATests)
[![Release](https://img.shields.io/github/release/i-love-flamingo/flamingo?style=flat-square)](https://github.com/i-love-flamingo/flamingo/releases)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/i-love-flamingo/flamingo)](https://www.tickgit.com/browse?repo=github.com/i-love-flamingo/flamingo)
[![Slack](https://img.shields.io/badge/Gophers_slack-%23flamingo-white?style=social&logo=slack&logoColor=E72064)](https://gophers.slack.com/archives/C04QYKWLVPD)


Flamingo is a web framework based on Go.  
It is designed to build pluggable and maintainable web projects.
It is production ready, field tested and has a growing ecosystem.


# Quick start

> See "examples/hello-world"

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

**Hello World Example:**

To extend this empty flamingo project with a "Hello World" output please create a new module "helloworld" like this:

```bash
mkdir helloworld
cat helloworld/module.go
``` 

With the following code in `module.go`:

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

This file now defines a very simple module, that can be used in the Flamingo bootstrap. 
In this case it registers a new handler that renders a simple "Hello World" message and binds the route "/" to this handler.
Now please include this new module in your existing `main.go` file:

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

And open http://localhost:3322 you will see your "Hello World!" output.



# Getting started
To learn more about Flamingo you can check out the full [hello-world example tutorial](https://github.com/i-love-flamingo/example-helloworld)
and read the documentation under [docs.flamingo.me](https://docs.flamingo.me/)

# Getting Help

The best way to ask a question is the [#flamingo channel](https://gophers.slack.com/messages/flamingo) on gophers.slack.com

If you are not yet in the Gophers slack, get your invitation here: https://invite.slack.golangbridge.org/

Other ways are:

* Ask in stackoverflow (we try to keep track of new questions)
* Write us an email: flamingo@aoe.com
* Open an issue in [github](https://github.com/i-love-flamingo/flamingo/issues) for bugs and feature requests


# Framework Details

## Feature List

* dependency injection with [Dingo](https://github.com/i-love-flamingo/dingo) 
* Flexible templating engines. (gotemplates and [pugtemplates](https://github.com/i-love-flamingo/pugtemplate))
* configuration concepts using [cue](https://cuelang.org/) with support for multiple config areas and additional config contexts
* A module concept for building modular and pluggable applications based on Dingo
* Authentication concepts and security middleware
* Flexible routing with support for prefix routes and reverse routing
* Web controller concept with request/response abstraction; form handling etc
* Operational readiness: logging, (distributed) tracing, metrics and healthchecks with separate endpoint
* Localisation support
* Commands using [Cobra](https://github.com/spf13/cobra)
* Event handling
* Sessionhandling and Management (By default uses [Gorilla](https://github.com/gorilla/sessions))

## Ecosystem

* GraphQL Module (and therefore support to build SPA and PWAs on top of it)
* Caching modules providing resilience and caching for external APIs calls.
* pugtemplate template engine for server side rendering with the related frontend tooling **[Flamingo Carotene](https://github.com/i-love-flamingo/flamingo-carotene)**
* **[Flamingo Commerce](https://github.com/i-love-flamingo/flamingo-commerce)**  is an active projects that offer rich and flexible features to build modern e-commerce applications.

