package core

type (
	// A Migration manages migration files.
	Migration struct {
		current uint64
	}
	// A Migrations collects Migration for sorting.
	Migrations []Migration
)

func (m *Migration) WithCurrent(v interface{}) *Migration {
	switch v := v.(type) {
	case uint64:
		m.current = v
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
	return m[i].Version < m[j].Version
}

// Swap swaps the elements with indexes i and j.
// Required by Sort Interface{}
func (m Migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
