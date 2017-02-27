package logger

// Logger logs logmessages...
// TODO
type Logger interface {
	Debug(fmt string, a ...interface{})
}
