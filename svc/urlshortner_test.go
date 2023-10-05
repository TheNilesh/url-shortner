package svc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thenilesh/url-shortner/store"

	"github.com/thenilesh/url-shortner/metrics/mocks"
	storemocks "github.com/thenilesh/url-shortner/store/mocks"
)

func TestURLShortner_CreateShortPath(t *testing.T) {
	ctx := context.Background()

	targetURLStore := new(storemocks.KVStore)
	shortPathStore := new(storemocks.KVStore)
	metrics := new(mocks.Metrics)
	collector := new(mocks.Collector)

	shortner := NewURLShortner(6, targetURLStore, shortPathStore, metrics)

	shortPathExpected := "abc123"
	targetURLExpected := "https://www.google.com"
	targetURLStore.On("Get", ctx, shortPathExpected).Return("", store.ErrKeyNotFound)
	targetURLStore.On("Exists", ctx, shortPathExpected).Return(false, nil)
	targetURLStore.On("Put", ctx, shortPathExpected, targetURLExpected).Return(nil)
	// targetURLStore.On("Delete", ctx, shortPathExpected).Return(nil)
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

func TestURLShortner_GetTargetURL(t *testing.T) {
	ctx := context.Background()

	targetURLStore := new(storemocks.KVStore)
	shortPathStore := new(storemocks.KVStore)
	metrics := new(mocks.Metrics)

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
}
