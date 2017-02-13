package backend

import "net/url"

// ProductServicerCreater psc
type ProductServicerCreater func(string) ProductServicer

var productServices map[string]ProductServicerCreater

func init() {
	productServices = make(map[string]ProductServicerCreater)
}

// RegisterProductService register a backend
func RegisterProductService(name string, psc ProductServicerCreater) {
	productServices[name] = psc
}

// CreateProductService factory
func CreateProductService(dsn string) ProductServicer {
	cfg, err := url.Parse(dsn)
	if err != nil {
		panic(err)
	}
	if productServices[cfg.Scheme] == nil {
		panic("unknown product service " + cfg.Scheme)
	}
	return productServices[cfg.Scheme](dsn)
}

// ProductServicer defines the goom default product backend
type ProductServicer interface {
	//WithProfiler(*profiler.Profile) ProductServicer
	Get(string) Producter
	GetBySkuList([]string) []Producter
}

// Producter default behaviour
type Producter interface {
	Sku() string
	Name() string
	Description() string
	Price() float64
}
