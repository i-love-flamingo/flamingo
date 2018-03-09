package logrus

import (
	"fmt"

	"go.aoe.com/flamingo/framework/flamingo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	// Logger is a Wrapper for the zap logger
	Logger struct {
		*zap.Logger
	}
)

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.Infof(format, args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.Warnf(format, args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatal(fmt.Sprintf(format, args...))
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	l.Logger.Panic(fmt.Sprintf(format, args...))
}

func (l *Logger) Debug(args ...interface{}) {
	l.Logger.Debug(fmt.Sprint(args...))
}

func (l *Logger) Info(args ...interface{}) {
	l.Logger.Info(fmt.Sprint(args...))
}

func (l *Logger) Print(args ...interface{}) {
	l.Logger.Info(fmt.Sprint(args...))
}

func (l *Logger) Warn(args ...interface{}) {
	l.Logger.Warn(fmt.Sprint(args...))
}

func (l *Logger) Warning(args ...interface{}) {
	l.Warn(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Logger.Error(fmt.Sprint(args...))
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(fmt.Sprint(args...))
}

func (l *Logger) Panic(args ...interface{}) {
	l.Logger.Panic(fmt.Sprint(args...))
}

func (l *Logger) Debugln(args ...interface{}) {
	l.Logger.Debug(fmt.Sprint(args...))
}

func (l *Logger) Infoln(args ...interface{}) {
	l.Logger.Info(fmt.Sprint(args...))
}

func (l *Logger) Println(args ...interface{}) {
	l.Logger.Info(fmt.Sprint(args...))
}

func (l *Logger) Warnln(args ...interface{}) {
	l.Logger.Warn(fmt.Sprint(args...))
}

func (l *Logger) Warningln(args ...interface{}) {
	l.Warning(args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.Logger.Error(fmt.Sprint(args...))
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.Logger.Fatal(fmt.Sprint(args...))
}

func (l *Logger) Panicln(args ...interface{}) {
	l.Logger.Panic(fmt.Sprint(args...))
}

func (l *Logger) WithField(key string, value interface{}) flamingo.Logger {
	return &Logger{
		Logger: l.Logger.With(
			zapcore.Field{
				Key:       key,
				Type:      zapcore.StringType,
				Integer:   0,
				String:    "",
				Interface: value,
			},
		),
	}
}

func (l *Logger) WithFields(fields map[string]interface{}) flamingo.Logger {
	zapFields := make([]zapcore.Field, len(fields))
	for key, value := range fields {
		zapFields = append(
			zapFields,
			zapcore.Field{
				Key:       key,
				Type:      zapcore.StringType,
				Integer:   0,
				String:    "",
				Interface: value,
			},
		)
	}

	return &Logger{
		Logger: l.Logger.With(zapFields...),
	}
}

func (l *Logger) WithError(err error) flamingo.Logger {
	return l.WithField("error", err)
}
