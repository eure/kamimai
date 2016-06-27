package direction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	assert := assert.New(t)

	candidates := []struct {
		value    string
		expected int
		message  string
	}{
		{value: "000_foo_up", expected: Up, message: ""},
		{value: "001_foo_up.sql", expected: Up, message: ""},
		{value: "000_foo_down", expected: Down, message: ""},
		{value: "001_foo_down.sql", expected: Down, message: ""},
		{value: "", expected: Unknown, message: ""},
		{value: "foo_up_bar", expected: Unknown, message: ""},
		{value: "foo_down_bar", expected: Unknown, message: ""},
	}

	for _, c := range candidates {
		assert.EqualValues(c.expected, Get(c.value), c.message)
	}
}

func BenchmarkGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Get("000_foo_up")
	}
}
