package cache

import (
	"context"

	"flamingo.me/flamingo/v3/framework/opencensus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	backendCacheKeyType, _        = tag.NewKey("backend_type")
	backendCacheKeyErrorReason, _ = tag.NewKey("error_reason")
	backendCacheHitCount          = stats.Int64("flamingo/cache/backend/hit", "Count of cache-backend hits", stats.UnitDimensionless)
	backendCacheMissCount         = stats.Int64("flamingo/cache/backend/miss", "Count of cache-backend misses", stats.UnitDimensionless)
	backendCacheErrorCount        = stats.Int64("flamingo/cache/backend/error", "Count of cache-backend errors", stats.UnitDimensionless)
)

type (
	// BackendMetrics representation
	BackendMetrics struct {
		backendType string
	}
)

// NewBackendMetrics creates an backend metrics helper instance
func NewBackendMetrics(backendType string) BackendMetrics {
	b := BackendMetrics{
		backendType: backendType,
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

func (bi BackendMetrics) countHit() {
	ctx, _ := tag.New(context.Background(), tag.Upsert(opencensus.KeyArea, "cacheBackend"), tag.Upsert(backendCacheKeyType, bi.backendType))
	stats.Record(ctx, backendCacheHitCount.M(1))
}

func (bi BackendMetrics) countMiss() {
	ctx, _ := tag.New(context.Background(), tag.Upsert(opencensus.KeyArea, "cacheBackend"), tag.Upsert(backendCacheKeyType, bi.backendType))
	stats.Record(ctx, backendCacheMissCount.M(1))
}

func (bi BackendMetrics) countError(reason string) {
	ctx, _ := tag.New(context.Background(), tag.Upsert(opencensus.KeyArea, "cacheBackend"), tag.Upsert(backendCacheKeyType, bi.backendType), tag.Upsert(backendCacheKeyErrorReason, reason))
	stats.Record(ctx, backendCacheErrorCount.M(1))
}
