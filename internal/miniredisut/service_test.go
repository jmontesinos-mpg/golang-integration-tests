package miniredisut

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

// TestMiniRedis is a test to demonstrate a simple scenario using Redis as the storage layer.
func TestMiniRedis(t *testing.T) {
	// Start a new MiniRedis instance, it registers itself into the t.CleanUp, so no need to do that manually.
	s := miniredis.RunT(t)
	s.Set("test-key", "test-val")

	ctx := context.Background()
	val, err := redisClientCall(ctx, s.Addr())

	require.NoError(t, err)

	require.Equal(t, "test-val", val)
}

// redisClientCall simulates the business logic of an existing service calling redis to retrieve data from the storage layers.
func redisClientCall(ctx context.Context, addr string) (string, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	val, err := rdb.Get(ctx, "test-key").Result()
	if err != nil {
		return "", err
	}

	return val, nil
}
