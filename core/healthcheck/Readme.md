# Healthcheck Module

The healthcheck module provides useful routes for:
1. check if the application is up (/status/ping endpoint)
1. Check the application status (/status/healthcheck endpoint)

## Usage

The healthcheck module requires the systemendpoint module and its second server.

If you use the default configuration you should then see an ok when calling [http://localhost:13210/status/ping](http://localhost:13210/status/ping)


## Healthcheck

The healthcheck endpoint allows to check different "things" from the view of your application.
There are some available checks that you can activate by config:

```yaml
healthcheck.checkSession: true
healthcheck.checkAuth: true
```

### Implement own Checks:

Just Implement the `healthcheck.Status` interface and register it via Dingo mapbinding:

```go
injector.BindMap(new(healthcheck.Status), "session").To(healthcheck.RedisSession{})
```
