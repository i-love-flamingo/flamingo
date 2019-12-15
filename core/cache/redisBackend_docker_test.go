// +build docker

package cache_test

import (
	"flamingo.me/flamingo/v3/framework/flamingo"
	"fmt"
	"log"
	"os"
	"testing"

	"flamingo.me/flamingo/v3/core/cache"
	"github.com/gomodule/redigo/redis"
	"github.com/ory/dockertest"
)

var (
	dockerTestPool     *dockertest.Pool
	dockerTestResource *dockertest.Resource
)

func redisConnector(network, address, password string, db int) (redis.Conn, error) {
	c, err := redis.Dial(network, address)
	if err != nil {
		return nil, err
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}
	if db != 0 {
		if _, err := c.Do("SELECT", db); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

// TestMain set
func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown() // comment out, if you want to keep the docker-instance running for debugging
	os.Exit(code)
}

// setup an redis docker-container for integration tests
func setup(m *testing.M) {

	var err error
	dockerTestPool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	dockerTestResource, err = dockerTestPool.Run("redis", "4-alpine", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// test connection while setup - no need to run other tests, if connection setup fails in setup
	connection, err := redis.Dial("tcp", fmt.Sprintf("%v:%v", "127.0.0.1", dockerTestResource.GetPort("6379/tcp")))
	if err != nil {
		log.Fatalf("Could not connect to redis-docker: %s", err)
	}
	err = redis.Conn.Close(connection)
	if err != nil {
		log.Fatalf("Could not close redis-docker: %s", err)
	}
}

// teardown the redis docker-container
func teardown() {
	err := dockerTestPool.Purge(dockerTestResource)
	if err != nil {
		log.Fatalf("Error purging docker resources: %s", err)
	}
}

func Test_RunDefaultBackendTestCase_RedisBackend(t *testing.T) {

	config := cache.RedisBackendConfig{
		MaxIdle:     8,
		IdleTimeOut: 30,
		Host:        "127.0.0.1",
		Port:        dockerTestResource.GetPort("6379/tcp"),
	}
	factory := cache.RedisBackendFactory{}
	backend := factory.Inject(flamingo.NullLogger{}).SetPoolByConfig(config).SetFrontendName("testfrontend").Build()
	testcase := cache.NewBackendTestCase(t, backend, false)
	testcase.RunTests()
}
