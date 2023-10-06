package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoMapStore(t *testing.T) {
	store := NewGoMapStore()

	err := store.Put(context.Background(), "key1", "value1")
	if err != nil {
		t.Errorf("Error putting key-value pair: %v", err)
	}
	val, err := store.Get(context.Background(), "key1")
	if err != nil {
		t.Errorf("Error getting value for key: %v", err)
	}
	if val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	val, err = store.Get(context.Background(), "notfoundkey")
	if val != "" {
		t.Errorf("Expected empty value, got %v", val)
	}
	if err != nil {
		assert.Equal(t, err, ErrKeyNotFound)
	}

	exists, err := store.Exists(context.Background(), "key1")
	if err != nil {
		t.Errorf("Error checking if key exists: %v", err)
	}
	if !exists {
		t.Errorf("Expected key1 to exist")
	}
	exists, err = store.Exists(context.Background(), "key2")
	if err != nil {
		t.Errorf("Error checking if key exists: %v", err)
	}
	if exists {
		t.Errorf("Expected key2 to not exist")
	}

	err = store.Delete(context.Background(), "key1")
	if err != nil {
		t.Errorf("Error deleting key: %v", err)
	}
	exists, err = store.Exists(context.Background(), "key1")
	if err != nil {
		t.Errorf("Error checking if key exists: %v", err)
	}
	if exists {
		t.Errorf("Expected key1 to not exist")
	}
}
