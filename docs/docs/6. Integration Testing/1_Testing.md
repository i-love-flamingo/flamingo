# Testing

## 1.) Testing with spock and gradle
### 1.1) Prerequisites
1.) Please be sure that groovy is installed on your machine. If it's not installed please download the latest groovy version (http://groovy-lang.org/download.html)
2.) 
### 1.2) 
In groovy/spock you place your tests in akl/integration-test/src/test/groovy
You can simply use the standard go testing tool.

To run tests of a certain package simply run the `go test` tool.
For example:

```
go test -v flamingo/core/cart/domain/
```

## Ginkgo

Most of Flamingo's package use [Gingko](https://github.com/onsi/ginkgo) for BDD based testing.

Take a look at the Dingo test suite to get an idea about the tests.


## Testing with Pact
Some tests are using Pact - so you need to prepare the test run:
Read on [5_Test and Mock external services.md](5_Test and Mock external services.md).

