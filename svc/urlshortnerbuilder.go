package svc

import (
	"github.com/thenilesh/url-shortner/metrics"
	"github.com/thenilesh/url-shortner/store"
)

type URLShortnerBuilder struct {
	minLength      int
	maxLength      int
	charset        string
	targetURLStore store.KVStore
	shortPathStore store.KVStore
	metrics        metrics.Metrics
}

func NewURLShortnerBuilder() *URLShortnerBuilder {
	return &URLShortnerBuilder{
		minLength: 4,
		maxLength: 7,
		charset:   "abcdefghijklmnopqrstuvwxyz0123456789",
	}
}

func (b *URLShortnerBuilder) SetMinLength(minLength int) *URLShortnerBuilder {
	b.minLength = minLength
	return b
}

func (b *URLShortnerBuilder) SetMaxLength(maxLength int) *URLShortnerBuilder {
	b.maxLength = maxLength
	return b
}

func (b *URLShortnerBuilder) SetCharset(charset string) *URLShortnerBuilder {
	b.charset = charset
	return b
}

func (b *URLShortnerBuilder) SetTargetURLStore(store store.KVStore) *URLShortnerBuilder {
	b.targetURLStore = store
	return b
}

func (b *URLShortnerBuilder) SetShortPathStore(store store.KVStore) *URLShortnerBuilder {
	b.shortPathStore = store
	return b
}

func (b *URLShortnerBuilder) SetMetrics(metrics metrics.Metrics) *URLShortnerBuilder {
	b.metrics = metrics
	return b
}

func (b *URLShortnerBuilder) Build() URLShortner {
	if b.targetURLStore == nil {
		panic("targetURLStore is nil")
	}
	if b.shortPathStore == nil {
		panic("shortPathStore is nil")
	}
	if b.metrics == nil {
		panic("metrics is nil")
	}
	if b.minLength > b.maxLength {
		panic("minLength is greater than maxLength")
	}
	if b.minLength <= 0 {
		panic("minLength is less than or equal to 0")
	}
	if b.maxLength <= 0 {
		panic("maxLength is less than or equal to 0")
	}
	if b.minLength > 50 || b.maxLength > 50 {
		panic("minLength or maxLength is greater than 50")
	}
	randomStrGen := NewRandomStrGen(b.minLength, b.maxLength, b.charset)
	return &urlShortner{
		randomStrGen:   randomStrGen,
		targetURLStore: b.targetURLStore,
		shortPathStore: b.shortPathStore,
		metrics:        b.metrics,
	}
}
