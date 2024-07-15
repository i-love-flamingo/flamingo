//go:build tracelog

package zap

// allowedLevels returns allowed levels for the trace logger.
func allowedLevels() string {
	return `"Trace" | *"Debug" | "Info" | "Warn" | "Error" | "DPanic" | "Panic" | "Fatal"`
}
