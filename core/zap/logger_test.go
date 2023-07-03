package zap_test

import (
	"testing"

	"flamingo.me/dingo"
	"github.com/stretchr/testify/require"

	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

func BenchmarkLogger(b *testing.B) {
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
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				fieldedLogger := logger.
					WithField("key1", "value1").
					WithField("key2", "value2").
					WithField("key3", "value3").
					WithField("key4", "value4").
					WithField("key5", "value5").
					WithField("key6", "value6")
				fieldedLogger.Info("Test Log")
			}
		})
	})
	b.Run("withFields-1-log", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				fieldedLogger := logger.
					WithFields(
						map[flamingo.LogKey]interface{}{
							"key1": "value1",
							"key2": "value2",
							"key3": "value3",
							"key4": "value4",
							"key5": "value5",
							"key6": "value6",
						})
				fieldedLogger.Info("Test Log")
			}
		})
	})
	b.Run("withField-each-log", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				fieldedLogger := logger.WithField("key1", "value1")
				fieldedLogger.Info("Test Log")

				fieldedLogger = fieldedLogger.WithField("key2", "value2")
				fieldedLogger.Info("Test Log")

				fieldedLogger = fieldedLogger.WithField("key3", "value3")
				fieldedLogger.Info("Test Log")

				fieldedLogger = fieldedLogger.WithField("key4", "value4")
				fieldedLogger.Info("Test Log")

				fieldedLogger = fieldedLogger.WithField("key5", "value5")
				fieldedLogger.Info("Test Log")

				fieldedLogger = fieldedLogger.WithField("key6", "value6")
				fieldedLogger.Info("Test Log")
			}
		})
	})

	b.Run("withFields-many-logs", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				fieldedLogger := logger.
					WithFields(
						map[flamingo.LogKey]interface{}{
							"key1": "value1",
							"key2": "value2",
							"key3": "value3",
							"key4": "value4",
							"key5": "value5",
							"key6": "value6",
						})

				fieldedLogger.Info("Test Log")
				fieldedLogger.Info("Test Log")
				fieldedLogger.Info("Test Log")
				fieldedLogger.Info("Test Log")
				fieldedLogger.Info("Test Log")
				fieldedLogger.Info("Test Log")
				fieldedLogger.Info("Test Log")
				fieldedLogger.Info("Test Log")
				fieldedLogger.Info("Test Log")
				fieldedLogger.Info("Test Log")
			}
		})
	})
}
