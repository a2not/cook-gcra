package main

import (
	"sync"
	"testing"

	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
)

func BenchmarkThrottled(b *testing.B) {
	store, err := memstore.NewCtx(65536)
	if err != nil {
		b.Fatal(err)
	}

	quota := throttled.RateQuota{
		MaxRate:  throttled.PerSec(RatePerSec),
		MaxBurst: BurstSize,
	}
	limiter, err := throttled.NewGCRARateLimiterCtx(store, quota)
	if err != nil {
		b.Fatal(err)
	}

	counter := struct {
		mutex sync.Mutex
		count int
	}{}

	for b.Loop() {
		limited, _, err := limiter.RateLimitCtx(b.Context(), "key", 1)
		if err != nil {
			b.Fatal(err)
		}

		if !limited {
			counter.mutex.Lock()
			counter.count++
			counter.mutex.Unlock()
		}
	}

	b.Logf("%d requests allowed", counter.count)
}
