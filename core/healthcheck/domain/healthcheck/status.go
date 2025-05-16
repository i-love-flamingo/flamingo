package healthcheck

type (
	// Status check interface
	Status interface {
		Status() (alive bool, details string)
	}

	MeasuredStatus interface {
		Status
		Name() string
	}
)
