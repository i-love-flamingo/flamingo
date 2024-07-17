//go:build tracelog

package silentzap

// allowedLevels returns allowed levels for the trace logger.
const allowedLevels = `"Trace" | *"Debug" | "Info" | "Warn" | "Error" | "DPanic" | "Panic" | "Fatal"`
