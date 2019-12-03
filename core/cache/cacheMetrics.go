package cache

import (
	"context"

	"flamingo.me/flamingo/v3/framework/opencensus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	backendTypeCacheKeyType, _        = tag.NewKey("backend_type")
	frontendNameCacheKeyType, _        = tag.NewKey("frontend_name")
	backendCacheKeyErrorReason, _ = tag.NewKey("error_reason")
	backendCacheHitCount          = stats.Int64("flamingo/cache/backend/hit", "Count of cache-backend hits", stats.UnitDimensionless)
	backendCacheMissCount         = stats.Int64("flamingo/cache/backend/miss", "Count of cache-backend misses", stats.UnitDimensionless)
	backendCacheErrorCount        = stats.Int64("flamingo/cache/backend/error", "Count of cache-backend errors", stats.UnitDimensionless)
)

type (
	// CacheMetrics take care of publishing metrics for a specific cache
	CacheMetrics struct {
		//backendType - the type of the cache backend
		backendType string
		//frontendName - the name if the cache frontend where the backend is attached
		frontendName string
	}
)

// NewCacheMetrics creates an backend metrics helper instance
func NewCacheMetrics(backendType string, frontendName string) CacheMetrics {
	b := CacheMetrics{
		backendType: backendType,
		frontendName: frontendName,
	}
	return b
}

func init() {
	if err := opencensus.View("flamingo/cache/backend/hit", backendCacheHitCount, view.Count()); err != nil {
		panic(err)
	}
	if err := opencensus.View("flamingo/cache/backend/miss", backendCacheMissCount, view.Count()); err != nil {
		panic(err)
	}
	if err := opencensus.View("flamingo/cache/backend/error", backendCacheErrorCount, view.Count()); err != nil {
		panic(err)
	}
}

func (bi CacheMetrics) countHit() {
	ctx, _ := tag.New(context.Background(), tag.Upsert(opencensus.KeyArea, "cacheBackend"), tag.Upsert(backendTypeCacheKeyType, bi.backendType),  tag.Upsert(frontendNameCacheKeyType, bi.frontendName))
	stats.Record(ctx, backendCacheHitCount.M(1))
}

func (bi CacheMetrics) countMiss() {
	ctx, _ := tag.New(context.Background(), tag.Upsert(opencensus.KeyArea, "cacheBackend"), tag.Upsert(backendTypeCacheKeyType, bi.backendType), tag.Upsert(frontendNameCacheKeyType, bi.frontendName))
	stats.Record(ctx, backendCacheMissCount.M(1))
}

func (bi CacheMetrics) countError(reason string) {
	ctx, _ := tag.New(context.Background(), tag.Upsert(opencensus.KeyArea, "cacheBackend"), tag.Upsert(backendTypeCacheKeyType, bi.backendType), tag.Upsert(frontendNameCacheKeyType, bi.frontendName), tag.Upsert(backendCacheKeyErrorReason, reason))
	stats.Record(ctx, backendCacheErrorCount.M(1))
}
