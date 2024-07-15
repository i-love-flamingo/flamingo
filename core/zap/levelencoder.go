package zap

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

func capitalLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if l == logLevels["Trace"] {
		enc.AppendString("TRACE")
		return
	}

	zapcore.CapitalLevelEncoder(l, enc)
}

func capitalColorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if l == logLevels["Trace"] {
		enc.AppendString(fmt.Sprintf("\x1b[36m%s\x1b[0m", "TRACE")) // Cyan colored TRACE
		return
	}

	zapcore.CapitalColorLevelEncoder(l, enc)
}
