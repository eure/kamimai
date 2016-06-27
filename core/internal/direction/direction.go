package direction

const (
	// Unknown XXX
	Unknown = iota

	// Up XXX
	Up

	// Down XXX
	Down
)

// Get returns a string which contains numbers.
func Get(name string) int {

	dotpos := len(name)
	for i := len(name) - 1; 0 <= i; i-- {
		switch name[i] {
		case '.':
			dotpos = i
		case '_':
			switch name[i+1 : dotpos] {
			case "up":
				return Up
			case "down":
				return Down
			}
			return Unknown
		}
	}

	return Unknown
}
