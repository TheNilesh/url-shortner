package svc

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thenilesh/url-shortner/store"

	"github.com/thenilesh/url-shortner/mocks"
)

func TestURLShortner_CreateShortPath(t *testing.T) {
	ctx := context.Background()

	targetURLStore := new(mocks.KVStore)
	shortPathStore := new(mocks.KVStore)
	metrics := new(mocks.Metrics)
	collector := new(mocks.Collector)

	shortner, _ := NewURLShortnerBuilder().
		SetTargetURLStore(targetURLStore).
		SetShortPathStore(shortPathStore).
		SetMetrics(metrics).
		Build()

	shortPathExpected := "abc123"
	targetURLExpected := "https://www.google.com"
	targetURLStore.On("Get", ctx, shortPathExpected).Return("", store.ErrKeyNotFound)
	targetURLStore.On("Put", ctx, shortPathExpected, targetURLExpected).Return(nil)
	shortPathStore.On("Get", ctx, targetURLExpected).Return("", store.ErrKeyNotFound)
	shortPathStore.On("Put", ctx, targetURLExpected, shortPathExpected).Return(nil)
	collector.On("Inc", mock.Anything).Once()
	metrics.On("GetCollector", "domain_shortens").Return(collector)

	shortPath, err := shortner.CreateShortPath(ctx, shortPathExpected, targetURLExpected)
	assert.NoError(t, err)
	assert.Equal(t, shortPathExpected, shortPath)

	targetURLStore.AssertExpectations(t)
	shortPathStore.AssertExpectations(t)
	metrics.AssertExpectations(t)
}

func TestURLShortner_CreateShortPath_random(t *testing.T) {
	ctx := context.Background()

	targetURLStore := new(mocks.KVStore)
	shortPathStore := new(mocks.KVStore)
	metrics := new(mocks.Metrics)
	collector := new(mocks.Collector)

	expMinLen := 2
	expMaxLen := 3
	expCharset := "abc"
	shortner, _ := NewURLShortnerBuilder().
		SetTargetURLStore(targetURLStore).
		SetShortPathStore(shortPathStore).
		SetMetrics(metrics).
		SetCharset(expCharset).
		SetMinLength(expMinLen).
		SetMaxLength(expMaxLen).
		Build()

	emptyShortPath := "" // auto generate
	targetURLExpected := "https://www.example.com/somerajsjf"
	targetURLStore.On("Exists", ctx, mock.Anything).Return(false, nil)
	targetURLStore.On("Put", ctx, mock.Anything, targetURLExpected).Return(nil)
	shortPathStore.On("Get", ctx, targetURLExpected).Return("", store.ErrKeyNotFound)
	shortPathStore.On("Put", ctx, targetURLExpected, mock.Anything).Return(nil)
	collector.On("Inc", mock.Anything).Once()
	metrics.On("GetCollector", "domain_shortens").Return(collector)

	shortPath, err := shortner.CreateShortPath(ctx, emptyShortPath, targetURLExpected)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(shortPath), expMinLen)
	assert.LessOrEqual(t, len(shortPath), expMaxLen)
	for _, c := range shortPath {
		assert.Contains(t, "abc", string(c))
	}

	targetURLStore.AssertExpectations(t)
	shortPathStore.AssertExpectations(t)
	metrics.AssertExpectations(t)
}

func TestURLShortner_CreateShortPath_http(t *testing.T) {
	ctx := context.Background()

	targetURLStore := new(mocks.KVStore)
	shortPathStore := new(mocks.KVStore)
	metrics := new(mocks.Metrics)
	collector := new(mocks.Collector)

	shortner, _ := NewURLShortnerBuilder().
		SetTargetURLStore(targetURLStore).
		SetShortPathStore(shortPathStore).
		SetMetrics(metrics).
		Build()

	shortPathExpected := "abc123"
	targetURLExpected := "http://www.google.com/"
	targetURL := "http://www.google.com"
	targetURLStore.On("Get", ctx, shortPathExpected).Return("", store.ErrKeyNotFound)
	targetURLStore.On("Put", ctx, shortPathExpected, targetURL).Return(nil)
	shortPathStore.On("Get", ctx, targetURL).Return("", store.ErrKeyNotFound)
	shortPathStore.On("Put", ctx, targetURL, shortPathExpected).Return(nil)
	collector.On("Inc", mock.Anything).Once()
	metrics.On("GetCollector", "domain_shortens").Return(collector)

	shortPath, err := shortner.CreateShortPath(ctx, shortPathExpected, targetURLExpected)
	assert.NoError(t, err)
	assert.Equal(t, shortPathExpected, shortPath)

	targetURLStore.AssertExpectations(t)
	shortPathStore.AssertExpectations(t)
	metrics.AssertExpectations(t)
}

func TestURLShortner_CreateShortPath_invalid_path_url(t *testing.T) {
	ctx := context.Background()

	targetURLStore := new(mocks.KVStore)
	shortPathStore := new(mocks.KVStore)
	metrics := new(mocks.Metrics)

	shortner, _ := NewURLShortnerBuilder().
		SetTargetURLStore(targetURLStore).
		SetShortPathStore(shortPathStore).
		SetMetrics(metrics).
		Build()

	targetURLExpected := "https://www.google.com"

	shortPath, err := shortner.CreateShortPath(ctx, "abc/123", targetURLExpected)
	assert.Equal(t, "", shortPath)
	assert.Error(t, err, "Expected an error")
	_, ok := err.(*ErrValidation)
	assert.True(t, ok, "Expected error of type ErrValidation")
	assert.EqualError(t, err, "short_path contains disallowed characters")

	shortPath, err = shortner.CreateShortPath(ctx, "sp ", targetURLExpected)
	assert.Equal(t, "", shortPath)
	assert.Error(t, err, "Expected an error")
	_, ok = err.(*ErrValidation)
	assert.True(t, ok, "Expected error of type ErrValidation")
	assert.EqualError(t, err, "short_path contains leading or trailing spaces")

	shortPath, err = shortner.CreateShortPath(ctx, " sp", targetURLExpected)
	assert.Equal(t, "", shortPath)
	assert.Error(t, err, "Expected an error")
	_, ok = err.(*ErrValidation)
	assert.True(t, ok, "Expected error of type ErrValidation")
	assert.EqualError(t, err, "short_path contains leading or trailing spaces")

	shortPath, err = shortner.CreateShortPath(ctx, "metrics", targetURLExpected)
	assert.Equal(t, "", shortPath)
	assert.Error(t, err, "Expected an error")
	_, ok = err.(*ErrValidation)
	assert.True(t, ok, "Expected error of type ErrValidation")
	assert.EqualError(t, err, "short_path is reserved")

	shortPath, err = shortner.CreateShortPath(ctx, strings.Repeat("a", 51), targetURLExpected)
	assert.Equal(t, "", shortPath)
	assert.Error(t, err, "Expected an error")
	_, ok = err.(*ErrValidation)
	assert.True(t, ok, "Expected error of type ErrValidation")
	assert.EqualError(t, err, "short_path is too long")

	shortPath, err = shortner.CreateShortPath(ctx, "@", targetURLExpected)
	assert.Equal(t, "", shortPath)
	assert.Error(t, err, "Expected an error")
	_, ok = err.(*ErrValidation)
	assert.True(t, ok, "Expected error of type ErrValidation")

	expectedShortPath := "my-short-path"
	shortPath, err = shortner.CreateShortPath(ctx, expectedShortPath, "http://www.google.com ")
	assert.Equal(t, "", shortPath)
	assert.Error(t, err, "Expected an error")
	_, ok = err.(*ErrValidation)
	assert.True(t, ok, "Expected error of type ErrValidation")
	assert.EqualError(t, err, "target_url contains leading or trailing spaces")

	shortPath, err = shortner.CreateShortPath(ctx, expectedShortPath, " http://www.google.com")
	assert.Equal(t, "", shortPath)
	assert.Error(t, err, "Expected an error")
	_, ok = err.(*ErrValidation)
	assert.True(t, ok, "Expected error of type ErrValidation")
	assert.EqualError(t, err, "target_url contains leading or trailing spaces")

	shortPath, err = shortner.CreateShortPath(ctx, expectedShortPath, "://www.google.com")
	assert.Equal(t, "", shortPath)
	assert.Error(t, err, "Expected an error")
	_, ok = err.(*ErrValidation)
	assert.True(t, ok, "Expected error of type ErrValidation")
	assert.EqualError(t, err, "target_url is not valid")

	shortPath, err = shortner.CreateShortPath(ctx, expectedShortPath, "ftp://hello")
	assert.Equal(t, "", shortPath)
	assert.Error(t, err, "Expected an error")
	_, ok = err.(*ErrValidation)
	assert.True(t, ok, "Expected error of type ErrValidation")
	assert.EqualError(t, err, "target_url is not valid")

	targetURLStore.AssertExpectations(t)
	shortPathStore.AssertExpectations(t)
	metrics.AssertExpectations(t)
}

func TestURLShortner_GetTargetURL(t *testing.T) {
	ctx := context.Background()

	targetURLStore := new(mocks.KVStore)
	shortPathStore := new(mocks.KVStore)
	metrics := new(mocks.Metrics)

	shortner, _ := NewURLShortnerBuilder().
		SetTargetURLStore(targetURLStore).
		SetShortPathStore(shortPathStore).
		SetMetrics(metrics).
		Build()

	targetURLStore.On("Get", ctx, "abc123").Return("https://www.google.com", nil)
	targetURLStore.On("Get", ctx, "unknown_shortpath").Return("", store.ErrKeyNotFound)
	targetURLStore.On("Get", ctx, "err_causing_key").Return("", errors.New("connection error"))

	targetURL, err := shortner.GetTargetURL(ctx, "abc123")
	assert.NoError(t, err)
	assert.Equal(t, "https://www.google.com", targetURL)

	_, err = shortner.GetTargetURL(ctx, "unknown_shortpath")
	assert.Error(t, err, "Expected an error")
	_, ok := err.(*ErrNotFound)
	assert.True(t, ok, "Expected error of type ErrNotFound")
	assert.EqualError(t, err, "shortpath mapping not found")

	_, err = shortner.GetTargetURL(ctx, "err_causing_key")
	assert.Error(t, err, "Expected an error")
	_, ok = err.(*ErrServerError)
	assert.True(t, ok, "Expected error of type ErrNotFound")
	assert.EqualError(t, err, "could not lookup shortpath")

	targetURLStore.AssertExpectations(t)
	shortPathStore.AssertExpectations(t)
	metrics.AssertExpectations(t)
}
