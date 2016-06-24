package core

type (
	// A Migration manages migration files for service.
	Migration struct {
		version uint64
	}

	// A Migrations collects Migration for sorting.
	Migrations []*Migration
)

// NewMigration returns a new Migration pointer that can be chained with builder methods to
// set multiple configuration values inline without using pointers.
func NewMigration() *Migration {
	return &Migration{}
}

// WithVersion sets a config version value returning a Config pointer
// for chaining.
func (m *Migration) WithVersion(v interface{}) *Migration {
	switch v := v.(type) {
	case int:
		m.version = (uint64)(v)
	case int8:
		m.version = (uint64)(v)
	case int16:
		m.version = (uint64)(v)
	case int32:
		m.version = (uint64)(v)
	case int64:
		m.version = (uint64)(v)
	case uint8:
		m.version = (uint64)(v)
	case uint16:
		m.version = (uint64)(v)
	case uint32:
		m.version = (uint64)(v)
	case uint64:
		m.version = v
	}
	return m
}

// Len is the number of elements in the collection.
// Required by Sort Interface{}
func (m Migrations) Len() int {
	return len(m)
}

// Less reports whether the element with
// index i should sort before the element with index j.
// Required by Sort Interface{}
func (m Migrations) Less(i, j int) bool {
	return m[i].version < m[j].version
}

// Swap swaps the elements with indexes i and j.
// Required by Sort Interface{}
func (m Migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
