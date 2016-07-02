package core

type (
	// Version interface.
	Version interface {
		// Insert inserts the given migration version.
		Insert(uint64) error

		// Delete deletes the given migration version.
		Delete(uint64) error

		// Count counts number of row the given migration version.
		Count(uint64) int

		// Current returns the current migration version.
		Current() (uint64, error)

		// Create creates
		Create() error

		// Drop drops
		Drop() error
	}
)
