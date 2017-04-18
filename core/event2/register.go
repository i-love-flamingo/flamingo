package event2

import (
	di "flamingo/framework/dependencyinjection"
)

// Register the EventDispatcher
func Register(c *di.Container) {
	c.Register(new(DefaultEventDispatcher))
}
