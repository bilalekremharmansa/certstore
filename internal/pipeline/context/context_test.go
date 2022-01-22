package context

import (
	"testing"

	"bilalekrem.com/certstore/internal/assert"
)

const TEST_KEY Key = "test-key"

func TestNew(t *testing.T) {
	ctx := New()
	assert.NotNil(t, ctx)
	assert.NotNil(t, ctx.values)
}

func TestStoreAndGetValue(t *testing.T) {
	ctx := New()

	value := "test value"
	ctx.StoreValue(TEST_KEY, value)

	retrievedValue := ctx.GetValue(TEST_KEY)
	assert.Equal(t, value, retrievedValue)
}
