package store

// KVStore is simple key value store interface that encapsulates
// underlying logic of kv-store functionality
type KVStore interface {
	Put(key string, value string) error
	Get(key string) (string, error)
	Exists(key string) (bool, error)
}

func NewKVStore() KVStore {
	return &kvStore{}
}

// TODO: Implement KV store using golang map for poc and then redis
type kvStore map[string]string

func (k *kvStore) Put(key string, value string) error {
	kv := *k
	kv[key] = value
	return nil
}

func (k *kvStore) Get(key string) (string, error) {
	kv := *k
	if val, ok := kv[key]; ok {
		return val, nil
	}
	return "", nil
}

func (k *kvStore) Exists(key string) (bool, error) {
	kv := *k
	if _, ok := kv[key]; ok {
		return true, nil
	}
	return false, nil
}
