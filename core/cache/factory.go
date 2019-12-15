package cache

import (
	"errors"
	"flamingo.me/flamingo/v3/framework/config"
)

type (
	//HTTPFrontendFactory that can be used to build caches
	HTTPFrontendFactory struct {
		provider               HTTPFrontendProvider
		redisBackendFactory    *RedisBackendFactory
		inMemoryBackendFactory *InMemoryBackendFactory
		twoLevelBackendFactory *TwoLevelBackendFactory
		cacheConfig            config.Map
	}

	//FactoryConfig typed configuration used to build Caches by the factory
	FactoryConfig map[string]BackendConfig

	//BackendConfig typed configuration used to build BackendCaches by the factory
	BackendConfig struct {
		BackendType           string
		InMemoryBackend       *InMemoryBackendConfig
		RedisBackend          *RedisBackendConfig
		TwoLevelBackendFirst  *BackendConfig
		TwoLevelBackendSecond *BackendConfig
	}

	//HTTPFrontendProvider - Dingo Provider func
	HTTPFrontendProvider func() *HTTPFrontend
)

//Inject for dependencies
func (f *HTTPFrontendFactory) Inject(
	provider HTTPFrontendProvider,
	redisBackendFactory *RedisBackendFactory,
	inMemoryBackendFactory *InMemoryBackendFactory,
	twoLevelBackendFactory *TwoLevelBackendFactory,
	config *struct {
		CacheConfig config.Map `inject:"config:core.cache.httpFrontendFactory,optional"`
	},
) {
	f.provider = provider
	f.inMemoryBackendFactory = inMemoryBackendFactory
	f.redisBackendFactory = redisBackendFactory
	f.twoLevelBackendFactory = twoLevelBackendFactory
	if config != nil {
		f.cacheConfig = config.CacheConfig
	}
}

//BuildConfiguredCache with the given name by injected configuration or error
func (f *HTTPFrontendFactory) BuildConfiguredCache(cacheName string) (*HTTPFrontend, error) {
	var cacheConfig FactoryConfig
	err := f.cacheConfig.MapInto(&cacheConfig)
	if err != nil {
		return nil, err
	}
	if v, found := cacheConfig[cacheName]; found {
		backend, err := f.BuildBackend(v, cacheName)
		if err != nil {
			return nil, err
		}
		return f.BuildWithBackend(backend), nil
	}
	return nil, errors.New("Cannot find config for " + cacheName)
}

//BuildWithBackend - returns new HTTPFrontend cache with given backend
func (f *HTTPFrontendFactory) BuildWithBackend(backend Backend) *HTTPFrontend {
	frontend := f.provider()
	frontend.backend = backend
	return frontend
}

//BuildBackend by given BackendConfig and frontendName
func (f *HTTPFrontendFactory) BuildBackend(bc BackendConfig, frontendName string) (Backend, error) {
	switch bc.BackendType {
	case "redis":
		if bc.RedisBackend == nil {
			return nil, errors.New("No RedisBackend config provided")
		}
		return f.RedisBackend(*bc.RedisBackend, frontendName), nil
	case "inmemory":
		if bc.InMemoryBackend == nil {
			return nil, errors.New("No InMemoryBackend config provided")
		}
		return f.MemoryBackend(*bc.InMemoryBackend, frontendName), nil
	case "twolevel":
		if bc.TwoLevelBackendFirst == nil || bc.TwoLevelBackendSecond == nil {
			return nil, errors.New("No TwoLevelBackendFirst config provided")
		}
		first, err := f.BuildBackend(*bc.TwoLevelBackendFirst, frontendName)
		if err != nil {
			return nil, err
		}
		second, err := f.BuildBackend(*bc.TwoLevelBackendSecond, frontendName)
		if err != nil {
			return nil, err
		}
		return f.TwoLevel(TwoLevelBackendConfig{first, second}), nil
	}
	return nil, errors.New("Unknown Backend Type")
}

//MemoryBackend - returns new MemoryBackend
func (f *HTTPFrontendFactory) MemoryBackend(config InMemoryBackendConfig, frontendName string) Backend {
	return f.inMemoryBackendFactory.SetConfig(config).SetFrontendName(frontendName).Build()
}

//RedisBackend returns new RedisBackend Backend
func (f *HTTPFrontendFactory) RedisBackend(config RedisBackendConfig, frontendName string) Backend {
	return f.redisBackendFactory.SetPoolByConfig(config).SetFrontendName(frontendName).Build()
}

//TwoLevel - returns new TwoLevel Backend Cache
func (f *HTTPFrontendFactory) TwoLevel(config TwoLevelBackendConfig) Backend {
	return f.twoLevelBackendFactory.SetConfig(config).Build()
}
