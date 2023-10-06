package svc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thenilesh/url-shortner/mocks"
)

func TestURLShortnerBuilder_Build(t *testing.T) {
	targetURLStore := new(mocks.KVStore)
	shortPathStore := new(mocks.KVStore)
	metrics := new(mocks.Metrics)

	tests := []struct {
		name       string
		builder    *URLShortnerBuilder
		expected   URLShortner
		errMessage string
	}{
		{
			name: "valid builder",
			builder: NewURLShortnerBuilder().
				SetTargetURLStore(targetURLStore).
				SetShortPathStore(shortPathStore).
				SetMetrics(metrics),
			expected: &urlShortner{
				randomStrGen:   NewRandomStrGen(4, 7, "abcdefghijklmnopqrstuvwxyz0123456789"),
				targetURLStore: targetURLStore,
				shortPathStore: shortPathStore,
				metrics:        metrics,
			},
			errMessage: "",
		},
		{
			name: "nil targetURLStore",
			builder: NewURLShortnerBuilder().
				SetShortPathStore(shortPathStore).
				SetMetrics(metrics),
			expected:   nil,
			errMessage: "targetURLStore is nil",
		},
		{
			name: "nil shortPathStore",
			builder: NewURLShortnerBuilder().
				SetTargetURLStore(targetURLStore).
				SetMetrics(metrics),
			expected:   nil,
			errMessage: "shortPathStore is nil",
		},
		{
			name: "nil metrics",
			builder: NewURLShortnerBuilder().
				SetTargetURLStore(targetURLStore).
				SetShortPathStore(shortPathStore),
			expected:   nil,
			errMessage: "metrics is nil",
		},
		{
			name: "minLength greater than maxLength",
			builder: NewURLShortnerBuilder().
				SetTargetURLStore(targetURLStore).
				SetShortPathStore(shortPathStore).
				SetMetrics(metrics).
				SetMinLength(10).
				SetMaxLength(5),
			expected:   nil,
			errMessage: "minLength is greater than maxLength",
		},
		{
			name: "minLength less than or equal to 0",
			builder: NewURLShortnerBuilder().
				SetTargetURLStore(targetURLStore).
				SetShortPathStore(shortPathStore).
				SetMetrics(metrics).
				SetMinLength(0),
			expected:   nil,
			errMessage: "minLength or maxLength is less than or equal to 0",
		},
		{
			name: "maxLength less than or equal to 0",
			builder: NewURLShortnerBuilder().
				SetTargetURLStore(targetURLStore).
				SetShortPathStore(shortPathStore).
				SetMetrics(metrics).
				SetMinLength(0).
				SetMaxLength(0),
			expected:   nil,
			errMessage: "minLength or maxLength is less than or equal to 0",
		},
		{
			name: "minLength greater than 50",
			builder: NewURLShortnerBuilder().
				SetTargetURLStore(targetURLStore).
				SetShortPathStore(shortPathStore).
				SetMetrics(metrics).
				SetMinLength(51),
			expected:   nil,
			errMessage: "minLength or maxLength is greater than 50",
		},
		{
			name: "maxLength greater than 50",
			builder: NewURLShortnerBuilder().
				SetTargetURLStore(targetURLStore).
				SetShortPathStore(shortPathStore).
				SetMetrics(metrics).
				SetMaxLength(51),
			expected:   nil,
			errMessage: "minLength or maxLength is greater than 50",
		},
		{
			name: "invalid charset",
			builder: NewURLShortnerBuilder().
				SetTargetURLStore(targetURLStore).
				SetShortPathStore(shortPathStore).
				SetMetrics(metrics).
				SetCharset("abc$"),
			expected:   nil,
			errMessage: "charset contains invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.builder.Build()
			if tt.errMessage != "" {
				assert.EqualError(t, err, tt.errMessage)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
