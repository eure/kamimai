package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	assert := assert.New(t)

	candidates := []struct {
		value    string
		expected string
		message  string
	}{
		{value: "000_foo", expected: "000", message: ""},
		{value: "001_foo", expected: "001", message: ""},
		{value: "99999_foo", expected: "99999", message: ""},
		{value: "-000_foo", expected: "", message: ""},
		{value: "-001_foo", expected: "", message: ""},
		{value: "", expected: "", message: ""},
		{value: "foo", expected: "", message: ""},
		{value: "foo_bar", expected: "", message: ""},
	}

	for _, c := range candidates {
		assert.EqualValues(c.expected, Get(c.value), c.message)
	}
}

func BenchmarkGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Get("000_foo")
	}
}
