package version

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	versionReg = regexp.MustCompile("^[0-9]+$")
)

// Get returns a string which contains numbers.
func Get(name string) string {
	for i := 0; i < len(name); i++ {
		if name[i] == '_' {
			str := name[:i]
			if versionReg.FindString(str) != "" {
				return str
			}
			return ""
		}
	}

	return ""
}

// Format returns a version format for printing.
func Format(name string) string {
	ver := Get(name)
	if len(ver) == 0 {
		return ""
	}

	if ver[0] != '0' {
		return "%d"
	}

	return fmt.Sprintf("%%0%dd", len(ver))
}

// IsTimestamp returns if value is timestamp or not
func IsTimestamp(value uint64) bool {
	const layout = "20060102150405"
	str := strconv.FormatUint(value, 10)
	t, err := time.Parse(layout, str)
	if err != nil {
		return false
	}
	// compare by zero-unixtime
	return t.After(time.Unix(0, 0))
}
