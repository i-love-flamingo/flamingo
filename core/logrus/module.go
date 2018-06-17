package logrus

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// Module for logrus logging
	Module struct {
		Area     string `inject:"config:area"`
		JSON     bool   `inject:"config:logrus.json,optional"`
		LogLevel string `inject:"config:logrus.loglevel,optional"`
	}

	// LogrusEntry is a Wrapper for the logrus.Entry logger fulfilling the flamingo.Logger interface
	LogrusEntry struct {
		*logrus.Entry
	}
	// LogrusLogger is a Wrapper for the logrus.Logger fulfilling the flamingo.Logger interface
	LogrusLogger struct {
		*logrus.Logger
	}

	ContextHook struct {
		area string
	}
)

// Levels returns all available logrus log levels
func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

var lock = new(sync.Mutex)

func (hook ContextHook) Fire(entry *logrus.Entry) error {
	lock.Lock()
	defer lock.Unlock()

	pc := make([]uintptr, 3, 3)
	cnt := runtime.Callers(6, pc)

	entry.Data["area"] = hook.area

	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !strings.Contains(name, "github.com/sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			entry.Data["source"] = fmt.Sprintf("File: %v Func: %v  Line: %v", lastPathAndFileName(file), path.Base(name), line)
			entry.Data["fileName"] = path.Base(file)
			break
		}
	}
	return nil
}

//lastPathAndFileName returns the filename and the last two folders of the Path
func lastPathAndFileName(completePath string) string {
	path, fileName := path.Split(completePath)
	dirNames := strings.Split(strings.Trim(path, "/"), "/")
	if len(dirNames) > 2 {
		return dirNames[len(dirNames)-2] + "/" + dirNames[len(dirNames)-1] + "/" + fileName
	}
	return fileName
}

// Configure the logrus logger as flamingo.Logger (in JSON mode kibana compatible)
func (m *Module) Configure(injector *dingo.Injector) {
	var l *logrus.Logger
	if m.JSON {
		l = &logrus.Logger{
			Out: os.Stderr,
			Formatter: &logrus.JSONFormatter{
				TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyTime:  "@timestamp",
					logrus.FieldKeyLevel: "level",
					logrus.FieldKeyMsg:   "message",
				},
			},
			Hooks: make(logrus.LevelHooks),
		}
	} else {
		l = &logrus.Logger{
			Out: os.Stderr,
			Formatter: &logrus.TextFormatter{
				TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
			},
			Hooks: make(logrus.LevelHooks),
		}
	}
	l.Level = logrus.ErrorLevel
	if m.LogLevel == "Info" {
		l.Level = logrus.InfoLevel
	} else if m.LogLevel == "Debug" {
		l.Level = logrus.DebugLevel
	}
	l.Hooks.Add(ContextHook{area: m.Area})
	injector.Bind((*flamingo.Logger)(nil)).ToInstance(&LogrusLogger{l})
}

// WithContext returns a logger with fields filled from the context
// businessId:    From Header X-Business-ID
// client_ip:     From Header X-Forwarded-For or request if header is empty
// correlationId: The ID of the context
// method:        HTTP verb from request
// path:          URL path from request
// referer:       referer from request
// request:       received payload from request
func (e *LogrusEntry) WithContext(ctx web.Context) flamingo.Logger {
	return appendContext(e, ctx)
}

// Flush does nothing because logrus is not buffered
func (e *LogrusEntry) Flush() {}

// WithField adds a single field to the Entry.
func (e *LogrusEntry) WithField(key flamingo.LogKey, value interface{}) flamingo.Logger {
	return &LogrusEntry{Entry: e.Entry.WithField(string(key), value)}
}

// WithFields adds a map of fields to the Entry.
func (e *LogrusEntry) WithFields(fields map[flamingo.LogKey]interface{}) flamingo.Logger {
	f := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		f[string(k)] = v
	}
	return &LogrusEntry{Entry: e.Entry.WithFields(f)}
}

// WithContext returns a logger with fields filled from the context
// businessId:    From Header X-Business-ID
// client_ip:     From Header X-Forwarded-For or request if header is empty
// correlationId: The ID of the context
// method:        HTTP verb from request
// path:          URL path from request
// referer:       referer from request
// request:       received payload from request
func (e *LogrusLogger) WithContext(ctx web.Context) flamingo.Logger {
	return appendContext(e, ctx)
}

// Flush does nothing because logrus is not buffered
func (e *LogrusLogger) Flush() {}

// WithField adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (e *LogrusLogger) WithField(key flamingo.LogKey, value interface{}) flamingo.Logger {
	return &LogrusEntry{Entry: e.Logger.WithField(string(key), value)}
}

// WithFields adds a struct of fields to the log entry. All it does is call `WithField` for each `Field`.
func (e *LogrusLogger) WithFields(fields map[flamingo.LogKey]interface{}) flamingo.Logger {
	f := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		f[string(k)] = v
	}
	return &LogrusEntry{Entry: e.Logger.WithFields(f)}
}

func appendContext(logger flamingo.Logger, ctx web.Context) flamingo.Logger {
	request := ctx.Request()
	clientIP := request.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = request.RemoteAddr
	}
	body, _ := ioutil.ReadAll(request.Body)

	return logger.WithFields(
		map[flamingo.LogKey]interface{}{
			flamingo.LogKeyBusinessID:    request.Header.Get("X-Business-ID"),
			flamingo.LogKeyClientIP:      clientIP,
			flamingo.LogKeyCorrelationID: ctx.ID(),
			flamingo.LogKeyMethod:        request.Method,
			flamingo.LogKeyPath:          request.URL.Path,
			flamingo.LogKeyReferer:       request.Referer(),
			flamingo.LogKeyRequest:       string(body),
		},
	)
}
