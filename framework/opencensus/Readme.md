# Opencensus

## General usage

First, if you're not sure what opencensus is, please visit [https://opencensus.io/](https://opencensus.io/) to learn more about it.

Opencensus allows you to collect data for your application and process them via e.g. Prometheus, Jaeger, etc.
The core package already collects metrics automatically for routers and prefixrouters rendering times out of the box.
Also traces for request handling are done by default.

The metrics endpoint is provided under the systemendpoint. Once the module is activate you can access them via `http://localhost:13210/metrics`
## Adding your own metrics

The most likely usecase is in a controller, but to have your own metric, please import the following packages:

* "flamingo.me/flamingo/framework/opencensus"
* "go.opencensus.io/stats"
* "go.opencensus.io/stats/view"

Then, create a variable for your metric. See below example:

```go
var stat = stats.Int64("flamingo/package/mystat", "my stat records 5 milliseconds per call", stats.UnitMilliseconds)
```

The Variable name of `stat` is totally arbitrary here. But "flamingo/package/mystat" is actually the metric id you're recording to. Normally, this should be a pattern of
"yourproject/subpackage/yourstat". Followed up by a description of your metric and the Unit which is recorded (use `stats.UnitDimensionless` for simple counters).

To be able to transmit this metric, register it in an init method:

```go
// register view in opencensus
func init() {
  opencensus.View("flamingo/package/mystat/sum", stat, view.Sum())
  opencensus.View("flamingo/package/mystat/count", stat, view.Count())
}
```

That's all that is needed as Setup. You can add to your metric in your code now. It's as simple as adding an increment:

```go
func (c *Controller) Get(ctx context.Context, r *web.Request) web.Result {
  // ...

  // record 5ms per call
  stats.Record(ctx, stat.M(5))

  // ...
}
```

For a simple counter where `stat` is a stats.UnitDimensionless metric, you'd simply increment by one, i.e. 

```go
// record increment of 1 per call
stats.Record(ctx, stat.M(1))
```
