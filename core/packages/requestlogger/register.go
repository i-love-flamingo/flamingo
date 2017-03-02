package requestlogger

import (
	di "flamingo/core/flamingo/dependencyinjection"
)

// Register a request logger
func Register(c *di.Container) {
	c.Register(new(Logger), "event.subscriber")
}
