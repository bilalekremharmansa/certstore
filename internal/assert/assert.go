package assert

import (
	"reflect"
	"strings"
	"testing"
)

func Equal(t *testing.T, expected interface{}, actual interface{}) {
	EqualM(t, expected, actual, "")
}

func EqualM(t *testing.T, expected interface{}, actual interface{}, msg string) {
	if expected != actual {
		t.Fatalf("[ERROR] %s - expected: %v, actual: %v", msg, expected, actual)
	}
}

func NotEqual(t *testing.T, expected interface{}, actual interface{}) {
	NotEqualM(t, expected, actual, "")
}

func NotEqualM(t *testing.T, notExpected interface{}, actual interface{}, msg string) {
	if notExpected == actual {
		t.Fatalf("[ERROR] %s - notExpected: %v, actual: %v", msg, notExpected, actual)
	}
}

func DeepEqual(t *testing.T, expected interface{}, actual interface{}) {
	DeepEqualM(t, expected, actual, "")
}

func DeepEqualM(t *testing.T, expected interface{}, actual interface{}, msg string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("[ERROR] %s - expected: %v, actual: %v", msg, expected, actual)
	}
}

func NotDeepEqual(t *testing.T, expected interface{}, actual interface{}) {
	NotDeepEqualM(t, expected, actual, "")
}

func NotDeepEqualM(t *testing.T, notExpected interface{}, actual interface{}, msg string) {
	if reflect.DeepEqual(notExpected, actual) {
		t.Fatalf("[ERROR] %s - notExpected: %v, actual: %v", msg, notExpected, actual)
	}
}

func True(t *testing.T, isTrue bool) {
	if !isTrue {
		t.Fatal("[ERROR] assertion failed - is not true")
	}
}

func TrueM(t *testing.T, isTrue bool, msg string) {
	if !isTrue {
		t.Fatalf("[ERROR] assertion failed - is not true - %s", msg)
	}
}

func False(t *testing.T, isTrue bool) {
	if isTrue {
		t.Fatal("[ERROR] assertion failed - is not false")
	}
}

func Falsef(t *testing.T, isTrue bool, msg string) {
	if isTrue {
		t.Fatalf("[ERROR] assertion failed - is not false - %s", msg)
	}
}

func Error(t *testing.T, err error, msg string) {
	if err == nil {
		t.Fatalf("[ERROR] error is expected, but not found - %s", msg)
	}
}

func NotError(t *testing.T, err error, msg string) {
	if err != nil {
		t.Fatalf("[ERROR] error is not expected, but found - %v - %s", err, msg)
	}
}

func ErrorContains(t *testing.T, err error, msg string) {
	if err == nil {
		t.Fatalf("[ERROR] error is expected, but not found - %s", msg)
	}

	if !strings.Contains(err.Error(), msg) {
		t.Fatalf("[ERROR] expected error message not found - expected:%s, actual: %s", err.Error(), msg)
	}
}

func Nil(t *testing.T, i interface{}) {
	if i != nil {
		t.Fatal("[ERROR] assertion failed - object is not nil")
	}
}

func NotNil(t *testing.T, i interface{}) {
	if i == nil {
		t.Fatal("[ERROR] assertion failed - object is nil")
	}
}
