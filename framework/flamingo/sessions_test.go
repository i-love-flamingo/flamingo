package flamingo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	redisURL              string
	redisHost             string
	redisPassword         string
	redisDatabase         string
	expectedRedisHost     string
	expectedRedisPassword string
	expectedRedisDatabase string
}

func TestGetRedisConnectionInformation(t *testing.T) {
	redisURLHost := "redis-url-host:68043"
	redisURLUser := "redis-url-user"
	redisURLPassword := "redis-url-pw"
	redisDatabase := "2"
	redisURL := fmt.Sprintf("redis://%s:%s@%s/%s", redisURLUser, redisURLPassword, redisURLHost, redisDatabase)
	redisURLWithDatabaseInQuery := fmt.Sprintf("redis://%s:%s@%s?db=%s", redisURLUser, redisURLPassword, redisURLHost, redisDatabase)
	redisURLWithoutUser := fmt.Sprintf("redis://:%s@%s/%s", redisURLPassword, redisURLHost, redisDatabase)
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
