package zap_test

import (
	"context"
	"net/http"
	"testing"

	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/trace"
	uberZap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogger_WithContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		argAndWant func(t *testing.T) (context.Context, map[flamingo.LogKey]string)
	}{
		{
			name: "trace and span id fields should be added",
			argAndWant: func(t *testing.T) (context.Context, map[flamingo.LogKey]string) {
				t.Helper()

				ctx, span := trace.StartSpan(context.Background(), "test")
				t.Cleanup(span.End)

				return ctx, map[flamingo.LogKey]string{
					flamingo.LogKeyTraceID: span.SpanContext().TraceID.String(),
					flamingo.LogKeySpanID:  span.SpanContext().SpanID.String(),
				}
			},
		},
		{
			name: "request method and path fields should be added",
			argAndWant: func(t *testing.T) (context.Context, map[flamingo.LogKey]string) {
				t.Helper()

				req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://example.com/deep/path", nil)
				require.NoError(t, err)

				ctx := web.ContextWithRequest(context.Background(), web.CreateRequest(req, web.EmptySession()))

				return ctx, map[flamingo.LogKey]string{
					flamingo.LogKeyMethod: http.MethodGet,
					flamingo.LogKeyPath:   "/deep/path",
				}
			},
		},
		{
			name: "sessoion id hash field should be added",
			argAndWant: func(t *testing.T) (context.Context, map[flamingo.LogKey]string) {
				t.Helper()

				session := web.EmptySession()
				ctx := web.ContextWithSession(context.Background(), session)

				return ctx, map[flamingo.LogKey]string{
					flamingo.LogKeySession: session.IDHash(),
				}
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			core, observedLogs := observer.New(zapcore.InfoLevel)
			l := zap.NewLogger(uberZap.New(core, uberZap.WithFatalHook(zapcore.WriteThenNoop)), zap.WithLogSession(true))

			ctx, expectedFields := tt.argAndWant(t)

			ctxLogger := l.WithContext(ctx)
			ctxLogger.Info()

			require.Equal(t, observedLogs.Len(), 1)

			foundFields := make(map[string]struct{})

			for _, field := range observedLogs.All()[0].Context {
				if expected, ok := expectedFields[flamingo.LogKey(field.Key)]; ok {
					foundFields[field.Key] = struct{}{}

					assert.Equal(t, expected, field.String, "field %q not as value", field.Key)
				}
			}

			for key := range expectedFields {
				_, found := foundFields[string(key)]
				assert.True(t, found, "expected log field %q not found in entry", key)
			}
		})
	}
}

func TestLogger_WithField(t *testing.T) {
	t.Parallel()

	type fields struct {
		fieldMap map[string]string
	}

	type args struct {
		key   flamingo.LogKey
		value string
	}

	type want struct {
		key   string
		value string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		want      want
		wantFound bool
	}{
		{
			name: "add given field as log key",
			args: args{
				key:   flamingo.LogKeyCategory,
				value: "test",
			},
			want: want{
				key:   string(flamingo.LogKeyCategory),
				value: "test",
			},
			wantFound: true,
		},
		{
			name: "add given field as log key using configured alias",
			fields: fields{
				fieldMap: map[string]string{
					string(flamingo.LogKeyCategory): "alias",
				},
			},
			args: args{
				key:   flamingo.LogKeyCategory,
				value: "test",
			},
			want: want{
				key:   "alias",
				value: "test",
			},
			wantFound: true,
		},
		{
			name: "ignore given field because of configured alias `-`",
			fields: fields{
				fieldMap: map[string]string{
					string(flamingo.LogKeyCategory): "-",
				},
			},
			args: args{
				key:   flamingo.LogKeyCategory,
				value: "test",
			},
			wantFound: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			core, observedLogs := observer.New(zapcore.InfoLevel)
			l := zap.NewLogger(uberZap.New(
				core,
				uberZap.WithFatalHook(zapcore.WriteThenNoop)),
				zap.WithFieldMap(tt.fields.fieldMap),
			).WithField(tt.args.key, tt.args.value)

			l.Info()
			require.Equal(t, observedLogs.Len(), 1)

			found := false

			for _, field := range observedLogs.All()[0].Context {
				if field.Key == tt.want.key {
					found = true

					assert.Equal(t, tt.want.value, field.String)
				}
			}

			assert.Equal(t, tt.wantFound, found, "key %q found: %t but wanted: %t", tt.args.key, found, tt.wantFound)
		})
	}
}

func TestLogger_WithFields(t *testing.T) {
	t.Parallel()

	type fields struct {
		fieldMap map[string]string
	}

	type args struct {
		fields map[flamingo.LogKey]any
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]string
	}{
		{
			name: "add given fields as log keys",
			args: args{
				fields: map[flamingo.LogKey]any{
					flamingo.LogKeyCategory:    "test",
					flamingo.LogKeySubCategory: "sub-test",
				},
			},
			want: map[string]string{
				string(flamingo.LogKeyCategory):    "test",
				string(flamingo.LogKeySubCategory): "sub-test",
			},
		},
		{
			name: "add given field as log key using configured alias",
			fields: fields{
				fieldMap: map[string]string{
					string(flamingo.LogKeyCategory):    "alias",
					string(flamingo.LogKeySubCategory): "alias2",
				},
			},
			args: args{
				fields: map[flamingo.LogKey]any{
					flamingo.LogKeyCategory:    "test",
					flamingo.LogKeySubCategory: "sub-test",
				},
			},
			want: map[string]string{
				"alias":  "test",
				"alias2": "sub-test",
			},
		},
		{
			name: "ignore given field because of configured alias `-`",
			fields: fields{
				fieldMap: map[string]string{
					string(flamingo.LogKeyCategory): "-",
				},
			},
			args: args{
				fields: map[flamingo.LogKey]any{
					flamingo.LogKeyCategory:    "test",
					flamingo.LogKeySubCategory: "sub-test",
				},
			},
			want: map[string]string{
				string(flamingo.LogKeySubCategory): "sub-test",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			core, observedLogs := observer.New(zapcore.InfoLevel)
			l := zap.NewLogger(uberZap.New(core, uberZap.WithFatalHook(zapcore.WriteThenNoop)),
				zap.WithFieldMap(tt.fields.fieldMap),
			).WithFields(tt.args.fields)

			l.Info()
			require.Equal(t, observedLogs.Len(), 1)

			foundFields := make(map[string]struct{})

			for _, field := range observedLogs.All()[0].Context {
				if expected, ok := tt.want[field.Key]; ok {
					foundFields[field.Key] = struct{}{}

					assert.Equal(t, expected, field.String, "field %q not as value", field.Key)
				}
			}

			for key := range tt.want {
				_, found := foundFields[key]
				assert.True(t, found, "expected log field %q not found in entry", key)
			}
		})
	}
}

func TestLogger_Trace(t *testing.T) {
	t.Parallel()

	core, observedLogs := observer.New(uberZap.DebugLevel - 1) // trace level
	l := zap.NewLogger(uberZap.New(core, uberZap.WithFatalHook(zapcore.WriteThenNoop)))

	l.Trace("test")

	require.Equal(t, observedLogs.Len(), 1)
	assert.Equal(t, "test", observedLogs.All()[0].Message)
}

func TestLogger_Tracef(t *testing.T) {
	t.Parallel()

	core, observedLogs := observer.New(uberZap.DebugLevel - 1) // trace level
	l := zap.NewLogger(uberZap.New(core, uberZap.WithFatalHook(zapcore.WriteThenNoop)))

	l.Tracef("test %s", "logger")

	require.Equal(t, observedLogs.Len(), 1)
	assert.Equal(t, "test logger", observedLogs.All()[0].Message)
}

func TestLogger_Debug(t *testing.T) {
	t.Parallel()

	core, observedLogs := observer.New(zapcore.DebugLevel)
	l := zap.NewLogger(uberZap.New(core, uberZap.WithFatalHook(zapcore.WriteThenNoop)))

	l.Debug("test")

	require.Equal(t, observedLogs.Len(), 1)
	assert.Equal(t, "test", observedLogs.All()[0].Message)
}

func TestLogger_Debugf(t *testing.T) {
	t.Parallel()

	core, observedLogs := observer.New(zapcore.DebugLevel)
	l := zap.NewLogger(uberZap.New(core, uberZap.WithFatalHook(zapcore.WriteThenNoop)))

	l.Debugf("test %s", "logger")

	require.Equal(t, observedLogs.Len(), 1)
	assert.Equal(t, "test logger", observedLogs.All()[0].Message)
}

func TestLogger_Info(t *testing.T) {
	t.Parallel()

	core, observedLogs := observer.New(zapcore.InfoLevel)
	l := zap.NewLogger(uberZap.New(core, uberZap.WithFatalHook(zapcore.WriteThenNoop)))

	l.Info("test")

	require.Equal(t, observedLogs.Len(), 1)
	assert.Equal(t, "test", observedLogs.All()[0].Message)
}

func TestLogger_Warn(t *testing.T) {
	t.Parallel()

	core, observedLogs := observer.New(zapcore.WarnLevel)
	l := zap.NewLogger(uberZap.New(core, uberZap.WithFatalHook(zapcore.WriteThenNoop)))

	l.Warn("test")

	require.Equal(t, observedLogs.Len(), 1)
	assert.Equal(t, "test", observedLogs.All()[0].Message)
}

func TestLogger_Error(t *testing.T) {
	t.Parallel()

	core, observedLogs := observer.New(zapcore.ErrorLevel)
	l := zap.NewLogger(uberZap.New(core, uberZap.WithFatalHook(zapcore.WriteThenNoop)))

	l.Error("test")

	require.Equal(t, observedLogs.Len(), 1)
	assert.Equal(t, "test", observedLogs.All()[0].Message)
}

func TestLogger_Fatal(t *testing.T) {
	t.Parallel()

	core, observedLogs := observer.New(zapcore.FatalLevel)
	l := zap.NewLogger(uberZap.New(core, uberZap.WithFatalHook(zapcore.WriteThenPanic)))

	assert.Panics(t, func() {
		l.Fatal("test")
	})

	require.Equal(t, observedLogs.Len(), 1)
	assert.Equal(t, "test", observedLogs.All()[0].Message)
}
