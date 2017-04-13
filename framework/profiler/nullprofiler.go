package profiler

// NullProfiler does not profile
type NullProfiler struct{}

// Profile nothing
func (np *NullProfiler) Profile(string, string) ProfileFinishFunc {
	return func() {}
}
