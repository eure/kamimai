package cast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint64(t *testing.T) {
	assert := assert.New(t)

	str := "123"

	candidates := []struct {
		value    interface{}
		expected uint64
		message  string
	}{
		{value: 0, expected: 0, message: ""},
		{value: 1, expected: 1, message: ""},
		{value: 0x10, expected: 16, message: ""},
		{value: "10", expected: 10, message: ""},
		{value: "-10", expected: 0, message: ""},
		{value: "-0", expected: 0, message: ""},
		{value: nil, expected: 0, message: ""},
		{value: (*string)(nil), expected: 0, message: ""},
		{value: str, expected: 123, message: ""},
		{value: &str, expected: 123, message: ""},
	}

	for _, c := range candidates {
		assert.EqualValues(c.expected, Uint64(c.value), c.message)
	}
}
