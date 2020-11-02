package web

import (
	"context"
	"net/http"
	"sort"
)

type (
	// Filter is an interface which can filter requests
	Filter interface {
		Filter(ctx context.Context, req *Request, w http.ResponseWriter, fc *FilterChain) Result
	}

	// PrioritizedFilter is an interface which allows Filter to be prioritized
	// In case Filter is not PrioritizedFilter, default priority is 0
	PrioritizedFilter interface {
		Filter
		Priority() int
	}

	// FilterChain defines the chain which contains all filters which will be worked off
	FilterChain struct {
		filters   []Filter
		final     lastFilter // special case for the final controller
		postApply []func(err error, result Result)
	}

	lastFilter func(ctx context.Context, req *Request, w http.ResponseWriter) Result

	sortableFilter struct {
		filter Filter
		index  int
	}

	sortableFilers []sortableFilter
)

func (fnc lastFilter) Filter(ctx context.Context, req *Request, w http.ResponseWriter, chain *FilterChain) Result {
	return fnc(ctx, req, w)
}

// NewFilterChain constructs and sets the final filter and optional filters
func NewFilterChain(final lastFilter, filters ...Filter) *FilterChain {
	sortable := newSortableFilters(filters)
	sort.Sort(sort.Reverse(sortable))

	return &FilterChain{
		final:   final,
		filters: sortable.toFilters(),
	}
}

// Next calls the next filter and deletes it of the chain
func (fc *FilterChain) Next(ctx context.Context, req *Request, w http.ResponseWriter) Result {
	if len(fc.filters) == 0 {
		// filter chain ended
		return fc.final(ctx, req, w)
	}

	next := fc.filters[0]
	fc.filters = fc.filters[1:]
	return next.Filter(ctx, req, w, fc)
}

// AddPostApply adds a callback to be called after the response has been applied to the responsewriter
func (fc *FilterChain) AddPostApply(callback func(err error, result Result)) {
	fc.postApply = append(fc.postApply, callback)
}

// Len supports implementation for sort.Interface
func (sf sortableFilers) Len() int {
	return len(sf)
}

// Less supports implementation for sort.Interface
func (sf sortableFilers) Less(i, j int) bool {
	firstPriority := 0
	if filter, ok := sf[i].filter.(PrioritizedFilter); ok {
		firstPriority = filter.Priority()
	}

	secondPriority := 0
	if filter, ok := sf[j].filter.(PrioritizedFilter); ok {
		secondPriority = filter.Priority()
	}

	return firstPriority < secondPriority || (firstPriority == secondPriority && sf[i].index > sf[j].index)
}

// Swap supports implementation for sort.Interface
func (sf sortableFilers) Swap(i, j int) {
	sf[i], sf[j] = sf[j], sf[i]
}

func (sf sortableFilers) toFilters() []Filter {
	var filters []Filter

	for _, f := range sf {
		filters = append(filters, f.filter)
	}

	return filters
}

func newSortableFilters(filters []Filter) sortableFilers {
	var result sortableFilers

	for i, f := range filters {
		result = append(result, sortableFilter{
			filter: f,
			index:  i,
		})
	}

	return result
}
