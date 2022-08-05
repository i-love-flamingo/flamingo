package flamingo

import (
	"fmt"
	"testing"

	"flamingo.me/flamingo/v3/framework/config"
	"github.com/stretchr/testify/assert"
)

type testData struct {
	redisURL              string
	redisHost             string
	redisPassword         string
	redisDatabase         int
	expectedRedisHost     string
	expectedRedisPassword string
	expectedRedisDatabase int
}

func TestModule_Configure(t *testing.T) {
	t.Run("empty additional configuration", func(t *testing.T) {
		if err := config.TryModules(nil, new(SessionModule)); err != nil {
			t.Error(err)
		}
	})
	t.Run("invalid redis timeout should lead to error", func(t *testing.T) {
		err := config.TryModules(config.Map{"flamingo.session.redis.timeout": "foo"}, new(SessionModule))
		assert.Error(t, err)
	})
}

func TestGetRedisConnectionInformation(t *testing.T) {
	redisURLHost := "redis-url-host:68043"
	redisURLUser := "redis-url-user"
	redisURLPassword := "redis-url-pw"
	redisDatabase := 2
	redisURL := fmt.Sprintf("redis://%s:%s@%s/%d", redisURLUser, redisURLPassword, redisURLHost, redisDatabase)
	redisURLWithDatabaseInQuery := fmt.Sprintf("redis://%s:%s@%s?db=%d", redisURLUser, redisURLPassword, redisURLHost, redisDatabase)
	redisURLWithoutUser := fmt.Sprintf("redis://:%s@%s/%d", redisURLPassword, redisURLHost, redisDatabase)
	redisHost := "redis-host"
	redisPassword := "pw1234"

	testSet := map[string]testData{
		"url only without user": {
			redisURL:              redisURLWithoutUser,
			expectedRedisHost:     redisURLHost,
			expectedRedisPassword: redisURLPassword,
			expectedRedisDatabase: redisDatabase,
		},
		"url only": {
			redisURL:              redisURL,
			expectedRedisHost:     redisURLHost,
			expectedRedisPassword: redisURLPassword,
			expectedRedisDatabase: redisDatabase,
		},
		"url and host": {
			redisURL:              redisURL,
			redisHost:             redisHost,
			expectedRedisHost:     redisURLHost,
			expectedRedisPassword: redisURLPassword,
			expectedRedisDatabase: redisDatabase,
		},
		"url and host and password": {
			redisURL:              redisURL,
			redisHost:             redisHost,
			redisPassword:         redisPassword,
			expectedRedisHost:     redisURLHost,
			expectedRedisPassword: redisURLPassword,
			expectedRedisDatabase: redisDatabase,
		},
		"url and host and password and database": {
			redisURL:              redisURL,
			redisHost:             redisHost,
			redisPassword:         redisPassword,
			redisDatabase:         redisDatabase,
			expectedRedisHost:     redisURLHost,
			expectedRedisPassword: redisURLPassword,
			expectedRedisDatabase: redisDatabase,
		},
		"host only": {
			redisHost:         redisHost,
			expectedRedisHost: redisHost,
		},
		"password only": {
			redisPassword:         redisPassword,
			expectedRedisPassword: redisPassword,
		},
		"database only": {
			redisDatabase:         redisDatabase,
			expectedRedisDatabase: redisDatabase,
		},
		"host and password": {
			redisHost:             redisHost,
			redisPassword:         redisPassword,
			redisDatabase:         redisDatabase,
			expectedRedisHost:     redisHost,
			expectedRedisPassword: redisPassword,
			expectedRedisDatabase: redisDatabase,
		},
		"database from query": {
			redisURL:              redisURLWithDatabaseInQuery,
			expectedRedisHost:     redisURLHost,
			expectedRedisPassword: redisURLPassword,
			expectedRedisDatabase: redisDatabase,
		},
	}

	t.Run("Get Redis connection information", func(t *testing.T) {
		for name, data := range testSet {
			t.Run(name, func(t *testing.T) {
				actualHost, actualPassword, actualDatabase := getRedisConnectionInformation(data.redisURL, data.redisHost, data.redisPassword, data.redisDatabase)
				assert.Equal(t, data.expectedRedisHost, actualHost)
				assert.Equal(t, data.expectedRedisPassword, actualPassword)
				assert.Equal(t, data.expectedRedisDatabase, actualDatabase)
			})
		}
	})
}
