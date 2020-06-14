package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	asrt := assert.New(t)

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
		asrt.EqualValues(c.expected, Get(c.value), c.message)
	}
}

func BenchmarkGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Get("000_foo")
	}
}

func TestFormat(t *testing.T) {
	asrt := assert.New(t)

	candidates := []struct {
		value    string
		expected string
		message  string
	}{
		{value: "000_foo", expected: "%03d", message: ""},
		{value: "001_foo", expected: "%03d", message: ""},
		{value: "99999_foo", expected: "%d", message: ""},
		{value: "-000_foo", expected: "", message: ""},
		{value: "-001_foo", expected: "", message: ""},
		{value: "", expected: "", message: ""},
		{value: "foo", expected: "", message: ""},
		{value: "foo_bar", expected: "", message: ""},
		{value: "20200101150405_foo", expected: "%d", message: ""},
		{value: "001_20200101150405_foo", expected: "%03d", message: ""},
		{value: "999_20200101150405_foo", expected: "%d", message: ""},
	}

	for _, c := range candidates {
		asrt.EqualValues(c.expected, Format(c.value), c.message)
	}
}

func BenchmarkFormat(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = Format("000_foo")
	}
}

func TestIsTimestamp(t *testing.T) {
	asrt := assert.New(t)

	candidates := []struct {
		value    uint64
		expected bool
		message  string
	}{
		{value: 0, expected: false, message: ""},
		{value: 1, expected: false, message: ""},
		{value: 99999, expected: false, message: ""},
		{value: 19700101000000, expected: false, message: ""},
		{value: 19700101000001, expected: true, message: ""},
		{value: 20200101150405, expected: true, message: ""},
	}

	for _, c := range candidates {
		asrt.EqualValues(c.expected, IsTimestamp(c.value), c.message)
	}
}
