package circuitbreaker

import (
	"errors"
	"testing"
	"time"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"flamingo.me/flamingo/core/logrus"
)

func TestGoBreakerExecute(t *testing.T) {
	tests := []struct {
		name                string
		settings            Settings
		tries               int
		expectedStatChanges int
		expectedResults     []interface{}
		expectedLogs        []string
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
			expectedLogs:        []string{},
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
			expectedLogs:        []string{"Circuit breaker 'Test Breaker' changed state from closed to open"},
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
			expectedLogs: []string{
				"Circuit breaker 'Test Breaker' changed state from closed to open",
				"Circuit breaker 'Test Breaker' changed state from open to half-open",
				"Circuit breaker 'Test Breaker' changed state from half-open to closed",
			},
			expectErrors: []bool{false, true, false, true, false},
		},
	}
	for _, testData := range tests {
		t.Run(testData.name, func(t *testing.T) {
			logger, hook := test.NewNullLogger()
			var stateChangedHandlerCalls = 0

			// a onStateChange function to count the state changes (and to assure it is used by the CircuitBreaker)
			testData.settings.OnStateChange = func() func(cb CircuitBreaker, from, to State) {
				return func(_ CircuitBreaker, _, _ State) {
					stateChangedHandlerCalls++
				}
			}()

			cb := NewCircuitBreaker(
				testData.settings,
				&logrus.LogrusEntry{Entry: logger.WithField("test", "test")},
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
			entries := hook.AllEntries()
			expectedLogMessages := testData.expectedLogs
			require.Len(t, entries, len(expectedLogMessages), "Number of Log messages not as expected")
			//for i := 0; i < len(entries); i++ {
			//assert.Equal(t, expectedLogMessages[i], entries[i].Message)
			//}
			//check the state changes
			assert.Equal(t, testData.expectedStatChanges, stateChangedHandlerCalls, "Number of expected state changes is wrong")
		})
	}
}
