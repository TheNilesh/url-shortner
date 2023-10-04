package store

import "context"

func NewGoMapStore() KVStore {
	return &goMapStore{}
}

type goMapStore map[string]string

func (k *goMapStore) Put(_ context.Context, key string, value string) error {
	kv := *k
	kv[key] = value
	return nil
}

func (k *goMapStore) Get(_ context.Context, key string) (string, error) {
	kv := *k
	if val, ok := kv[key]; ok {
		return val, nil
	}
	return "", ErrKeyNotFound
}

func (k *goMapStore) Exists(_ context.Context, key string) (bool, error) {
	kv := *k
	if _, ok := kv[key]; ok {
		return true, nil
	}
	return false, nil
}

func (k *goMapStore) Delete(_ context.Context, key string) error {
	kv := *k
	delete(kv, key)
	return nil
}
