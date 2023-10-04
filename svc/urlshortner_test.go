package svc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thenilesh/url-shortner/metrics"
	"github.com/thenilesh/url-shortner/store"
)

type mockTargetURLStore struct {
	mock.Mock
}

func (m *mockTargetURLStore) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *mockTargetURLStore) Put(ctx context.Context, key, value string) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *mockTargetURLStore) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *mockTargetURLStore) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

type mockShortPathStore struct {
	mock.Mock
}

func (m *mockShortPathStore) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *mockShortPathStore) Put(ctx context.Context, key, value string) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *mockShortPathStore) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *mockShortPathStore) Exists(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

type mockMetrics struct {
	mock.Mock
}

func (m *mockMetrics) Start() {
	m.Called()
}

func (m *mockMetrics) IncDomainCount(domain string) {
	m.Called(domain)
}

func (m *mockMetrics) GetDomainCounts() []metrics.KeyValuePair {
	args := m.Called()
	return args.Get(0).([]metrics.KeyValuePair)
}

func TestURLShortner_CreateShortPath(t *testing.T) {
	ctx := context.Background()

	targetURLStore := new(mockTargetURLStore)
	shortPathStore := new(mockShortPathStore)
	metrics := new(mockMetrics)

	shortner := NewURLShortner(6, targetURLStore, shortPathStore, metrics)

	shortPathExpected := "abc123"
	targetURLExpected := "https://www.google.com"
	targetURLStore.On("Get", ctx, shortPathExpected).Return("", store.ErrKeyNotFound)
	targetURLStore.On("Exists", ctx, shortPathExpected).Return(false, nil)
	targetURLStore.On("Put", ctx, shortPathExpected, targetURLExpected).Return(nil)
	// targetURLStore.On("Delete", ctx, shortPathExpected).Return(nil)
	shortPathStore.On("Get", ctx, targetURLExpected).Return("", store.ErrKeyNotFound)
	shortPathStore.On("Put", ctx, targetURLExpected, shortPathExpected).Return(nil)
	metrics.On("IncDomainCount", "www.google.com").Once()

	shortPath, err := shortner.CreateShortPath(ctx, shortPathExpected, targetURLExpected)
	assert.NoError(t, err)
	assert.Equal(t, shortPathExpected, shortPath)

	// targetURLStore.On("Get", ctx, "google").Return("https://www.google.com", nil)

	// metrics.On("IncDomainCount", ctx, "shortpath_created").Once()
	// metrics.On("IncDomainCount", ctx, "shortpath_create_failed").Once()
	// _, err = shortner.CreateShortPath(ctx, "abc123", "https://www.github.com")
	// assert.Error(t, err)

	// shortPath, err = shortner.CreateShortPath(ctx, "", "https://www.google.com")
	// assert.NoError(t, err)
	// assert.Len(t, shortPath, 6)

	// shortPath, err = shortner.CreateShortPath(ctx, "", "https://www.github.com")
	// assert.NoError(t, err)
	// assert.Len(t, shortPath, 6)

	targetURLStore.AssertExpectations(t)
	shortPathStore.AssertExpectations(t)
	metrics.AssertExpectations(t)
}

func TestURLShortner_GetTargetURL(t *testing.T) {
	ctx := context.Background()

	targetURLStore := new(mockTargetURLStore)
	shortPathStore := new(mockShortPathStore)
	metrics := new(mockMetrics)

	shortner := NewURLShortner(6, targetURLStore, shortPathStore, metrics)

	targetURLStore.On("Get", ctx, "abc123").Return("https://www.google.com", nil)
	targetURLStore.On("Get", ctx, "invalid_shortpath").Return("", errors.New("not found"))

	targetURL, err := shortner.GetTargetURL(ctx, "abc123")
	assert.NoError(t, err)
	assert.Equal(t, "https://www.google.com", targetURL)

	_, err = shortner.GetTargetURL(ctx, "invalid_shortpath")
	assert.Error(t, err)

	targetURLStore.AssertExpectations(t)
	shortPathStore.AssertExpectations(t)
	metrics.AssertExpectations(t)
}
