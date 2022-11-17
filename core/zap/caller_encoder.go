package zap

import (
	"os"
	"runtime"
	"strings"

	"go.uber.org/zap/zapcore"
)

func short(file string) string {
	parts := strings.Split(file, string(os.PathSeparator))

	file = ""

	for i, part := range parts {
		if i == len(parts)-1 || len(part) == 0 {
			file += part
		} else if i == len(parts)-2 {
			file += part + "/"
		} else {
			file += string(part[0]) + "/"
		}
	}

	return file
}

func smartCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(short(runtime.FuncForPC(caller.PC).Name()))
}
