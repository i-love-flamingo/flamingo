# Changelog

## v3

- `framework/web.Response` is not `framework/web.Result`
- `framework/web.Request` is heavily condensed
- `framework/web.Session` does not expose `.GS()` for the gorilla session anymore
- `framework/router` package merged into `framework/web`
- `framework/event` package merged into `framework/flamingo`
- `core/requestTask` is renamed to `core/requesttask`
- `core/canonicalUrl` is renamed to `core/canonicalurl`
- `event.Subscriber` is getting `context.Context` as the first argument: `Notify(ctx context.Context, e flamingo.Event)`
- `event.Subscriber` are registered via `framework/flamingo.BindEventSubscriber(injector).To(...)`
