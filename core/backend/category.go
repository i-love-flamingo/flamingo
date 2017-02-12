package backend

import (
	"net/url"
)

// CategoryServicerCreater psc
type CategoryServicerCreater func(string) CategoryServicer

var categoryServices map[string]CategoryServicerCreater

func init() {
	categoryServices = make(map[string]CategoryServicerCreater)
}

// RegisterCategoryService register a backend
func RegisterCategoryService(name string, psc CategoryServicerCreater) {
	categoryServices[name] = psc
}

// CreateCategoryService factory
func CreateCategoryService(dsn string) CategoryServicer {
	cfg, err := url.Parse(dsn)
	if err != nil {
		panic(err)
	}

	if categoryServices[cfg.Scheme] == nil {
		panic("unknown category service " + cfg.Scheme)
	}
	return categoryServices[cfg.Scheme](dsn)
}

// CategoryServicer defines the goom default Category backend
type CategoryServicer interface {
	//WithProfiler(*profiler.Profile) CategoryServicer

	Root() Categoryer
	Get(string) Categoryer
}

// Categoryer default behaviour
type Categoryer interface {
	Name() string
	Key() string
	Childs() []Categoryer
	ProductSkuList() []string
}
