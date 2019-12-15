package cache

import (
	config2 "flamingo.me/flamingo/v3/framework/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHTTPFrontendFactory_ConfigUnmarshalling(t *testing.T) {
	testconfig := config2.Map{
		"one": config2.Map{
			"backendType": "inmemory",
			"inMemoryBackend": config2.Map{
				"size": 100.0,
			},
		},
		"two": config2.Map{
			"backendType": "inmemory",
			"inMemoryBackend": config2.Map{
				"size": 100.0,
			},
		},
	}

	var typedCacheConfig FactoryConfig
	testconfig.MapInto(&typedCacheConfig)

	assert.Contains(t, typedCacheConfig, "one")
	assert.Contains(t, typedCacheConfig, "two")

	one := typedCacheConfig["one"]
	assert.Equal(t, "inmemory", one.BackendType)
	require.NotNil(t, one.InMemoryBackend)
	assert.Equal(t, one.InMemoryBackend.Size, 100)

}

func TestHTTPFrontendFactory_BuildBackend(t *testing.T) {
	provider := func() *HTTPFrontend {
		return &HTTPFrontend{}
	}
	f := &HTTPFrontendFactory{}
	f.Inject(provider, &RedisBackendFactory{}, &InMemoryBackendFactory{}, &TwoLevelBackendFactory{}, nil)

	t.Run("inmemory", func(t *testing.T) {
		testConfig := BackendConfig{
			BackendType:     "inmemory",
			InMemoryBackend: &InMemoryBackendConfig{Size: 10},
		}
		backend, err := f.BuildBackend(testConfig, "test")
		assert.NoError(t, err)
		assert.IsType(t, &inMemoryBackend{}, backend)
	})

	t.Run("inmemory error", func(t *testing.T) {
		testConfig := BackendConfig{
			BackendType: "inmemory",
		}
		_, err := f.BuildBackend(testConfig, "test")
		assert.Error(t, err)
	})

	t.Run("redis", func(t *testing.T) {
		testConfig := BackendConfig{
			BackendType:  "redis",
			RedisBackend: &RedisBackendConfig{},
		}
		backend, err := f.BuildBackend(testConfig, "test")
		assert.NoError(t, err)
		assert.IsType(t, &redisBackend{}, backend)
	})
}
