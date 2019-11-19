# Cache module
The Cache module provides an easy interface to cache things in flamingo.

The basic concept is, that there is a so called "cache frontend" - that offers an interface to cache certain types, and a "cache backend" that takes care about storing(persiting) the cache entry.

## Caching HTTP responses from APIs

A typical use case is, to cache responses from (slow) backends that you need to call.

First define an injection to get a HTTPFrontend cache injected:

```go
MyApiClient struct {
      Cache            *cache.HTTPFrontend                  `inject:"myservice"`
}
```


We use annotation to be able to individually configure the requested Cache. So our binding may look like:

```go
injector.Bind((*cache.HTTPFrontend)(nil)).AnnotatedWith("myservice").In(dingo.Singleton)
```

Then when you request your API you can wrap the result in the cache and provide a cache loader function.
The HTTPFrontend Cache then makes sure that:
- If there is a cache hit - return it within the given "GraceTime" - and eventually do a new request if cache "LiveTime" is over.
- Requests for the same key are done in "single flight" if cache is empty meaning that even if there may be 1000 parallel requests only the first one will be executed against the backend service and the other wait for the result

Example:
```go
loadData :=  func(ctx context.Context) (*http.Response, *cache.Meta, error) {
    r, err := http.DefaultClient.Do(req.WithContext(ctx))
    if err != nil {
        return nil, nil, err
    }
    //cache semantic errors for certain time to avoid recalling the same request
    if r.StatusCode == http.StatusNotFound {
        return r, &cache.Meta{
            Lifetime:  5 * time.Minute,
            Gracetime: 300 * time.Second,
        }, nil
    }
    //cache semantic errors for certain time to avoid recalling the same request
    if r.StatusCode != http.StatusOK {
        return r, &cache.Meta{
            Lifetime:  10 * time.Second,
            Gracetime: 30 * time.Second,
        }, nil
    }
    return r, nil, nil
}
response, err := apiclient.Cache.Get(requestContext, u.String(), loadData)
```

## Cache backends

Currently there are the following backends available:

### inMemoryCache

Caches in memory - and therefore is a very fast cache.

It is base on the LRU-Strategy witch drops least used entries. For this reason the cache will be no overcommit your memory and will automicly fit the need of your current traffic.

### redisBackend

Is using [redis](https://redis.io/) as an shared inMemory cache.
Since all cache-fetched has an overhead to the inMemoryBackend, the redis is a little slower.
The benefit of redis is the shared storage an the high efficiency in reading and writing keys. especialy if you need scale fast horizonaly, it helps to keep your backend-systems healthy.

Be ware of using redis (or any other shared cache backend) as an single backend, because of network latency. (have a loo at the multiLevelBackend)

### fileBackend

Writes the cache content to the local filesystem.

### nullBackend

Caches nothing.

### multiLevelBackend

The multiLevelBackend was introduced to get the benefit of the extrem fast inMemorybackend and a shared backend.
Using the inMemoryBackend in combination with an shared backend, gives you blazing fast responces and helps you to protect you backend in case of fast scaleout-scenarios.

@TODO: Write example code.
