package healthcheck

type (
	// Status is used to provide a healthcheck
	Status interface {
		Status() (alive bool, details string)
	}
)
