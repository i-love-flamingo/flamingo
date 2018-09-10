# Opencensus

## General usage

First, if you're not sure what opencensus is, please visit [https://opencensus.io/](https://opencensus.io/) to learn more about it.

Opencensus allows you to collect metrics for your application and process them via e.g. Prometheus, Grafana, etc.
The core package already collects metrics automatically for routers, prefixrouters and pugtemplate render buckets out of the box, with no additional code necessary.

## Adding your own metrics

The most likely usecase is in a controller, but to have your own metric, please import the following packages:

* "flamingo.me/flamingo/framework/opencensus"
* "go.opencensus.io/stats"
* "go.opencensus.io/stats/view"

Then, create a variable for your metric. See below example:

```go
var rt = stats.Int64("flamingo/package/mystat", "my stat records 5 milliseconds per call", stats.UnitMilliseconds)
```

The Variable name of rt is totally arbitrary here. But "flamingo/package/mystat" is actually the metric id you're recording to. Normally, this should be a pattern of
"yourproject/subpackage/yourstat". Followed up by a description of your metric and the Unit which is recorded (use stats.UnitDimensionless for simple counters).

To be actuall able to transmit this metic, register it in an init method:

```go
	// register view in opencensus
	func init() {
		opencensus.View("flamingo/package/mystat/sum", rt, view.Sum())
		opencensus.View("flamingo/package/mystat/count", rt, view.Count())
	}
```

Thats all that is needed as Setup. You can add to your metric in your code now. It's as simple as adding an increment:

```go
    // measure mystat
	func (c *Controller) Get(ctx context.Context, r *web.Request) web.Response {
		// ...

		// record 5ms per call
		stats.Record(ctx, rt.M(5))

		// ...
	}
```

For a simple counter where rt is a stats.UnitDimensionless metric, you´d simply increment by one, i.e. 

```go
		// record increment of 1 per call
		stats.Record(ctx, rt.M(1))
```
