# Healthcheck Module

The healthcheck module provides useful routes for:
1) check application is up (/status/ping endpoint)
2) Check application status (/status/healthcheck endpoint)

## Usage

The healthcheck module requires the systemendpoint module - so make sure both are loaded in your bootstrap:

```go
  ...
  new(systemendpoint.Module),
  new(healthcheck.Module),
  ...
``` 

If you use the default configuration you should then see an ok when calling *http://localhost:13210/status/ping*


## Healthcheck

The healthcheck endpoint allows to check different "things" from the view of your application.
There are some available checks that you can activate by config:

```
healthcheck.checkSession: true
healthcheck.checkAuth: true
```

### Implement own Checks:

Register new checks with Dingo:

```
injector.BindMap(new(healthcheck.Status), "session").To(healthcheck.RedisSession{})
```
