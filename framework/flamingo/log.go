package flamingo

//go:generate mockery -name "Logger"

// Logger defines a standard Flamingo logger interfaces
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})

	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
	WithError(err error) Logger
}

// NullLogger does not log
type NullLogger struct{}

func (NullLogger) Debugf(format string, args ...interface{})         {}
func (NullLogger) Infof(format string, args ...interface{})          {}
func (NullLogger) Printf(format string, args ...interface{})         {}
func (NullLogger) Warnf(format string, args ...interface{})          {}
func (NullLogger) Warningf(format string, args ...interface{})       {}
func (NullLogger) Errorf(format string, args ...interface{})         {}
func (NullLogger) Fatalf(format string, args ...interface{})         {}
func (NullLogger) Panicf(format string, args ...interface{})         {}
func (NullLogger) Debug(args ...interface{})                         {}
func (NullLogger) Info(args ...interface{})                          {}
func (NullLogger) Print(args ...interface{})                         {}
func (NullLogger) Warn(args ...interface{})                          {}
func (NullLogger) Warning(args ...interface{})                       {}
func (NullLogger) Error(args ...interface{})                         {}
func (NullLogger) Fatal(args ...interface{})                         {}
func (NullLogger) Panic(args ...interface{})                         {}
func (NullLogger) Debugln(args ...interface{})                       {}
func (NullLogger) Infoln(args ...interface{})                        {}
func (NullLogger) Println(args ...interface{})                       {}
func (NullLogger) Warnln(args ...interface{})                        {}
func (NullLogger) Warningln(args ...interface{})                     {}
func (NullLogger) Errorln(args ...interface{})                       {}
func (NullLogger) Fatalln(args ...interface{})                       {}
func (NullLogger) Panicln(args ...interface{})                       {}
func (n NullLogger) WithField(key string, value interface{}) Logger  { return n }
func (n NullLogger) WithFields(fields map[string]interface{}) Logger { return n }
func (n NullLogger) WithError(err error) Logger                      { return n }
