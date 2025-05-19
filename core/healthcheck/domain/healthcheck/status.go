package healthcheck

type (
	// Status check interface
	Status interface {
		Status() (alive bool, details string)
	}

	// MeasuredStatus healthcheck interface which will be used in metrics gathering
	MeasuredStatus interface {
		Status
		Name() string
	}
)
