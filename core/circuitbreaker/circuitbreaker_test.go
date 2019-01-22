package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo/mocks"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
)

func TestGoBreakerExecute(t *testing.T) {
	tests := []struct {
		name                string
		settings            Settings
		tries               int
		expectedStatChanges int
		expectedResults     []interface{}
		expectedLogs        [][]interface{}
		expectErrors        []bool
	}{
		{
			name: "positive case (no state change)",
			settings: Settings{
				Name:          "Test Breaker",
				MaxRequests:   0,
				Interval:      0,
				Timeout:       0,
				OnStateChange: nil,
				ReadyToTrip:   nil,
			},
			tries:               5,
			expectedStatChanges: 0,
			expectedResults:     []interface{}{1, nil, 3, nil, 5},
			expectedLogs:        [][]interface{}{},
			expectErrors:        []bool{false, true, false, true, false},
		},
		{
			name: "switch to open state after second error",
			settings: Settings{
				Name:          "Test Breaker",
				MaxRequests:   0,
				Interval:      0,
				Timeout:       0,
				OnStateChange: nil,
				ReadyToTrip: func(counts gobreaker.Counts) bool {
					return counts.TotalFailures > 1
				},
			},
			tries:               5,
			expectedStatChanges: 1,
			expectedResults:     []interface{}{1, nil, 3, nil, nil},
			expectedLogs:        [][]interface{}{{"Circuit breaker '%s' changed state from %s to %s", "Test Breaker", "closed", "open"}},
			expectErrors:        []bool{false, true, false, true, true},
		},
		{
			name: "switch to open state after second error and back to closed",
			settings: Settings{
				Name:          "Test Breaker",
				MaxRequests:   1,
				Interval:      0,
				Timeout:       time.Millisecond,
				OnStateChange: nil,
				ReadyToTrip: func(counts gobreaker.Counts) bool {
					return counts.TotalFailures > 1
				},
			},
			tries:               5,
			expectedStatChanges: 3,
			expectedResults:     []interface{}{1, nil, 3, nil, 5},
			expectedLogs: [][]interface{}{
				{"Circuit breaker '%s' changed state from %s to %s", "Test Breaker", "closed", "open"},
				{"Circuit breaker '%s' changed state from %s to %s", "Test Breaker", "open", "half-open"},
				{"Circuit breaker '%s' changed state from %s to %s", "Test Breaker", "half-open", "closed"},
			},
			expectErrors: []bool{false, true, false, true, false},
		},
	}
	for _, testData := range tests {
		t.Run(testData.name, func(t *testing.T) {
			var logger = &mocks.Logger{}
			var stateChangedHandlerCalls = 0

			// a onStateChange function to count the state changes (and to assure it is used by the CircuitBreaker)
			testData.settings.OnStateChange = func() func(cb CircuitBreaker, from, to State) {
				return func(_ CircuitBreaker, _, _ State) {
					stateChangedHandlerCalls++
				}
			}()

			for _, expLog := range testData.expectedLogs {
				logger.On("Info", expLog...)
			}

			cb := NewCircuitBreaker(
				testData.settings,
				logger,
			)
			// The BreakableFunction simulates a request which will fail every second time
			f := func() BreakableFunction {
				i := 0
				return func() (interface{}, error) {
					i++
					if i%2 == 0 {
						return nil, errors.New("an error")
					}
					return i, nil
				}
			}()

			for try := 0; try < testData.tries; try++ {
				// give the test a change to hit the timeout, if we have one =)
				if testData.settings.Timeout > 0 {
					time.Sleep(testData.settings.Timeout)
				}

				res, err := cb.Execute(f)

				assert.Equal(t, testData.expectedResults[try], res, "result not as expected in try %d", try+1)
				if testData.expectErrors[try] {
					assert.Error(t, err)
				} else {
					assert.Nil(t, err, "An error was not expected in try %d", try+1)
				}

			}
			//check the log messages
			logger.AssertExpectations(t)

			//check the state changes
			assert.Equal(t, testData.expectedStatChanges, stateChangedHandlerCalls, "Number of expected state changes is wrong")
		})
	}
}
