package normalize_key_test

import (
	"testing"

	"github.com/bearer/bearer/internal/util/normalize_key"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		key      string
		expected string
	}{
		{"APIName", "api name"},
		{"userHost", "user host"},
		{"USER_HOST", "user host"},
		{"customer-url", "customer url"},
		{"customer.url", "customer url"},
		{"customer:url", "customer url"},
	}
	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			assert.Equal(t, tt.expected, normalize_key.Normalize(tt.key))
		})
	}
}
