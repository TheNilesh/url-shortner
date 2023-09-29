package store

// KVStore is simple key value store interface that encapsulates
// underlying logic of kv-store functionality
type KVStore interface {
	Put(key string, value string) error
	Get(key string) (string, error)
}

// TODO: Implement KV store using golang map for poc and then redis
