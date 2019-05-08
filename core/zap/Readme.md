# Zap Logger

The Zap module implements the flamingo `Logger` interface and uses Zap logger.

## Usage

Just add the module to your bootstrap, the module will Bind on the interface `flamingo.Logger` which is the interface you should inject in your application code to do logging.

Use the method `.WithContext(ctx)` whenever you have a context, this will make sure that traceId, spanId and session (if configured) are added to the log automatically.


## Configuration

```
zap:
  loglevel: Info
  json: true
  colored: false
  devmode: false
  logsession: true # to enable logging sessionid
  sampling: # see zap sampling doc
    enabled: false
    initial: 100
    thereafter: 100
  fieldmap:
    key: value
  
```

## More about Zap

 * https://godoc.org/go.uber.org/zap
 * https://github.com/uber-go/zap
 * https://github.com/uber-go/zap/blob/master/FAQ.md