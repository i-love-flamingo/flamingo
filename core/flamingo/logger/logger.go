package logger

type Logger interface {
	Debug(fmt string, a ...interface{})
}
