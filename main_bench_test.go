package main

import (
	"log"
	"sync"
	"testing"

	"github.com/throttled/throttled/v2"
	"github.com/throttled/throttled/v2/store/memstore"
)

func Benchmark(b *testing.B) {
	store, err := memstore.NewCtx(65536)
	if err != nil {
		log.Fatal(err)
	}

	quota := throttled.RateQuota{
		MaxRate:  throttled.PerMin(20),
		MaxBurst: 5,
	}
	rateLimiter, err := throttled.NewGCRARateLimiterCtx(store, quota)
	if err != nil {
		log.Fatal(err)
	}

	counter := struct {
		mutex sync.Mutex
		count int
	}{}

	for b.Loop() {
		limited, _, err := rateLimiter.RateLimitCtx(b.Context(), "key", 1)
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
