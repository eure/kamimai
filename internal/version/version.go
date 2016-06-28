package version

import (
	"fmt"
	"regexp"
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
