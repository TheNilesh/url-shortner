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

func NewInMemoryKVStore() KVStore {
	return &kvStore{}
}

type kvStore map[string]string

func (k *kvStore) Put(_ context.Context, key string, value string) error {
	kv := *k
	kv[key] = value
	return nil
}

func (k *kvStore) Get(_ context.Context, key string) (string, error) {
	kv := *k
	if val, ok := kv[key]; ok {
		return val, nil
	}
	return "", ErrKeyNotFound
}

func (k *kvStore) Exists(_ context.Context, key string) (bool, error) {
	kv := *k
	if _, ok := kv[key]; ok {
		return true, nil
	}
	return false, nil
}

func (k *kvStore) Delete(_ context.Context, key string) error {
	kv := *k
	delete(kv, key)
	return nil
}
