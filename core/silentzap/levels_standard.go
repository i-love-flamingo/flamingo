//go:build !tracelog

package silentzap

// allowedLevels returns allowed levels for the standard logger.
func allowedLevels() string {
	return `*"Debug" | "Info" | "Warn" | "Error" | "DPanic" | "Panic" | "Fatal"`
}
