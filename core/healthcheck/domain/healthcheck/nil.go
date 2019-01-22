package healthcheck

// Nil healtcheck
type Nil struct{}

var _ Status = &Nil{}

// Status is always healthy
func (s *Nil) Status() (bool, string) {
	return true, "success"
}
