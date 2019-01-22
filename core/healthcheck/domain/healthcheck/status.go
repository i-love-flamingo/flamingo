package healthcheck

// Status check interface
type Status interface {
	Status() (alive bool, details string)
}
