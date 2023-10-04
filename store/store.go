package store

import (
	"context"
	"errors"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

// KVStore is simple key value store interface that encapsulates
// underlying logic of kv-store functionality
type KVStore interface {
	Put(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
	Delete(ctx context.Context, key string) error
}
