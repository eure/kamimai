package core

type (
	// Version interface.
	Version interface {
		// Insert inserts the given migration version.
		Insert(uint64) error

		// Delete deletes the given migration version.
		Delete(uint64) error

		// Current returns the current migration version.
		Current() (uint64, error)

		// Create creates
		Create() error

		// Drop drops
		Drop() error
	}
)
