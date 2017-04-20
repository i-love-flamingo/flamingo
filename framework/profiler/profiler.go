package profiler

type (
	// ProfileFinishFunc is used to finish a profiling-block
	ProfileFinishFunc func()

	// Profiler is used to measure the time certain things need
	Profiler interface {
		Profile(key, msg string) ProfileFinishFunc
	}
)
