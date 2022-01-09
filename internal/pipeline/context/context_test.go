package context

import (
	"testing"
)

const TEST_KEY Key = "test-key"

func TestNew(t *testing.T) {
	ctx := New()
	if ctx == nil {
		t.Fatalf("New context creation failed")
	}

	if ctx.values == nil {
		t.Fatalf("values map is nil, should've been initialized")
	}
}

func TestStoreAndGetValue(t *testing.T) {
	ctx := New()

	value := "test value"
	ctx.StoreValue(TEST_KEY, value)

	retrievedValue := ctx.GetValue(TEST_KEY)
	if retrievedValue != value {
		t.Fatalf("Stored value is not correct, expected: %s, found: %s", value, retrievedValue)
	}
}
