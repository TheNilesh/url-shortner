package metrics

import (
	"container/heap"
)

type KeyValuePair struct {
	Key   string
	Value int
}

type Heap interface {
	heap.Interface
	IncOrPush(key string)
	GetMaxValuePairs(n int) []KeyValuePair
}

func NewHeap() Heap {
	h := &maxHeap{}
	heap.Init(h)
	return h
}

type maxHeap []KeyValuePair

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].Value > h[j].Value }
func (h maxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *maxHeap) Push(x interface{}) {
	kv := x.(KeyValuePair)
	*h = append(*h, kv)
}

func (h *maxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	kv := old[n-1]
	*h = old[0 : n-1]
	return kv
}

func (h maxHeap) findIndexByKey(key string) int {
	for i, kv := range h {
		if kv.Key == key {
			return i
		}
	}
	return -1
}

func (h *maxHeap) incValue(key string) bool {
	if index := h.findIndexByKey(key); index != -1 {
		(*h)[index].Value++
		heap.Fix(h, index)
		return true
	}
	return false
}

// IncOrPush push a new entry with a value of 1 if the key doesn't exist,
// or increments the value by one if the key exists.
func (h *maxHeap) IncOrPush(key string) {
	if !h.incValue(key) {
		heap.Push(h, KeyValuePair{Key: key, Value: 1})
	}
}

func (h maxHeap) GetMaxValuePairs(n int) []KeyValuePair {
	top3 := make([]KeyValuePair, 0)
	for i := 0; i < n && i < len(h); i++ {
		top3 = append(top3, h[i])
	}
	return top3
}
