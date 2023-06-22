package zap_test

import (
	"fmt"
	"testing"

	"flamingo.me/dingo"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

func BenchmarkLogger(b *testing.B) {
	fields := make(map[flamingo.LogKey]any, 100)

	for i := 0; i < 100; i++ {
		fields[flamingo.LogKey(fmt.Sprintf("field-%01d", i))] = fmt.Sprintf("value-%01d", i)
	}

	zapModule := new(zap.Module)
	area := config.NewArea("test", []dingo.Module{zapModule})
	err := config.Load(area, "",
		// Warn Level so that the benchmark does not print to output
		config.AdditionalConfig([]string{"core.zap.loglevel: Warn"}))
	require.NoError(b, err)

	injector, err := area.GetInitializedInjector()
	require.NoError(b, err)

	l, err := injector.GetInstance(new(flamingo.Logger))
	require.NoError(b, err)

	require.IsTypef(b, new(zap.Logger), l, "logger must be a *zap.Logger")
	logger := l.(flamingo.Logger)

	b.Run("withField-1-log", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fieldedLogger := logger
			for k, v := range fields {
				fieldedLogger = fieldedLogger.WithField(k, v)
			}
			fieldedLogger.Info("Test Log")
		}
	})

	b.Run("withFields-1-log", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fieldedLogger := logger
			fieldedLogger = fieldedLogger.WithFields(fields)
			fieldedLogger.Info("Test Log")
		}
	})
	b.Run("withField-each-log", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fieldedLogger := logger
			for k, v := range fields {
				fieldedLogger = fieldedLogger.WithField(k, v)
				fieldedLogger.Info("Test Log")
			}
		}
	})

	b.Run("withFields-many-logs", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fieldedLogger := logger
			fieldedLogger = fieldedLogger.WithFields(fields)
			for i := 0; i < len(fields); i++ {
				fieldedLogger.Info("Test Log")
			}
		}
	})
}
