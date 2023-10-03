package store

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string, password string, db int) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	}), nil
}

func CheckRedisConnection(client *redis.Client) error {
	_, err := client.Ping(context.Background()).Result()
	return err
}

// redisKVStore represents a key-value store backed by Redis.
type redisKVStore struct {
	client    *redis.Client
	namespace string
}

func NewRedisKVStore(client *redis.Client, namespace string) (*redisKVStore, error) {
	return &redisKVStore{
		client:    client,
		namespace: namespace,
	}, nil
}

func (store *redisKVStore) Put(ctx context.Context, key string, value string) error {
	return store.client.Set(ctx, fmt.Sprintf("%s:%s", store.namespace, key), value, 0).Err()
}

func (store *redisKVStore) Get(ctx context.Context, key string) (string, error) {
	val, err := store.client.Get(ctx, fmt.Sprintf("%s:%s", store.namespace, key)).Result()
	if err == redis.Nil {
		return "", ErrKeyNotFound
	}
	return val, err
}

func (store *redisKVStore) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := store.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (store *redisKVStore) Delete(ctx context.Context, key string) error {
	return store.client.Del(ctx, fmt.Sprintf("%s:%s", store.namespace, key)).Err()
}
