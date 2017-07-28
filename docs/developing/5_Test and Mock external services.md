# Testing and Mocking external Services

## Introduction / Context
An important role of Flamingo is to consume external services - often via REST APIs.

Flamingo is a **Consumer** of Services provided by a **Provider**

During development you usally want to fake or mock this external services. So lets clearify the naming:
* Faking:  Means instead of Calling the actual API we are faking a result internal, e.g. by using an implementation that dont calls an API but just answers with (fixed) faked results.
* Mocking: Means we try to use the real code that later (in production) should call the real service. But instead of using the real service we call a seperate mock. This mock need to run.

Faking is good for start. Mocking is the right way if you want more production like development setup.

## Contract Testing 
If you use a Mock, then someone should verify that the Mock behaves like the real Service (or the other way around).
This is what is called contract testing. 
You should establish this.

## Consumer Based Contract Testing
In case you have control over the provider you can use consumer based contract testing. 
That means that a contract test provided by flamingo is executed in the build pipeline of the provider.

We support pact for this.

# Faking in Flamingo
The [Ports and Adapters](2_Ports and Adapters.md) concept allows to register a "fake" implementation.

Adding the fake implementation:

```
###PROJECT###
│---internalmock
     |--- servicepackagename
     |     |---- service.go
     |module.go  
```

In module.go register the fake implementation:

```
injector.Override((*productdomain.BrandService)(nil), "").To(product.FakeService{})
```

Better is to make this even configurable. And then add feature flags to your configuration in dev context (context_dev.yml).
This way you automatically use the fake implementation when starting flamingo in CONTEXT dev.

# Mocking in Flamingo

## Getting started with Pact
Read more about Pact here: https://docs.pact.io/

Pact tests rely on pact_go (https://github.com/pact-foundation/pact-go/) and flamingo comes with a "testutil" package with useful functions to work with pact.

So you might need to have the Pact daemon running:

1. Install pact_go ( https://github.com/pact-foundation/pact-go/#installation )
2. Run `pact-go daemon`. Now you have the daemon running on default port 6666




