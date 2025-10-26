package main

import (
	"sync"
	"testing"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func BenchmarkRedisRate(b *testing.B) {
	ctx := b.Context()
	redisC, err := testcontainers.Run(
		ctx, "redis:latest",
		testcontainers.WithExposedPorts("6379/tcp"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("6379/tcp"),
			wait.ForLog("Ready to accept connections"),
		),
	)
	if err != nil {
		b.Fatal(err)
	}
	endpoint, err := redisC.Endpoint(ctx, "")
	if err != nil {
		b.Error(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: endpoint,
	})
	_ = rdb.FlushDB(ctx).Err()

	quota := redis_rate.PerSecond(100)
	limiter := redis_rate.NewLimiter(rdb)

	counter := struct {
		mutex sync.Mutex
		count int
	}{}

	for b.Loop() {
		res, err := limiter.Allow(ctx, "key", quota)
		if err != nil {
			b.Fatal(err)
		}

		if res.Allowed > 0 {
			counter.mutex.Lock()
			counter.count++
			counter.mutex.Unlock()
		}
	}

	b.Logf("%d requests allowed", counter.count)

	testcontainers.CleanupContainer(b, redisC)
	require.NoError(b, err)
}
