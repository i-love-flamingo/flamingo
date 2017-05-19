# Testing

## Go Testing

In go you place your tests directly in the package.
You can simply use the standard go testing tool.

To run tests of a certain package simply run the `go test` tool.
For example:

```
go test -v flamingo/core/cart/domain/
```

## Ginkgo

Most of Flamingo's package use [Gingko](https://github.com/onsi/ginkgo) for BDD based testing.

Take a look at the Dingo test suite to get an idea about the tests.
