package healthcheck

type (
	Nil struct{}
)

var (
	_ Status = &Nil{}
)

func (s *Nil) Status() (bool, string) {
	return true, "success"
}
