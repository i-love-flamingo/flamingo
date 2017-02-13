package mock

import (
	"flamingo/core/backend"
	"strings"
)

// CategoryService cs
type CategoryService struct {
	dsn string
	//profiler *profiler.Profile
}

// NewCategoryService ncs
func NewCategoryService(dsn string) backend.CategoryServicer {
	return CategoryService{dsn: dsn}
}

// WithProfiler make it profileable
/*func (cs CategoryService) WithProfiler(p *profiler.Profile) backend.CategoryServicer {
	ncs := cs
	ncs.profiler = p
	return ncs
}*/

// Root r
func (cs CategoryService) Root() backend.Categoryer {
	/*if cs.profiler != nil {
		cs.profiler.Start("API.mock", "Get Category Root")
		defer cs.profiler.End()
	}*/

	return Category{
		name: "Root",
		childs: []Category{
			Category{name: "Shirts", childs: []Category{Category{name: "T-Shirts"}, Category{name: "Pullover"}}},
			Category{name: "Bags"},
			Category{name: "Jeans"},
			Category{name: "Backpacks"},
		},
	}
}

// Get category
func (cs CategoryService) Get(name string) backend.Categoryer {
	/*if cs.profiler != nil {
		cs.profiler.Start("API.mock", "Get Category "+name)
		defer cs.profiler.End()
	}*/
	return cs.Root()
}

// Category c
type Category struct {
	name   string
	key    string
	childs []Category
}

// Name n
func (c Category) Name() string {
	return c.name
}

// Key k
func (c Category) Key() string {
	return strings.ToLower(c.name)
}

// Childs c
func (c Category) Childs() []backend.Categoryer {
	childs := make([]backend.Categoryer, len(c.childs))
	for i, v := range c.childs {
		childs[i] = v
	}
	return childs
}

// ProductSkuList psl
func (c Category) ProductSkuList() []string {
	return []string{
		"prod-1",
		"prod-2",
		"prod-3",
		"prod-4",
		"prod-5",
		"prod-6",
		"prod-7",
		"prod-8",
		"prod-9",
		"prod-10",
		"prod-11",
		"prod-12",
	}
}
