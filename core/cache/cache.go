package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Alexigbokwe/gonext-framework/core/config"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store defines the contract for caching
type Store interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Forget(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}

// NewStore creates a new cache store based on configuration
// For simplicity, if Redis config is present, we use Redis, otherwise Memory.
// Or we can add specific CACHE_DRIVER env.
func NewStore(cfg *config.Config) Store {
	// Simple heuristic: if Redis host is set, use Redis.
	// In a real app, introduce CACHE_DRIVER=redis|memory
	if cfg.Redis.Host != "" {
		addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
		rdb := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: cfg.Redis.Password,
		})
		return &RedisStore{client: rdb}
	}
	return NewMemoryStore()
}

// MemoryStore implementation
type MemoryStore struct {
	items sync.Map
}

type item struct {
	Value     []byte
	ExpiresAt int64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Get(ctx context.Context, key string, dest interface{}) error {
	val, ok := m.items.Load(key)
	if !ok {
		return fmt.Errorf("key not found")
	}

	it := val.(item)
	if it.ExpiresAt > 0 && time.Now().UnixNano() > it.ExpiresAt {
		m.items.Delete(key)
		return fmt.Errorf("key expired")
	}

	return json.Unmarshal(it.Value, dest)
}

func (m *MemoryStore) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	m.items.Store(key, item{Value: data, ExpiresAt: exp})
	return nil
}

func (m *MemoryStore) Forget(ctx context.Context, key string) error {
	m.items.Delete(key)
	return nil
}

func (m *MemoryStore) Flush(ctx context.Context) error {
	m.items = sync.Map{}
	return nil
}

// RedisStore implementation
type RedisStore struct {
	client *redis.Client
}

func (r *RedisStore) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(val, dest)
}

func (r *RedisStore) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisStore) Forget(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisStore) Flush(ctx context.Context) error {
	return r.client.FlushAll(ctx).Err()
}
