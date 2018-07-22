package profiler

type (
	// ProfileFinishFunc is used to finish a profiling-block
	// deprecated: use opencensus
	ProfileFinishFunc func()

	// Profiler is used to measure the time certain things need
	// deprecated: use opencensus
	Profiler interface {
		// deprecated: use opencensus
		Profile(key, msg string) ProfileFinishFunc
	}
)
