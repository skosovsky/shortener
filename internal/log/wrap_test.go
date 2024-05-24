package log_test

import (
	"errors"
	"testing"
	"time"

	"shortener/internal/log"
)

func TestFloat32Attr(t *testing.T) {
	t.Parallel()

	key := "testKey"
	val := float32(123.456)

	attr := log.Float32Attr(key, val)

	// Checking if the attribute correctly converts float32 to float64
	if attr.Key != key || attr.Value.Float64() != float64(val) {
		t.Errorf("Float32Attr() = %v, want %v", attr, log.Float64Attr(key, float64(val)))
	}
}

func TestStringAttr(t *testing.T) {
	t.Parallel()

	key := "testKey"
	val := "testValue"
	attr := log.StringAttr(key, val)
	if attr.Key != key || attr.Value.String() != val {
		t.Errorf("key = %v, want %v", attr, log.StringAttr(key, val))
	}
}

func TestErrAttr(t *testing.T) {
	t.Parallel()

	err := errors.New("test error") //nolint:err113 // test
	attr := log.ErrAttr(err)
	if attr.Key != "error" || attr.Value.String() != err.Error() {
		t.Errorf("ErrAttr() = %v, want %v", attr, log.StringAttr("error", err.Error()))
	}
}

func TestUInt32Attr(t *testing.T) {
	t.Parallel()

	key := "testKey"
	val := uint32(123456789)

	attr := log.UInt32Attr(key, val)

	// The expected behavior is that UInt32Attr converts the uint32 to an int.
	// We need to verify that this conversion is correct.
	expectedValue := int(val)
	if attr.Key != key || int(attr.Value.Int64()) != expectedValue {
		t.Errorf("UInt32Attr() = {Key: %s, Value: %v}, want {Key: %s, Value: %v}",
			attr.Key, attr.Value, key, expectedValue)
	}
}

func TestInt32Attr(t *testing.T) {
	t.Parallel()

	key := "testKey"
	val := int32(12345)

	attr := log.Int32Attr(key, val)

	// Verify that the attribute has the correct key and value
	expectedValue := int(val)
	if attr.Key != key || int(attr.Value.Int64()) != expectedValue {
		t.Errorf("Int32Attr() = {Key: %s, Value: %v}, want {Key: %s, Value: %v}",
			attr.Key, attr.Value, key, expectedValue)
	}
}

func TestTimeAttr(t *testing.T) {
	t.Parallel()

	key := "timestamp"
	val := time.Now()

	attr := log.TimeAttr(key, val)

	// Verify that the attribute has the correct key and the time in string format
	expectedValue := val.String()
	if attr.Key != key || attr.Value.String() != expectedValue {
		t.Errorf("TimeAttr() = {Key: %s, Value: %v}, want {Key: %s, Value: %v}",
			attr.Key, attr.Value, key, expectedValue)
	}
}
