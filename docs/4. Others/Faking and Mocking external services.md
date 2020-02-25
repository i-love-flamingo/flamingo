# Testing and Mocking external Services

## Introduction / Context
An important typical work in a Flamingo application is to consume external services - often via Restful APIs.

Flamingo acts then as a **Consumer** of Services provided by a **Provider**

During development and also in certain test steps - you usually want to fake or mock this external services. 

So lets clarify the naming:

* ***Faking:***  Means instead of calling the actual API we are faking a result internally, e.g. by using an implementation that dont calls an API but just answers with (fixed) faked results.
* ***Mocking:*** Means we try to use the real code that later (in production) should call the real service. But instead of using the real external service we call a separate mock. This mock needs to run.

Faking is easy and also allows for local independent testing.
Mocking has the advantage to test your application in a more production-like setup.

### Contract Testing
If you use a Mock, then someone should verify that the Mock behaves like the real service (or the other way around).
This is what is called contract testing.
You should establish this.

### Consumer Based Contract Testing
In case you have control over the provider you can use consumer based contract testing.
That means that a contract test provided by Flamingo is executed in the build pipeline of the provider.

We support pact for this.

## Faking in Flamingo
The [Ports and Adapters](../1. Flamingo Basics/4. Ports and Adapters.md) concept allows to register a "fake" implementation.

Adding the fake implementation:

```
###PROJECT###
└───fakeservices
    └─── servicepackagename
    │    └─── service.go
    module.go
```

In module.go, register the fake implementation:

```
injector.Override(new(productdomain.BrandService), "").To(product.FakeService{})
```

Better is to make this even configurable. And then add feature flags to your configuration in dev context (context_dev.yml).
This way you automatically use the fake implementation when starting Flamingo in with `CONTEXT=dev`.

## Mocking in Flamingo

### Getting started with Pact
Read more about Pact here: [docs.pact.io](https://docs.pact.io/)

Pact tests rely on [pact_go](https://github.com/pact-foundation/pact-go/) and Flamingo comes with a `testutil` package with useful functions to work with pact.

So you need to have the Pact daemon running:

1. [Install pact_go](https://github.com/pact-foundation/pact-go/#installation)
2. Run `pact-go daemon`. Now you have the daemon running on default port 6666

And can start writing PACT based tests... 

Checkout the example - e.g. in the example project  *"openweather"*

[todo]: <> (todo: deeplink openweather)

