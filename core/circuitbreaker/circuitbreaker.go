package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
	"go.aoe.com/flamingo/framework/flamingo"
)

type (
	CircuitBreaker interface {
		Execute(function BreakableFunction) (interface{}, error)
	}

	goBreaker struct {
		Name        string
		MaxRequests uint32          `inject:"config:circuitbreaker.maxrequests,optional"`
		Interval    time.Duration   `inject:"config:circuitbreaker.interval,optional"`
		Timeout     time.Duration   `inject:"config:circuitbreaker.timeout,optional"`
		Logger      flamingo.Logger `inject:""`
		breaker     *gobreaker.CircuitBreaker
	}
	BreakableFunction func() (interface{}, error)
)

func NewCircuitBreaker(name string, logger flamingo.Logger) CircuitBreaker {
	var settings gobreaker.Settings
	b := &goBreaker{
		Name:   name,
		Logger: logger,
	}

	settings.Name = name
	settings.ReadyToTrip = readyToTrip
	settings.OnStateChange = b.onStateChange
	settings.Interval = time.Second * 5
	settings.Timeout = time.Second * 10
	b.breaker = gobreaker.NewCircuitBreaker(settings)

	return b
}

func readyToTrip(counts gobreaker.Counts) bool {
	failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
	return counts.Requests >= 3 && failureRatio >= 0.6
}

func (b *goBreaker) onStateChange(name string, from, to gobreaker.State) {
	b.Logger.Printf("Circuit breaker '%s' changed state from %s to %s", name, from.String(), to.String())
}

func (b *goBreaker) Execute(f BreakableFunction) (interface{}, error) {
	return b.breaker.Execute(f)
}
