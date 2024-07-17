//go:build !tracelog

package zap

// allowedLevels returns allowed levels for the standard logger.
const allowedLevels = `*"Debug" | "Info" | "Warn" | "Error" | "DPanic" | "Panic" | "Fatal"`
