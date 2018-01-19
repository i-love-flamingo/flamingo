package healthcheck

type (
	Status interface {
		Status() (alive bool, details string)
	}
)
