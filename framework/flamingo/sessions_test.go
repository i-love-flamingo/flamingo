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
	expectedRedisHost     string
	expectedRedisPassword string
}

func TestGetRedisConnectionInformation(t *testing.T) {
	redisURLHost := "redis-url-host:68043"
	redisURLUser := "redis-url-user"
	redisURLPassword := "redis-url-pw"
	redisURL := fmt.Sprintf("redis://%s:%s@%s/0", redisURLUser, redisURLPassword, redisURLHost)
	redisURLWithoutUser := fmt.Sprintf("redis://:%s@%s/0", redisURLPassword, redisURLHost)
	redisHost := "redis-host"
	redisPassword := "pw1234"

	testSet := map[string]testData{
		"url only without user": {
			redisURL:              redisURLWithoutUser,
			expectedRedisHost:     redisURLHost,
			expectedRedisPassword: redisURLPassword,
		},
		"url only": {
			redisURL:              redisURL,
			expectedRedisHost:     redisURLHost,
			expectedRedisPassword: redisURLPassword,
		},
		"url and host": {
			redisURL:              redisURL,
			redisHost:             redisHost,
			expectedRedisHost:     redisURLHost,
			expectedRedisPassword: redisURLPassword,
		},
		"url and host and password": {
			redisURL:              redisURL,
			redisHost:             redisHost,
			redisPassword:         redisPassword,
			expectedRedisHost:     redisURLHost,
			expectedRedisPassword: redisURLPassword,
		},
		"host only": {
			redisHost:         redisHost,
			expectedRedisHost: redisHost,
		},
		"password only": {
			redisPassword:         redisPassword,
			expectedRedisPassword: redisPassword,
		},
		"host and password": {
			redisHost:             redisHost,
			redisPassword:         redisPassword,
			expectedRedisHost:     redisHost,
			expectedRedisPassword: redisPassword,
		},
	}

	t.Run("Get Redis connection information", func(t *testing.T) {
		for name, data := range testSet {
			t.Run(name, func(t *testing.T) {
				actualHost, actualPassword := getRedisConnectionInformation(data.redisURL, data.redisHost, data.redisPassword)
				assert.Equal(t, data.expectedRedisHost, actualHost)
				assert.Equal(t, data.expectedRedisPassword, actualPassword)
			})
		}
	})
}
