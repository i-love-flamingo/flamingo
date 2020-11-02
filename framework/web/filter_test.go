package web

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockedFilter struct {
	name string
}

func (f *mockedFilter) Filter(context.Context, *Request, http.ResponseWriter, *FilterChain) Result {
	return new(Responder).Data(f.name)
}

type mockedPrioritizedFilter struct {
	mockedFilter
	priority int
}

func (f *mockedPrioritizedFilter) Priority() int {
	return f.priority
}

func TestNewFilterChain(t *testing.T) {
	testCases := []struct {
		filters []Filter
		sorted  []Filter
	}{
		{},
		{
			filters: []Filter{
				&mockedFilter{
					name: "first",
				},
				&mockedFilter{
					name: "second",
				},
				&mockedFilter{
					name: "third",
				},
			},
			sorted: []Filter{
				&mockedFilter{
					name: "first",
				},
				&mockedFilter{
					name: "second",
				},
				&mockedFilter{
					name: "third",
				},
			},
		},
		{
			filters: []Filter{
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "first",
					},
					priority: -1,
				},
				&mockedFilter{
					name: "second",
				},
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "third",
					},
					priority: 1,
				},
			},
			sorted: []Filter{
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "third",
					},
					priority: 1,
				},
				&mockedFilter{
					name: "second",
				},
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "first",
					},
					priority: -1,
				},
			},
		},
		{
			filters: []Filter{
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "first",
					},
					priority: -1,
				},
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "second",
					},
					priority: -1,
				},
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "third",
					},
					priority: 0,
				},
				&mockedFilter{
					name: "fourth",
				},
				&mockedFilter{
					name: "fifth",
				},
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "sixth",
					},
					priority: 1,
				},
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "seventh",
					},
					priority: 1,
				},
			},
			sorted: []Filter{
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "sixth",
					},
					priority: 1,
				},
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "seventh",
					},
					priority: 1,
				},
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "third",
					},
					priority: 0,
				},
				&mockedFilter{
					name: "fourth",
				},
				&mockedFilter{
					name: "fifth",
				},
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "first",
					},
					priority: -1,
				},
				&mockedPrioritizedFilter{
					mockedFilter: mockedFilter{
						name: "second",
					},
					priority: -1,
				},
			},
		},
	}

	for _, tc := range testCases {
		chain := NewFilterChain(nil, tc.filters...)
		assert.Equal(t, tc.sorted, chain.filters)
	}
}
