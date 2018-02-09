package application

import (
	"go.aoe.com/flamingo/core/auth/domain"
	"go.aoe.com/flamingo/framework/web"
)

type (
	// EventPublisher
	EventPublisher struct {
	}
)

func (e *EventPublisher) PublishLoginEvent(ctx web.Context, event *domain.LoginEvent) {
	//publish to Flamingo default Event Router
	ctx.EventRouter().Dispatch(event)
}

func (e *EventPublisher) PublishLogoutEvent(ctx web.Context, event *domain.LogoutEvent) {
	//publish to Flamingo default Event Router
	ctx.EventRouter().Dispatch(event)
}
