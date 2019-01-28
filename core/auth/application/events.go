package application

import (
	"context"
	"flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/event"
)

// EventPublisher struct
type EventPublisher struct {
	router event.Router
}

func (e *EventPublisher) Inject(router event.Router) {
	e.router = router
}

// PublishLoginEvent dispatches the login event on the contexts event router
func (e *EventPublisher) PublishLoginEvent(ctx context.Context, event *domain.LoginEvent) {
	e.router.Dispatch(ctx, event)
}

// PublishLogoutEvent dispatches the logout event on the contexts event router
func (e *EventPublisher) PublishLogoutEvent(ctx context.Context, event *domain.LogoutEvent) {
	e.router.Dispatch(ctx, event)
}


type (
	EventHandler struct {
		authManager *AuthManager
	}
)

func (e *EventHandler) Inject(authManager *AuthManager) {
	e.authManager = authManager
}

// Notify calls AuthManager on each logout, so it can destroy data stored for previously logged in user
func (e *EventHandler) Notify(event event.Event) {
	logoutEvent, ok := event.(*domain.LogoutEvent)
	if ok {
		e.authManager.DeleteTokenDetails(logoutEvent.Session)
		e.authManager.DeleteAuthState(logoutEvent.Session)
	}
}
