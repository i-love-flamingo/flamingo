package flamingo_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

func createLogger(recorder io.Writer) flamingo.Logger {
	var l = new(flamingo.StdLogger)
	l.Logger = *log.New(recorder, "TEST--", 0)

	return l
}

// In this example we see the normal usecase of a logger
func ExampleLogger() {
	recorder := &strings.Builder{}

	// the flamingo.Logger is normally obtained via injection
	logger := createLogger(recorder)

	// create a new child logger witch additional fields
	logger = logger.WithFields(map[flamingo.LogKey]any{
		flamingo.LogKeyBusinessID: "example",
		flamingo.LogKeyCategory:   "test",
	})

	// we can also add single fields
	logger = logger.WithField(flamingo.LogKeySubCategory, "another field")

	// we can add the current context. It is up to the logger implementation what to do with it
	// the standard flamingo zap.Logger can add the session hash aon tracing information as fields
	// the StdLogger we use here, does nothing with the context
	logger = logger.WithContext(context.Background())

	// log some messages in different levels
	logger.Info("my", "log", "args")
	logger.Warn("my", "log", "args")
	logger.Debug("my", "log", "args")
	// logger.Trace("my", "log", "args") // this is possible with build tag `tracelog`

	fmt.Println(recorder.String())

	// Output:
	// TEST--WithFields map[businessId:example category:test]
	// TEST--WithField sub_category another field
	// TEST--info: mylogargs
	// TEST--warn: mylogargs
	// TEST--debug: mylogargs
}

func ExampleLogFunc() {
	recorder := &strings.Builder{}

	// the flamingo.Logger is normally obtained via injection
	logger := createLogger(recorder)

	// we can create a log function which we use later to log messages with the fields we define here
	logFunc := flamingo.LogFunc(logger, map[flamingo.LogKey]any{
		flamingo.LogKeyBusinessID: "example",
		flamingo.LogKeyCategory:   "test",
	})

	// the wanted log level is defined by passing the corresponding method from the flamingo.Logger interface
	logFunc(flamingo.Logger.Info, "my", "log", "args")
	logFunc(flamingo.Logger.Warn, "my", "log", "args")
	logFunc(flamingo.Logger.Debug, "my", "log", "args")
	// logFunc(flamingo.Logger.Trace, "my", "log", "args") // this is possible with build tag `tracelog`
	// if we do not give the level, it will be Info
	logFunc(nil, "my", "log", "args")

	fmt.Println(recorder.String())

	// Output:
	// TEST--WithFields map[businessId:example category:test]
	// TEST--info: mylogargs
	// TEST--WithFields map[businessId:example category:test]
	// TEST--warn: mylogargs
	// TEST--WithFields map[businessId:example category:test]
	// TEST--debug: mylogargs
	// TEST--WithFields map[businessId:example category:test]
	// TEST--info: mylogargs
}

func ExampleLogFuncWithContext() {
	recorder := &strings.Builder{}

	// the flamingo.Logger is normally obtained via injection
	logger := createLogger(recorder)

	// we can create a log function which we use later to log messages with the fields we define here
	// the log function will expect a context on each call
	logFunc := flamingo.LogFuncWithContext(logger, map[flamingo.LogKey]any{
		flamingo.LogKeyBusinessID: "example",
	})

	ctx := context.Background()

	// the current context is the first argument. It is up to the logger implementation what to do with it
	// the standard flamingo zap.Logger can add the session hash aon tracing information as fields
	// the StdLogger we use here, does nothing with the context
	// the wanted log level is defined by passing the corresponding method from the flamingo.Logger interface
	logFunc(ctx, flamingo.Logger.Info, "my", "log", "args")
	logFunc(ctx, flamingo.Logger.Warn, "my", "log", "args")
	logFunc(ctx, flamingo.Logger.Debug, "my", "log", "args")
	// logFunc(ctx, flamingo.Logger.Trace, "my", "log", "args") // this is possible with build tag `tracelog`
	// if we do not give the level, it will be Info
	logFunc(ctx, nil, "my", "log", "args")

	fmt.Println(recorder.String())

	// Output:
	// TEST--WithFields map[businessId:example]
	// TEST--info: mylogargs
	// TEST--WithFields map[businessId:example]
	// TEST--warn: mylogargs
	// TEST--WithFields map[businessId:example]
	// TEST--debug: mylogargs
	// TEST--WithFields map[businessId:example]
	// TEST--info: mylogargs
}
