package flamingo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/framework/config"
)

type testData struct {
	redisURL              string
	redisHost             string
	redisUsername         string
	redisPassword         string
	redisDatabase         int
	expectedRedisHost     string
	expectedRedisUsername string
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

	t.Run("invalid empty username", func(t *testing.T) {
		t.Parallel()

		err := config.TryModules(config.Map{"flamingo.session.redis.username": ""}, new(SessionModule))
		assert.Error(t, err)
	})
}

func TestGetRedisConnectionInformation(t *testing.T) {
	redisURLHost := "redis-url-host:68043"
	redisURLUser := "redis-url-user"
	redisURLPassword := "redis-url-pw"
	redisURLDatabase := 2
	redisURL := fmt.Sprintf("redis://%s:%s@%s/%d", redisURLUser, redisURLPassword, redisURLHost, redisURLDatabase)
	redisHost := "redis-host"
	redisUser := "user1234"
	redisPassword := "pw1234"
	redisDatabase := 4

	type args struct {
		redisURL      string
		redisUsername string
		redisPassword string
		redisDatabase int
		redisHost     string
	}

	tests := []struct {
		args args
		want struct {
			redisUsername   string
			redisPassword   string
			redisHost       string
			redisDatabase   int
			panicsOnDBParse bool
		}
		name string
	}{
		{
			name: "url only without user",
			args: args{
				redisURL: fmt.Sprintf("redis://:%s@%s/%d", redisURLPassword, redisURLHost, redisURLDatabase),
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername: "",
				redisPassword: redisURLPassword,
				redisHost:     redisURLHost,
				redisDatabase: redisURLDatabase,
			},
		},
		{
			name: "url only",
			args: args{
				redisURL: redisURL,
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername: redisURLUser,
				redisPassword: redisURLPassword,
				redisHost:     redisURLHost,
				redisDatabase: redisURLDatabase,
			},
		},
		{
			name: "url with db in query",
			args: args{
				redisURL: fmt.Sprintf("redis://%s:%s@%s?db=%d", redisURLUser, redisURLPassword, redisURLHost, redisURLDatabase),
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername: redisURLUser,
				redisPassword: redisURLPassword,
				redisHost:     redisURLHost,
				redisDatabase: redisURLDatabase,
			},
		},
		{
			name: "url with configured host",
			args: args{
				redisURL:  redisURL,
				redisHost: redisHost,
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername: redisURLUser,
				redisPassword: redisURLPassword,
				redisHost:     redisURLHost,
				redisDatabase: redisURLDatabase,
			},
		},
		{
			name: "url with configured host and username",
			args: args{
				redisURL:      redisURL,
				redisHost:     redisHost,
				redisUsername: redisUser,
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername: redisURLUser,
				redisPassword: redisURLPassword,
				redisHost:     redisURLHost,
				redisDatabase: redisURLDatabase,
			},
		},
		{
			name: "url with configured host, username and password",
			args: args{
				redisURL:      redisURL,
				redisHost:     redisHost,
				redisUsername: redisUser,
				redisPassword: redisPassword,
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername: redisURLUser,
				redisPassword: redisURLPassword,
				redisHost:     redisURLHost,
				redisDatabase: redisURLDatabase,
			},
		},
		{
			name: "url with configured host, username, password and database",
			args: args{
				redisURL:      redisURL,
				redisHost:     redisHost,
				redisUsername: redisUser,
				redisPassword: redisPassword,
				redisDatabase: redisDatabase,
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername: redisURLUser,
				redisPassword: redisURLPassword,
				redisHost:     redisURLHost,
				redisDatabase: redisURLDatabase,
			},
		},
		{
			name: "host only",
			args: args{
				redisHost: redisHost,
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername: "",
				redisPassword: "",
				redisHost:     redisHost,
				redisDatabase: 0,
			},
		},
		{
			name: "username only",
			args: args{
				redisUsername: redisUser,
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername: redisUser,
				redisPassword: "",
				redisHost:     "",
				redisDatabase: 0,
			},
		},
		{
			name: "password only",
			args: args{
				redisPassword: redisPassword,
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername: "",
				redisPassword: redisPassword,
				redisHost:     "",
				redisDatabase: 0,
			},
		},
		{
			name: "database only",
			args: args{
				redisDatabase: redisDatabase,
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername: "",
				redisPassword: "",
				redisHost:     "",
				redisDatabase: redisDatabase,
			},
		},
		{
			name: "empty url",
			args: args{
				redisURL: "",
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername: "",
				redisPassword: "",
				redisHost:     "",
				redisDatabase: 0,
			},
		},
		{
			name: "wrong url",
			args: args{
				redisURL: "1231://redis_url",
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername: "",
				redisPassword: "",
				redisHost:     "",
				redisDatabase: 0,
			},
		},
		{
			name: "broken db in url",
			args: args{
				redisURL: fmt.Sprintf("redis://%s:%s@%s/%s", redisURLUser, redisURLPassword, redisURLHost, "broken"),
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername:   redisURLUser,
				redisPassword:   redisURLPassword,
				redisHost:       redisURLHost,
				redisDatabase:   0,
				panicsOnDBParse: true,
			},
		},
		{
			name: "broken db in query url",
			args: args{
				redisURL: fmt.Sprintf("redis://%s:%s@%s?db=%s", redisURLUser, redisURLPassword, redisURLHost, "broken"),
			},
			want: struct {
				redisUsername   string
				redisPassword   string
				redisHost       string
				redisDatabase   int
				panicsOnDBParse bool
			}{
				redisUsername:   redisURLUser,
				redisPassword:   redisURLPassword,
				redisHost:       redisURLHost,
				redisDatabase:   0,
				panicsOnDBParse: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actualHost := getRedisHost(tt.args.redisURL, tt.args.redisHost)
			actualUsername := getRedisUsername(tt.args.redisURL, tt.args.redisUsername)
			actualPassword := getRedisPassword(tt.args.redisURL, tt.args.redisPassword)

			assert.Equal(t, tt.want.redisHost, actualHost)
			assert.Equal(t, tt.want.redisUsername, actualUsername)
			assert.Equal(t, tt.want.redisPassword, actualPassword)

			if tt.want.panicsOnDBParse {
				assert.Panics(t, func() {
					getRedisDatabase(tt.args.redisURL, tt.args.redisDatabase)
				})
			} else {
				require.NotPanics(t, func() {
					getRedisDatabase(tt.args.redisURL, tt.args.redisDatabase)
				})

				actualDatabase := getRedisDatabase(tt.args.redisURL, tt.args.redisDatabase)

				assert.Equal(t, tt.want.redisDatabase, actualDatabase)
			}
		})
	}
}
