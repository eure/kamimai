package version

import (
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
