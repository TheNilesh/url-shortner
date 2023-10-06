package svc

import (
	"errors"

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

func (b *URLShortnerBuilder) Build() (URLShortner, error) {
	if b.targetURLStore == nil {
		return nil, errors.New("targetURLStore is nil")
	}
	if b.shortPathStore == nil {
		return nil, errors.New("shortPathStore is nil")
	}
	if b.metrics == nil {
		return nil, errors.New("metrics is nil")
	}
	if b.minLength <= 0 || b.maxLength <= 0 {
		return nil, errors.New("minLength or maxLength is less than or equal to 0")
	}
	if b.minLength > 50 || b.maxLength > 50 {
		return nil, errors.New("minLength or maxLength is greater than 50")
	}
	if b.minLength > b.maxLength {
		return nil, errors.New("minLength is greater than maxLength")
	}
	if !isValidPathSegment(b.charset) {
		return nil, errors.New("charset contains invalid characters")
	}

	randomStrGen := NewRandomStrGen(b.minLength, b.maxLength, b.charset)
	return &urlShortner{
		randomStrGen:   randomStrGen,
		targetURLStore: b.targetURLStore,
		shortPathStore: b.shortPathStore,
		metrics:        b.metrics,
	}, nil
}
