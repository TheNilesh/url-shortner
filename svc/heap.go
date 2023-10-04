package svc

import (
	"container/heap"
)

type KeyValuePair struct {
	Key   string
	Value int
}

type MaxHeap []KeyValuePair

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i].Value >= h[j].Value }
func (h MaxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MaxHeap) Push(x interface{}) {
	kv := x.(KeyValuePair)
	*h = append(*h, kv)
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	kv := old[n-1]
	*h = old[0 : n-1]
	return kv
}

func (h MaxHeap) FindIndexByKey(key string) int {
	for i, kv := range h {
		if kv.Key == key {
			return i
		}
	}
	return -1
}

func (h *MaxHeap) IncValue(key string) bool {
	if index := h.FindIndexByKey(key); index != -1 {
		(*h)[index].Value++
		heap.Fix(h, index)
		return true
	}
	return false
}

// IncOrPush push a new entry with a value of 1 if the key doesn't exist,
// or increments the value by one if the key exists.
func (h *MaxHeap) IncOrPush(key string) {
	if !h.IncValue(key) {
		heap.Push(h, KeyValuePair{Key: key, Value: 1})
	}
}

func (h MaxHeap) Top3MaxHeapKeysValues() []KeyValuePair {
	top3 := make([]KeyValuePair, 0)
	for i := 0; i < 3 && i < len(h); i++ {
		top3 = append(top3, h[i])
	}
	return top3
}
