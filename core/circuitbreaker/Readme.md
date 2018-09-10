# Circuit Breaker

The circuitbreaker package wraps an implementation of the [Circuit Breaker pattern](https://docs.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker) for usage within flamingo.

The used implementation is [sony/gobreaker](https://github.com/sony/gobreaker).

It can be used to wrap all kind of requests that are likely to fail.

## Basic usage

To use the circuit breaker, you will have to wrap your request into a `BreakableFunction`, which is a `func() (interface{}, error)`.

Then just create an instance of `CircuitBreaker` an pass the function to the `Execute` function.
Execute runs the given request if the CircuitBreaker accepts it.
It returns an error instantly if the CircuitBreaker rejects the request.
Otherwise, it returns the result of the request.
If a panic occurs in the request, the CircuitBreaker handles it as an error and causes the same panic again.

## Settings

To create an instance of `CircuitBreaker` you will have to pass a `Settings` object:

* `Name` is the name of the CircuitBreaker.
* `MaxRequests` is the maximum number of requests allowed to pass through when the CircuitBreaker is half-open. If MaxRequests is 0, CircuitBreaker allows only 1 request.
* `Interval` is the cyclic period of the closed state for CircuitBreaker to clear the internal Counts. If Interval is 0, CircuitBreaker doesn't clear the internal Counts during the closed state.
* `Timeout` is the period of the open state, after which the state of CircuitBreaker becomes half-open. If Timeout is 0, the timeout value of CircuitBreaker is set to 60 seconds.
* `OnStateChange` is called whenever the state of CircuitBreaker changes.
* `ReadyToTrip` is called with a copy of Counts whenever a request fails in the closed state. If ReadyToTrip returns true, CircuitBreaker will be placed into the open state. If ReadyToTrip is nil, readyToTrip from this package is used.

If you pass a `flamingo.Logger` into the constructor, the CircuitBreaker will log each state change. Pass `nil` if you do not want this.

## Example

```go
s := circuitbreaker.Settings{
  Name:          "Test",
  MaxRequests:   0,
  Interval:      10 * time.Second,
  Timeout:       5 * time.Second,
  OnStateChange: nil,
  ReadyToTrip:   nil,
}
b := circuitbreaker.NewCircuitBreaker(s, logger)

for i := 0; i < 100; i++ {
  time.Sleep(time.Second)
  fmt.Println("Get Request", i)

  body, err := b.Execute(
    func()(interface{}, error) {
      resp, err := http.Get("http://localhost:8081")
      if err != nil {
        return nil, err
      }
  
      defer resp.Body.Close()
      body, err := ioutil.ReadAll(resp.Body)
      if err != nil {
        return nil, err
      }
  
      return body, nil
    },
  )

  if err != nil {
    fmt.Println("CB ERR:", err)
  } else {
    fmt.Println("CB BODY:", string(body.([]byte)))
  }
}
```
