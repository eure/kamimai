package core

import (
	"io/ioutil"
	"os"

	"github.com/mr04vv/kamimai/internal/cast"
)

type (
	// A Migration manages migration files.
	Migration struct {
		version uint64
		name    string
	}

	// A Migrations collects Migration for sorting.
	Migrations []*Migration
)

//////////////////////////////
// Migration
//////////////////////////////

// NewMigration returns a new Migration pointer that can be chained with builder methods to
// set multiple configuration values inline without using pointers.
func NewMigration() *Migration {
	return &Migration{}
}

// WithVersion sets a config version value returning a Config pointer
// for chaining.
func (m *Migration) WithVersion(v interface{}) *Migration {
	m.version = cast.Uint64(v)
	return m
}

// Version returns a migration version value.
func (m Migration) Version() uint64 {
	return m.version
}

// Read reads from file until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.
func (m Migration) Read() ([]byte, error) {
	file, err := os.Open(m.name)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(file)
}

// Name returns a file name.
func (m Migration) Name() string {
	return m.name
}

// IsValid reports whether the migration is valid.
func (m *Migration) IsValid() bool {
	if m == nil {
		return false
	}
	return m.version > 0 && m.name != ""
}

//////////////////////////////
// Migrations
//////////////////////////////

func (m Migrations) index(mig Migration) int {
	for i, v := range m {
		if v.version == mig.version {
			return i
		}
	}
	return int(notFoundIndex)
}

func (m Migrations) first() *Migration {
	if m.Len() == 0 {
		return nil
	}
	return m[0]
}

func (m Migrations) last() *Migration {
	c := m.Len()
	if c == 0 {
		return nil
	}
	return m[c-1]
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
