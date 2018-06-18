package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
	"flamingo.me/flamingo/framework/flamingo"
)

type (
	// CircuitBreaker represents a state machine to prevent sending requests that are likely to fail.
	CircuitBreaker interface {
		Execute(function BreakableFunction) (interface{}, error)
	}

	// State represents the current state of the CircuitBreaker
	State interface {
		String() string
	}

	// Settings is used to pass settings to the CircuitBreaker
	//
	// Name is the name of the CircuitBreaker.
	//
	// MaxRequests is the maximum number of requests allowed to pass through
	// when the CircuitBreaker is half-open.
	// If MaxRequests is 0, the CircuitBreaker allows only 1 request.
	//
	// Interval is the cyclic period of the closed state
	// for the CircuitBreaker to clear the internal Counts.
	// If Interval is 0, the CircuitBreaker doesn't clear internal Counts during the closed state.
	//
	// Timeout is the period of the open state,
	// after which the state of the CircuitBreaker becomes half-open.
	// If Timeout is 0, the timeout value of the CircuitBreaker is set to 60 seconds.
	//
	// OnStateChange is called whenever the state of the CircuitBreaker changes.
	//
	// ReadyToTrip is called with a copy of Counts whenever a request fails in the closed state.
	// If ReadyToTrip returns true, the CircuitBreaker will be placed into the open state.
	// If ReadyToTrip is nil, readyToTrip from this package is used.
	Settings struct {
		Name          string
		MaxRequests   uint32
		Interval      time.Duration
		Timeout       time.Duration
		OnStateChange StateChangeFunction
		ReadyToTrip   func(counts gobreaker.Counts) bool
	}

	// BreakableFunction is a request which will be performed inside the CircuitBreaker
	BreakableFunction func() (interface{}, error)

	// StateChangeFunction is called when the CircuitBreaker changes its state
	StateChangeFunction func(cb CircuitBreaker, from, to State)

	goBreaker struct {
		Logger          flamingo.Logger
		breaker         *gobreaker.CircuitBreaker
		stateChangeFunc StateChangeFunction
	}
)

// NewCircuitBreaker returns a new CircuitBreaker configured with the given Settings.
func NewCircuitBreaker(settings Settings, logger flamingo.Logger) CircuitBreaker {
	b := &goBreaker{
		Logger:          logger,
		stateChangeFunc: settings.OnStateChange,
	}

	rTT := readyToTrip
	if settings.ReadyToTrip != nil {
		rTT = settings.ReadyToTrip
	}

	b.breaker = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:          settings.Name,
		MaxRequests:   settings.MaxRequests,
		Interval:      settings.Interval,
		Timeout:       settings.Timeout,
		ReadyToTrip:   rTT,
		OnStateChange: b.onStateChange,
	})

	return b
}

func readyToTrip(counts gobreaker.Counts) bool {
	failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
	return counts.Requests >= 3 && failureRatio >= 0.6
}

func (b *goBreaker) onStateChange(name string, from, to gobreaker.State) {
	if b.Logger != nil {
		b.Logger.Info("Circuit breaker '%s' changed state from %s to %s", name, from.String(), to.String())
	}
	if b.stateChangeFunc != nil {
		b.stateChangeFunc(b, from, to)
	}
}

// Execute runs the given request if the CircuitBreaker accepts it.
// Execute returns an error instantly if the CircuitBreaker rejects the request.
// Otherwise, Execute returns the result of the request.
// If a panic occurs in the request, the CircuitBreaker handles it as an error
// and causes the same panic again.
func (b *goBreaker) Execute(f BreakableFunction) (interface{}, error) {
	return b.breaker.Execute(f)
}
